package inpu

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"slices"
	"time"
)

var retriableClientErrors = []int{http.StatusTooManyRequests}
var nonRetriableServerErrors = []int{
	http.StatusNotImplemented, http.StatusHTTPVersionNotSupported,
	http.StatusLoopDetected, http.StatusVariantAlsoNegotiates,
	http.StatusNetworkAuthenticationRequired}

type RetryConfig struct {
	MaxRetries        int
	InitialBackoff    time.Duration
	MaxBackoff        time.Duration
	BackoffMultiplier float64
}

type retryMiddleware struct {
	config RetryConfig
	next   http.RoundTripper
}

// RetryMiddleware creates a retry middleware with default config
func RetryMiddleware(maxRetries int) Middleware {
	return RetryMiddlewareWithConfig(RetryConfig{
		MaxRetries:        maxRetries,
		InitialBackoff:    500 * time.Millisecond,
		MaxBackoff:        30 * time.Second,
		BackoffMultiplier: 2.0,
	})
}

// RetryMiddlewareWithConfig creates a retry middleware with custom config
func RetryMiddlewareWithConfig(config RetryConfig) Middleware {
	return &retryMiddleware{
		config: config,
	}
}

func (t *retryMiddleware) ID() string {
	return "retry-middleware"
}

func (t *retryMiddleware) Priority() int {
	return 2
}

func (t *retryMiddleware) Apply(next http.RoundTripper) http.RoundTripper {
	t.next = next

	return t
}

func (t *retryMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	backoff := t.config.InitialBackoff

	for attempt := 0; attempt <= t.config.MaxRetries; attempt++ {
		// Clone request for retry (important for body)
		clonedReq := t.cloneRequest(req)
		resp, err = t.next.RoundTrip(clonedReq)

		// Check if we should retry
		if !t.shouldRetry(resp, err, attempt) {
			return resp, err
		}
		// TODO close body
		// Don't sleep after last attempt
		if attempt < t.config.MaxRetries {
			// Check context cancellation
			select {
			case <-req.Context().Done():
				return resp, req.Context().Err()
			case <-time.After(backoff):
				log.Printf("[RETRY] Attempt %d/%d for %s %s (waiting %v)",
					attempt+1, t.config.MaxRetries, req.Method, req.URL, backoff)
			}

			// Exponential backoff
			backoff = time.Duration(float64(backoff) * t.config.BackoffMultiplier)
			if backoff > t.config.MaxBackoff {
				backoff = t.config.MaxBackoff
			}
		}
	}

	return resp, err
}

func (t *retryMiddleware) shouldRetry(resp *http.Response, err error, attempt int) bool {
	// No more retries left
	if attempt >= t.config.MaxRetries {
		return false
	}

	if err != nil {
		return checkRetryBasedOnConnectionError(err)
	}

	return checkRetryBasedOnStatusCode(resp)
}

func checkRetryBasedOnConnectionError(connectionError error) bool {
	var certificateVerificationError *tls.CertificateVerificationError
	if errors.As(connectionError, &certificateVerificationError) {
		return false
	}

	return false
}

func checkRetryBasedOnStatusCode(response *http.Response) bool {
	statusCode := response.StatusCode
	if slices.Contains(retriableClientErrors, statusCode) {
		return true
	} else if slices.Contains(nonRetriableServerErrors, statusCode) {
		return false
	}

	if statusCode >= http.StatusInternalServerError {
		return true
	}

	return false
}

func (t *retryMiddleware) cloneRequest(req *http.Request) *http.Request {
	clonedReq := req.Clone(req.Context())

	// If body exists, we need to cache and restore it
	if req.Body != nil && req.GetBody != nil {
		body, _ := req.GetBody()
		clonedReq.Body = body
	}

	return clonedReq
}
