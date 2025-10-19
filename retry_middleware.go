package inpu

import (
	"crypto/tls"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"time"
)

var (
	retriableClientErrors    = []int{http.StatusTooManyRequests}
	nonRetriableServerErrors = []int{
		http.StatusNotImplemented, http.StatusHTTPVersionNotSupported,
		http.StatusLoopDetected, http.StatusVariantAlsoNegotiates,
		http.StatusNetworkAuthenticationRequired,
	}
)

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
	logger := GetLoggerFromContext(req.Context())
	ctx := req.Context()

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
			retryAfterDuration := t.extractBackoffFromHeader(resp)
			timeToWait := favorRetryAfterValueIfNotEmpty(retryAfterDuration, backoff)
			// Check context cancellation
			select {
			case <-req.Context().Done():
				return resp, req.Context().Err()
			case <-time.After(timeToWait):
				logger.Info(ctx, "[RETRY] Attempt %d/%d for %s %s (waiting %v)", attempt+1, t.config.MaxRetries,
					req.Method, req.URL, backoff)
				// drain the body and close the connection because
				// it is going to send another request soon
				err := DrainBodyAndClose(resp.Body)
				if err != nil {
					logger.Error(ctx, err, "could not drain the body")
				}
			}

			// Exponential backoff
			backoff = time.Duration(float64(backoff) * t.config.BackoffMultiplier)
			backoff = t.getMaxBackoffTimeIfBigger(backoff)
		}
	}

	return resp, err
}

func (t *retryMiddleware) getMaxBackoffTimeIfBigger(d time.Duration) time.Duration {
	if d > t.config.MaxBackoff {
		return t.config.MaxBackoff
	}

	return d
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

func favorRetryAfterValueIfNotEmpty(retryAfter time.Duration, backoff time.Duration) time.Duration {
	if retryAfter > 0 {
		return retryAfter
	}

	return backoff
}

func (t *retryMiddleware) extractBackoffFromHeader(response *http.Response) time.Duration {
	if response != nil {
		if response.StatusCode == http.StatusTooManyRequests || response.StatusCode == http.StatusServiceUnavailable {
			if sleep, ok := parseRetryAfterHeader(response.Header.Get(HeaderRetryAfter)); ok {
				return t.getMaxBackoffTimeIfBigger(sleep)
			}
		}
	}

	return 0
}

func parseRetryAfterHeader(header string) (time.Duration, bool) {
	if len(header) == 0 {
		return 0, false
	}
	// Retry-After: 120
	if sleep, err := strconv.ParseInt(header, 10, 64); err == nil {
		if sleep < 0 { // a negative sleep doesn't make sense
			return 0, false
		}
		return time.Second * time.Duration(sleep), true
	}

	// Retry-After: Fri, 31 Dec 1999 23:59:59 GMT
	retryTime, err := http.ParseTime(header)
	if err != nil {
		return 0, false
	}
	if until := retryTime.Sub(time.Now()); until > 0 {
		return until, true
	}

	// date is in the past
	return 0, true
}
