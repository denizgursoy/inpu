package otel

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/denizgursoy/inpu"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	instrumentationName = "github.com/denizgursoy/inpu/middlewares/otel"

	// Custom attribute keys
	attrKeyRequestID   = attribute.Key("inpu.request.id")
	attrKeyResendCount = attribute.Key("http.resend_count")
)

type otelMiddleware struct {
	next       http.RoundTripper
	tracer     trace.Tracer
	propagator propagation.TextMapPropagator
	cfg        config

	// Metric instruments
	requestDuration  metric.Float64Histogram
	requestBodySize  metric.Int64Histogram
	responseBodySize metric.Int64Histogram
	activeRequests   metric.Int64UpDownCounter
	requestTotal     metric.Int64Counter
	retryTotal       metric.Int64Counter
}

// NewMiddleware creates an OTel observability middleware that collects metrics and traces
// for outgoing HTTP requests. It implements the inpu.Middleware interface.
//
// By default, it uses the global OTel MeterProvider, TracerProvider, and TextMapPropagator.
// Use WithMeterProvider, WithTracerProvider, and WithPropagator to override.
//
// Usage:
//
//	client := inpu.NewClient("https://api.example.com").
//		Use(otel.NewMiddleware())
func NewMiddleware(opts ...Option) inpu.Middleware {
	cfg := config{
		meterProvider:  otel.GetMeterProvider(),
		tracerProvider: otel.GetTracerProvider(),
		propagator:     otel.GetTextMapPropagator(),
		metricsEnabled: true,
		tracingEnabled: true,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	meter := cfg.meterProvider.Meter(instrumentationName)
	tracer := cfg.tracerProvider.Tracer(instrumentationName)

	m := &otelMiddleware{
		tracer:     tracer,
		propagator: cfg.propagator,
		cfg:        cfg,
	}

	if cfg.metricsEnabled {
		m.requestDuration = must(meter.Float64Histogram(
			"http.client.request.duration",
			metric.WithDescription("Duration of HTTP client requests"),
			metric.WithUnit("s"),
		))
		m.requestBodySize = must(meter.Int64Histogram(
			"http.client.request.body.size",
			metric.WithDescription("Size of HTTP client request bodies"),
			metric.WithUnit("By"),
		))
		m.responseBodySize = must(meter.Int64Histogram(
			"http.client.response.body.size",
			metric.WithDescription("Size of HTTP client response bodies"),
			metric.WithUnit("By"),
		))
		m.activeRequests = must(meter.Int64UpDownCounter(
			"http.client.active_requests",
			metric.WithDescription("Number of active HTTP client requests"),
			metric.WithUnit("{request}"),
		))
		m.requestTotal = must(meter.Int64Counter(
			"http.client.request.total",
			metric.WithDescription("Total number of HTTP client requests"),
			metric.WithUnit("{request}"),
		))
		m.retryTotal = must(meter.Int64Counter(
			"http.client.request.retry.count",
			metric.WithDescription("Total number of HTTP client request retries"),
			metric.WithUnit("{retry}"),
		))
	}

	return m
}

func (m *otelMiddleware) ID() string {
	return "otel-collector-middleware"
}

func (m *otelMiddleware) Priority() int {
	return 2
}

func (m *otelMiddleware) Apply(next http.RoundTripper) http.RoundTripper {
	m.next = next

	return m
}

func (m *otelMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	// Build common attributes
	baseAttrs := m.baseAttributes(req)

	// Extract request ID if available
	if requestID := inpu.ExtractRequestIDFromContext(ctx); requestID != nil {
		baseAttrs = append(baseAttrs, attrKeyRequestID.String(*requestID))
	}

	// Extract retry attempt
	resendCount := inpu.ExtractRetryAttemptFromContext(ctx)

	// Start span if tracing is enabled
	if m.cfg.tracingEnabled {
		spanName := fmt.Sprintf("%s %s", req.Method, req.URL.Hostname())
		var span trace.Span
		ctx, span = m.tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(baseAttrs...),
			trace.WithAttributes(attrKeyResendCount.Int(resendCount)),
		)
		defer span.End()

		// Inject trace context into outgoing request headers
		m.propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

		// Update request with new context
		req = req.WithContext(ctx)
	}

	// Record pre-request metrics
	if m.cfg.metricsEnabled {
		activeAttrs := metric.WithAttributes(baseAttrs...)
		m.activeRequests.Add(ctx, 1, activeAttrs)

		// Record request body size
		if req.ContentLength > 0 {
			m.requestBodySize.Record(ctx, req.ContentLength, activeAttrs)
		}

		// Record retry metric
		if resendCount > 0 {
			m.retryTotal.Add(ctx, 1, activeAttrs)
		}
	}

	// Execute the request
	start := time.Now()
	resp, err := m.next.RoundTrip(req)
	duration := time.Since(start).Seconds()

	// Build response attributes
	responseAttrs := make([]attribute.KeyValue, 0, len(baseAttrs)+3)
	responseAttrs = append(responseAttrs, baseAttrs...)
	responseAttrs = append(responseAttrs, attrKeyResendCount.Int(resendCount))

	if err != nil {
		responseAttrs = append(responseAttrs, attribute.String("error.type", errorType(err)))
	}

	if resp != nil {
		responseAttrs = append(responseAttrs, semconv.HTTPResponseStatusCode(resp.StatusCode))
	}

	// Record post-request metrics
	if m.cfg.metricsEnabled {
		respMetricAttrs := metric.WithAttributes(responseAttrs...)

		m.activeRequests.Add(ctx, -1, metric.WithAttributes(baseAttrs...))
		m.requestDuration.Record(ctx, duration, respMetricAttrs)
		m.requestTotal.Add(ctx, 1, respMetricAttrs)

		// Record response body size
		if resp != nil && resp.ContentLength > 0 {
			m.responseBodySize.Record(ctx, resp.ContentLength, respMetricAttrs)
		}
	}

	// Update span with response info
	if m.cfg.tracingEnabled {
		span := trace.SpanFromContext(ctx)
		if resp != nil {
			span.SetAttributes(semconv.HTTPResponseStatusCode(resp.StatusCode))
			span.SetAttributes(attrKeyResendCount.Int(resendCount))

			if resp.StatusCode >= http.StatusBadRequest {
				span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", resp.StatusCode))
			}
		}

		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.SetAttributes(attribute.String("error.type", errorType(err)))
		}
	}

	return resp, err
}

// baseAttributes returns common OTel attributes for a request.
func (m *otelMiddleware) baseAttributes(req *http.Request) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.HTTPRequestMethodKey.String(req.Method),
		semconv.ServerAddress(req.URL.Hostname()),
		semconv.URLScheme(req.URL.Scheme),
	}

	port := req.URL.Port()
	if port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			attrs = append(attrs, semconv.ServerPort(p))
		}
	}

	return attrs
}

// errorType returns a short error type string suitable for the error.type attribute.
func errorType(err error) string {
	if err == nil {
		return ""
	}

	t := fmt.Sprintf("%T", err)
	if t == "*errors.errorString" || t == "*fmt.wrapError" {
		return err.Error()
	}

	return t
}

// must panics if the metric instrument creation returns an error.
// OTel instrument creation only fails on invalid parameters, which indicates a programming error.
func must[T any](val T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("otel middleware: failed to create metric instrument: %v", err))
	}

	return val
}
