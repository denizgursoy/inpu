package inpu

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

type loggingMiddleware struct {
	verbose  bool
	disabled bool
}

// LoggingMiddleware creates a logging middleware
func LoggingMiddleware(verbose, disabled bool) Middleware {
	return &loggingMiddleware{
		verbose:  verbose,
		disabled: disabled,
	}
}

func (t *loggingMiddleware) ID() string {
	return "default-logging-middleware"
}

func (t *loggingMiddleware) Priority() int {
	return 1
}

func (t *loggingMiddleware) Apply(next http.RoundTripper) http.RoundTripper {
	return &loggingTransport{
		next: next,
		mv:   t,
	}
}

type loggingTransport struct {
	next http.RoundTripper
	mv   *loggingMiddleware
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.mv.disabled {
		return t.next.RoundTrip(req)
	}

	start := time.Now()

	// Log request
	log.Printf("→ [%s] %s", req.Method, req.URL.Redacted())

	if t.mv.verbose {
		log.Printf("  Headers: %v", req.Header)
		if req.Body != nil {
			body, _ := io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(body))
			log.Printf("  Body: %s", string(body))
		}
	}

	// Execute request
	resp, err := t.next.RoundTrip(req)
	duration := time.Since(start)

	// Log response
	if err != nil {
		log.Printf("← [%s] %s - ERROR: %v (took %v)", req.Method, req.URL.Redacted(), err, duration)

		return resp, err
	}

	log.Printf("← [%s] %s - Status: %d - Duration: %v", req.Method, req.URL.Redacted(), resp.StatusCode, duration)

	if t.mv.verbose {
		log.Printf("  Response Headers: %v", resp.Header)
		if resp.Body != nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body = io.NopCloser(bytes.NewBuffer(body))
			log.Printf("  Response Body: %s", string(body))
		}
	}

	return resp, nil
}
