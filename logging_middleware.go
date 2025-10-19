package inpu

import (
	"bytes"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"
)

var sensitiveHeaders = []string{HeaderAuthorization, HeaderAPISecret, HeaderAPIKey, HeaderAPIToken, HeaderCookie}

// LogLevel represents the logging verbosity level
type LogerLevel int

const (
	LogLevelDisabled LogerLevel = iota
	LogLevelSimple
	LogLevelVerbose
)

type loggingMiddleware struct {
	logLevel LogerLevel
	next     http.RoundTripper
}

// LoggingMiddleware creates a logging middleware
func LoggingMiddleware(logLevel LogerLevel) Middleware {
	return &loggingMiddleware{
		logLevel: logLevel,
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
	if t.logLevel == LogLevelDisabled {
		return t.next.RoundTrip(req)
	}
	ctx := req.Context()
	logger := GetLoggerFromContext(ctx)
	start := time.Now()

	// Log request
	logger.Info(ctx, "→ [%s] %s", req.Method, req.URL.Redacted())

	isVerbose := t.logLevel == LogLevelVerbose
	if isVerbose {
		logger.Info(ctx, "  Headers: %v", headersToString(req.Header))
		if req.Body != nil {
			body, _ := io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(body))
			logger.Info(ctx, "  Body: %s", string(body))
		}
	}

	// Execute request
	resp, err := t.next.RoundTrip(req)
	duration := time.Since(start)

	// Log response
	if err != nil {
		logger.Error(ctx, err, "← [%s] %s - ERROR: %v (took %v)", req.Method, req.URL.Redacted(), err, duration)

		return resp, err
	}

	logger.Info(ctx, "← [%s] %s - Status: %d - Duration: %v", req.Method, req.URL.Redacted(), resp.StatusCode, duration)

	if isVerbose {
		logger.Info(ctx, "  Response Headers: %v", headersToString(resp.Header))
		if resp.Body != nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body = io.NopCloser(bytes.NewBuffer(body))
			logger.Info(ctx, "  Response Body: %s", string(body))
		}
	}

	return resp, nil
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
