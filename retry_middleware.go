package inpu

import (
	"log"
	"net/http"
	"time"
)

type RetryConfig struct {
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	// Multiplier for exponential backoff
	BackoffMultiplier float64
	// Retry on these status codes
	RetryStatusCodes map[int]bool
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
		RetryStatusCodes: map[int]bool{
			http.StatusTooManyRequests:    true, // 429
			http.StatusServiceUnavailable: true, // 503
			http.StatusGatewayTimeout:     true, // 504
		},
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

	// Network error - retry
	if err != nil {
		return true
	}

	// Check if status code is retryable
	if t.config.RetryStatusCodes[resp.StatusCode] {
		return true
	}

	// Default: retry on 5xx errors
	if resp.StatusCode >= 500 {
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
