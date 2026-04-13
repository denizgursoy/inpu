package otel

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/denizgursoy/inpu"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// helpers

func setupTestProviders(t *testing.T) (*tracetest.InMemoryExporter, *sdkmetric.ManualReader, []Option) {
	t.Helper()

	spanExporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(spanExporter))
	t.Cleanup(func() { tp.Shutdown(context.Background()) })

	metricReader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(metricReader))
	t.Cleanup(func() { mp.Shutdown(context.Background()) })

	opts := []Option{
		WithTracerProvider(tp),
		WithMeterProvider(mp),
		WithPropagator(propagation.TraceContext{}),
	}

	return spanExporter, metricReader, opts
}

func collectMetrics(t *testing.T, reader *sdkmetric.ManualReader) metricdata.ResourceMetrics {
	t.Helper()

	var rm metricdata.ResourceMetrics
	if err := reader.Collect(context.Background(), &rm); err != nil {
		t.Fatalf("failed to collect metrics: %v", err)
	}

	return rm
}

func findMetric(rm metricdata.ResourceMetrics, name string) *metricdata.Metrics {
	for _, sm := range rm.ScopeMetrics {
		for i := range sm.Metrics {
			if sm.Metrics[i].Name == name {
				return &sm.Metrics[i]
			}
		}
	}

	return nil
}

func hasAttribute(attrs attribute.Set, key attribute.Key, val string) bool {
	v, ok := attrs.Value(key)
	if !ok {
		return false
	}

	return v.AsString() == val
}

func hasIntAttribute(attrs attribute.Set, key attribute.Key, val int64) bool {
	v, ok := attrs.Value(key)
	if !ok {
		return false
	}

	return v.AsInt64() == val
}

// Tests

func TestMiddleware_MetricsRecorded(t *testing.T) {
	spanExporter, metricReader, opts := setupTestProviders(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "13")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	mw := NewMiddleware(opts...)
	client := inpu.New().BasePath(server.URL).Use(mw)

	err := client.Get("/test").Send()
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	rm := collectMetrics(t, metricReader)

	// Verify all expected metrics exist
	expectedMetrics := []string{
		"http.client.request.duration",
		"http.client.active_requests",
		"http.client.request.total",
	}
	for _, name := range expectedMetrics {
		if m := findMetric(rm, name); m == nil {
			t.Errorf("expected metric %q not found", name)
		}
	}

	// Verify request.total has correct attributes
	totalMetric := findMetric(rm, "http.client.request.total")
	if totalMetric == nil {
		t.Fatal("http.client.request.total not found")
	}

	sumData, ok := totalMetric.Data.(metricdata.Sum[int64])
	if !ok {
		t.Fatalf("expected Sum[int64], got %T", totalMetric.Data)
	}

	if len(sumData.DataPoints) == 0 {
		t.Fatal("no data points in request.total")
	}

	dp := sumData.DataPoints[0]
	if dp.Value != 1 {
		t.Errorf("expected request.total value 1, got %d", dp.Value)
	}

	if !hasAttribute(dp.Attributes, "http.request.method", "GET") {
		t.Error("missing http.request.method=GET attribute")
	}

	if !hasIntAttribute(dp.Attributes, "http.response.status_code", 200) {
		t.Error("missing http.response.status_code=200 attribute")
	}

	// Verify span was created
	spans := spanExporter.GetSpans()
	if len(spans) == 0 {
		t.Fatal("expected at least one span")
	}

	span := spans[0]
	if span.SpanKind != trace.SpanKindClient {
		t.Errorf("expected SpanKindClient, got %v", span.SpanKind)
	}
}

