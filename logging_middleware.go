package inpu

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"
)

var sensitiveHeaders = []string{HeaderAuthorization, HeaderAPISecret, HeaderAPIKey, HeaderAPIToken, HeaderCookie}

type Logger interface {
	Printf(string, ...interface{})
}

// LogLevel represents the logging verbosity level
type LogLevel int

const (
	LogLevelDisabled LogLevel = iota
	LogLevelInfo
	LogLevelVerbose
)

type loggingMiddleware struct {
	logger   Logger
	logLevel LogLevel
	next     http.RoundTripper
}

// LoggingMiddleware creates a logging middleware
func LoggingMiddleware(logLevel LogLevel) Middleware {
	return &loggingMiddleware{
		logLevel: logLevel,
		logger:   log.Default(),
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

	start := time.Now()

	// Log request
	t.logger.Printf("→ [%s] %s", req.Method, req.URL.Redacted())

	isVerbose := t.logLevel == LogLevelVerbose
	if isVerbose {
		t.logger.Printf("  Headers: %v", headersToString(req.Header))
		if req.Body != nil {
			body, _ := io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(body))
			t.logger.Printf("  Body: %s", string(body))
		}
	}

	// Execute request
	resp, err := t.next.RoundTrip(req)
	duration := time.Since(start)

	// Log response
	if err != nil {
		t.logger.Printf("← [%s] %s - ERROR: %v (took %v)", req.Method, req.URL.Redacted(), err, duration)

		return resp, err
	}

	t.logger.Printf("← [%s] %s - Status: %d - Duration: %v", req.Method, req.URL.Redacted(), resp.StatusCode, duration)

	if isVerbose {
		t.logger.Printf("  Response Headers: %v", headersToString(resp.Header))
		if resp.Body != nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body = io.NopCloser(bytes.NewBuffer(body))
			t.logger.Printf("  Response Body: %s", string(body))
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
