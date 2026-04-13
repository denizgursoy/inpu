package otel

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type config struct {
	meterProvider  metric.MeterProvider
	tracerProvider trace.TracerProvider
	propagator     propagation.TextMapPropagator
	metricsEnabled bool
	tracingEnabled bool
}

// Option applies a configuration option to the OTel middleware.
type Option func(*config)

// WithMeterProvider sets a custom OTel MeterProvider.
// If not set, the global MeterProvider is used.
func WithMeterProvider(mp metric.MeterProvider) Option {
	return func(c *config) {
		c.meterProvider = mp
	}
}

// WithTracerProvider sets a custom OTel TracerProvider.
// If not set, the global TracerProvider is used.
func WithTracerProvider(tp trace.TracerProvider) Option {
	return func(c *config) {
		c.tracerProvider = tp
	}
}

// WithPropagator sets a custom OTel TextMapPropagator for trace context injection.
// If not set, the global TextMapPropagator is used.
func WithPropagator(p propagation.TextMapPropagator) Option {
	return func(c *config) {
		c.propagator = p
	}
}

// WithoutMetrics disables metric collection.
// Tracing will still be active unless also disabled with WithoutTracing.
func WithoutMetrics() Option {
	return func(c *config) {
		c.metricsEnabled = false
	}
}

// WithoutTracing disables distributed tracing and context propagation.
// Metrics will still be active unless also disabled with WithoutMetrics.
func WithoutTracing() Option {
	return func(c *config) {
		c.tracingEnabled = false
	}
}
