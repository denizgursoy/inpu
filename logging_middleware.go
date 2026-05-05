package inpu

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"
)

const defaultMaxBodyLogSize = 4096

var sensitiveHeaders = []string{HeaderAuthorization, HeaderAPISecret, HeaderAPIKey, HeaderAPIToken, HeaderCookie}

// LoggingOption configures the logging middleware.
type LoggingOption func(*loggingMiddleware)

// WithVerbose enables verbose logging of headers and bodies.
func WithVerbose() LoggingOption {
	return func(l *loggingMiddleware) {
		l.verbose = true
	}
}

// WithDisabled creates the middleware in a disabled state (no-op passthrough).
func WithDisabled() LoggingOption {
	return func(l *loggingMiddleware) {
		l.disabled = true
	}
}

// WithMaxBodyLogSize sets the maximum number of bytes to log for request/response bodies.
// Bodies exceeding this size are truncated. Default is 4096 bytes.
func WithMaxBodyLogSize(n int) LoggingOption {
	return func(l *loggingMiddleware) {
		l.maxBodyLogSize = n
	}
}

type loggingMiddleware struct {
	verbose        bool
	disabled       bool
	maxBodyLogSize int
	next           http.RoundTripper
}

// NewLoggingMiddleware creates a logging middleware with the provided options.
func NewLoggingMiddleware(opts ...LoggingOption) Middleware {
	m := &loggingMiddleware{
		maxBodyLogSize: defaultMaxBodyLogSize,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// LoggingMiddleware creates a logging middleware.
//
// Deprecated: Use NewLoggingMiddleware(WithVerbose(), ...) instead.
func LoggingMiddleware(verbose, disabled bool) Middleware {
	return &loggingMiddleware{
		verbose:        verbose,
		disabled:       disabled,
		maxBodyLogSize: defaultMaxBodyLogSize,
	}
}

func (t *loggingMiddleware) ID() string {
	return "default-logging-middleware"
}

func (t *loggingMiddleware) Priority() int {
	return 1
}

func (t *loggingMiddleware) Apply(next http.RoundTripper) http.RoundTripper {
	t.next = next

	return t
}

func (t *loggingMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.disabled {
		return t.next.RoundTrip(req)
	}
	ctx := req.Context()
	logger := ExtractLoggerFromContext(ctx)
	start := time.Now()

	// Log request
	logger.Info(ctx, "→ [%s] %s", req.Method, req.URL.Redacted())

	if t.verbose {
		logger.Info(ctx, "  Headers: %v", headersToString(req.Header))
		if req.Body != nil {
			body, _ := io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(body))
			logger.Info(ctx, "  Body: %s", t.truncateBody(body))
		}
	}

	// Execute request
	resp, err := t.next.RoundTrip(req)
	duration := time.Since(start)

	// Log response
	if err != nil {
		if t.verbose && resp != nil && resp.Body != nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body = io.NopCloser(bytes.NewBuffer(body))
			logger.Info(ctx, "  Error Response Body: %s", t.truncateBody(body))
		}
		logger.Error(ctx, err, "← [%s] %s - ERROR: %v (took %v)", req.Method, req.URL.Redacted(), err, duration)

		return resp, err
	}

	logger.Info(ctx, "← [%s] %s - Status: %d - Duration: %v", req.Method, req.URL.Redacted(), resp.StatusCode, duration)

	if t.verbose {
		logger.Info(ctx, "  Response Headers: %v", headersToString(resp.Header))
		if resp.Body != nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body = io.NopCloser(bytes.NewBuffer(body))
			logger.Info(ctx, "  Response Body: %s", t.truncateBody(body))
		}
	}

	return resp, nil
}

func (t *loggingMiddleware) truncateBody(body []byte) string {
	if t.maxBodyLogSize > 0 && len(body) > t.maxBodyLogSize {
		return fmt.Sprintf("%s... (truncated, %d bytes total)", string(body[:t.maxBodyLogSize]), len(body))
	}
	return string(body)
}

func headersToString(headers http.Header) string {
	var parts []string

	for key, values := range headers {
		for _, value := range values {
			if slices.Contains(sensitiveHeaders, key) {
				value = strings.Repeat("X", len(value))
			}
			parts = append(parts, key+"="+value)
		}
	}
	return strings.Join(parts, "; ")
}
