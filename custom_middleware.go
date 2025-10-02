package inpu

import (
	"net/http"
)

// requestModifierTransport modifies the request only
type requestModifierTransport struct {
	next     http.RoundTripper
	modifier RequestModifier
}

// RequestModifierMiddleware creates a middleware that allows request to be modified
func RequestModifierMiddleware(modifier RequestModifier) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &requestModifierTransport{
			next:     next,
			modifier: modifier,
		}
	}
}

func (t *requestModifierTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// modify request
	t.modifier(req)
	// Execute request
	return t.next.RoundTrip(req)
}