func TestMiddleware_SpanNameAndAttributes(t *testing.T) {
	spanExporter, _, opts := setupTestProviders(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	mw := NewMiddleware(opts...)
	client := inpu.New().BasePath(server.URL).Use(mw)

	err := client.Get("/resource").Send()
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	spans := spanExporter.GetSpans()
	if len(spans) != 1 {
		t.Fatalf("expected 1 span, got %d", len(spans))
	}

	span := spans[0]

	// Span name should be "METHOD hostname"
	expectedName := "GET 127.0.0.1"
	if span.Name != expectedName {
		t.Errorf("expected span name %q, got %q", expectedName, span.Name)
	}

	// Check span attributes contain http.request.method
	found := false
	for _, attr := range span.Attributes {
		if attr.Key == "http.request.method" && attr.Value.AsString() == "GET" {
			found = true

			break
		}
	}

	if !found {
		t.Error("span missing http.request.method=GET attribute")
	}
}

func TestMiddleware_TraceContextPropagation(t *testing.T) {
	_, _, opts := setupTestProviders(t)

	var receivedTraceparent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedTraceparent = r.Header.Get("Traceparent")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	mw := NewMiddleware(opts...)
	client := inpu.New().BasePath(server.URL).Use(mw)

	err := client.Get("/propagate").Send()
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if receivedTraceparent == "" {
		t.Error("expected traceparent header to be injected, got empty")
	}
}

func TestMiddleware_ErrorSetsSpanStatus(t *testing.T) {
	spanExporter, metricReader, opts := setupTestProviders(t)

	// Use a server that will be closed to trigger a connection error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	serverURL := server.URL
	server.Close() // close immediately to cause connection error

	mw := NewMiddleware(opts...)
	client := inpu.New().BasePath(serverURL).Use(mw)

	err := client.Get("/fail").Send()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Check span has error status
	spans := spanExporter.GetSpans()
	if len(spans) == 0 {
		t.Fatal("expected at least one span")
	}

	span := spans[0]
	if span.Status.Code != codes.Error {
		t.Errorf("expected span status Error, got %v", span.Status.Code)
	}

	// Check error.type attribute exists in metrics
	rm := collectMetrics(t, metricReader)
	totalMetric := findMetric(rm, "http.client.request.total")

	if totalMetric == nil {
		t.Fatal("http.client.request.total not found")
	}

	sumData, ok := totalMetric.Data.(metricdata.Sum[int64])
	if !ok {
		t.Fatalf("expected Sum[int64], got %T", totalMetric.Data)
	}

	if len(sumData.DataPoints) == 0 {
		t.Fatal("no data points")
	}

	// Verify error.type attribute is present
	dp := sumData.DataPoints[0]
	_, hasErrType := dp.Attributes.Value("error.type")
	if !hasErrType {
		t.Error("expected error.type attribute on metrics for failed request")
	}
}

func TestMiddleware_HTTP4xxSetsSpanError(t *testing.T) {
	spanExporter, _, opts := setupTestProviders(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	mw := NewMiddleware(opts...)
	client := inpu.New().BasePath(server.URL).Use(mw)

	_ = client.Get("/not-found").Send()

	spans := spanExporter.GetSpans()
	if len(spans) == 0 {
		t.Fatal("expected at least one span")
	}

	span := spans[0]
	if span.Status.Code != codes.Error {
		t.Errorf("expected span status Error for 404, got %v", span.Status.Code)
	}

	if span.Status.Description != "HTTP 404" {
		t.Errorf("expected span status description 'HTTP 404', got %q", span.Status.Description)
	}
}

func TestMiddleware_RequestDurationHistogram(t *testing.T) {
	_, metricReader, opts := setupTestProviders(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	mw := NewMiddleware(opts...)
	client := inpu.New().BasePath(server.URL).Use(mw)

	err := client.Get("/duration").Send()
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	rm := collectMetrics(t, metricReader)
	durationMetric := findMetric(rm, "http.client.request.duration")

	if durationMetric == nil {
		t.Fatal("http.client.request.duration metric not found")
	}

	if durationMetric.Unit != "s" {
		t.Errorf("expected unit 's', got %q", durationMetric.Unit)
	}

	histData, ok := durationMetric.Data.(metricdata.Histogram[float64])
	if !ok {
		t.Fatalf("expected Histogram[float64], got %T", durationMetric.Data)
	}

	if len(histData.DataPoints) == 0 {
		t.Fatal("no data points in duration histogram")
	}

	dp := histData.DataPoints[0]
	if dp.Count != 1 {
		t.Errorf("expected count 1, got %d", dp.Count)
	}

	if dp.Sum <= 0 {
		t.Error("expected positive duration sum")
	}
}

func TestMiddleware_WithoutMetrics(t *testing.T) {
	spanExporter, metricReader, opts := setupTestProviders(t)
	opts = append(opts, WithoutMetrics())

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	mw := NewMiddleware(opts...)
	client := inpu.New().BasePath(server.URL).Use(mw)

	err := client.Get("/no-metrics").Send()
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	// Metrics should be empty
	rm := collectMetrics(t, metricReader)
	for _, sm := range rm.ScopeMetrics {
		if len(sm.Metrics) > 0 {
			t.Errorf("expected no metrics when WithoutMetrics is set, got %d", len(sm.Metrics))
		}
	}

	// But spans should still be created
	spans := spanExporter.GetSpans()
	if len(spans) == 0 {
		t.Error("expected spans even with WithoutMetrics")
	}
}

func TestMiddleware_WithoutTracing(t *testing.T) {
	spanExporter, metricReader, opts := setupTestProviders(t)
	opts = append(opts, WithoutTracing())

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	mw := NewMiddleware(opts...)
	client := inpu.New().BasePath(server.URL).Use(mw)

	err := client.Get("/no-tracing").Send()
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	// Spans should be empty
	spans := spanExporter.GetSpans()
	if len(spans) != 0 {
		t.Errorf("expected no spans when WithoutTracing is set, got %d", len(spans))
	}

	// But metrics should still be recorded
	rm := collectMetrics(t, metricReader)
	if findMetric(rm, "http.client.request.total") == nil {
		t.Error("expected http.client.request.total metric even with WithoutTracing")
	}
}

func TestMiddleware_ResendCountAttribute(t *testing.T) {
	_, metricReader, opts := setupTestProviders(t)

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount <= 1 {
			w.WriteHeader(http.StatusServiceUnavailable)

			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	mw := NewMiddleware(opts...)
	retryMw := inpu.RetryMiddleware(1)
	client := inpu.New().BasePath(server.URL).Use(mw, retryMw)

	err := client.Get("/retry").Send()
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if callCount != 2 {
		t.Fatalf("expected 2 calls, got %d", callCount)
	}

	rm := collectMetrics(t, metricReader)
	totalMetric := findMetric(rm, "http.client.request.total")

	if totalMetric == nil {
		t.Fatal("http.client.request.total not found")
	}

	sumData, ok := totalMetric.Data.(metricdata.Sum[int64])
	if !ok {
		t.Fatalf("expected Sum[int64], got %T", totalMetric.Data)
	}

	// Should have 2 data points (one for attempt 0, one for attempt 1)
	// or they could be aggregated differently based on attributes
	if len(sumData.DataPoints) < 1 {
		t.Fatal("expected at least 1 data point")
	}

	// Verify total requests is 2
	var total int64
	for _, dp := range sumData.DataPoints {
		total += dp.Value
	}

	if total != 2 {
		t.Errorf("expected total of 2 requests, got %d", total)
	}
}

func TestMiddleware_RetryCountMetric(t *testing.T) {
	_, metricReader, opts := setupTestProviders(t)

	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount <= 1 {
			w.WriteHeader(http.StatusServiceUnavailable)

			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	mw := NewMiddleware(opts...)
	retryMw := inpu.RetryMiddlewareWithConfig(inpu.RetryConfig{
		MaxRetries:        1,
		InitialBackoff:    0, // no delay for test speed
		BackoffMultiplier: 1.0,
	})
	client := inpu.New().BasePath(server.URL).Use(mw, retryMw)

	err := client.Get("/retry-count").Send()
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	rm := collectMetrics(t, metricReader)
	retryMetric := findMetric(rm, "http.client.request.retry.count")

	if retryMetric == nil {
		t.Fatal("http.client.request.retry.count metric not found")
	}

	sumData, ok := retryMetric.Data.(metricdata.Sum[int64])
	if !ok {
		t.Fatalf("expected Sum[int64], got %T", retryMetric.Data)
	}

	if len(sumData.DataPoints) == 0 {
		t.Fatal("no data points in retry count")
	}

	// The retry count should be 1 (only the second attempt is a retry)
	var retryCount int64
	for _, dp := range sumData.DataPoints {
		retryCount += dp.Value
	}

	if retryCount != 1 {
		t.Errorf("expected retry count of 1, got %d", retryCount)
	}
}

func TestMiddleware_RequestIDAttribute(t *testing.T) {
	_, metricReader, opts := setupTestProviders(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	mw := NewMiddleware(opts...)
	requestIDMw := inpu.RequestIDMiddleware()
	client := inpu.New().BasePath(server.URL).Use(mw, requestIDMw)

	err := client.Get("/with-id").Send()
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	rm := collectMetrics(t, metricReader)
	totalMetric := findMetric(rm, "http.client.request.total")

	if totalMetric == nil {
		t.Fatal("http.client.request.total not found")
	}

	sumData, ok := totalMetric.Data.(metricdata.Sum[int64])
	if !ok {
		t.Fatalf("expected Sum[int64], got %T", totalMetric.Data)
	}

	if len(sumData.DataPoints) == 0 {
		t.Fatal("no data points")
	}

	dp := sumData.DataPoints[0]
	_, hasRequestID := dp.Attributes.Value("inpu.request.id")
	if !hasRequestID {
		t.Error("expected inpu.request.id attribute on metrics when RequestIDMiddleware is active")
	}
}

func TestMiddleware_ActiveRequestsUpDown(t *testing.T) {
	_, metricReader, opts := setupTestProviders(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	mw := NewMiddleware(opts...)
	client := inpu.New().BasePath(server.URL).Use(mw)

	err := client.Get("/active").Send()
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	rm := collectMetrics(t, metricReader)
	activeMetric := findMetric(rm, "http.client.active_requests")

	if activeMetric == nil {
		t.Fatal("http.client.active_requests metric not found")
	}

	sumData, ok := activeMetric.Data.(metricdata.Sum[int64])
	if !ok {
		t.Fatalf("expected Sum[int64], got %T", activeMetric.Data)
	}

	if len(sumData.DataPoints) == 0 {
		t.Fatal("no data points")
	}

	// After request completes, active requests should be back to 0
	dp := sumData.DataPoints[0]
	if dp.Value != 0 {
		t.Errorf("expected active_requests to be 0 after request completes, got %d", dp.Value)
	}
}

func TestMiddleware_IDAndPriority(t *testing.T) {
	mw := NewMiddleware()

	if mw.ID() != "otel-collector-middleware" {
		t.Errorf("expected ID 'otel-collector-middleware', got %q", mw.ID())
	}

	if mw.Priority() != 2 {
		t.Errorf("expected Priority 2, got %d", mw.Priority())
	}
}
