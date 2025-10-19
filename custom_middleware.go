package inpu

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type customMiddleware struct {
	requestModifier  RequestModifier
	responseModifier ResponseModifier
	middlewareID     string
	priority         int
	next             http.RoundTripper
}

// CustomMiddleware creates a logging middleware
func newCustomMiddleware(requestModifier RequestModifier, responseModifier ResponseModifier,
	middlewareID string, priority int,
) Middleware {
	return &customMiddleware{
		requestModifier:  requestModifier,
		responseModifier: responseModifier,
		middlewareID:     middlewareID,
		priority:         priority,
	}
}

// RequestModifierMiddleware creates a middleware that allows request to be modified
func RequestModifierMiddleware(modifier RequestModifier, middlewareID string, priority int) Middleware {
	return newCustomMiddleware(modifier, nil, middlewareID, priority)
}

// ResponseModifierMiddleware creates a middleware that allows request to be modified
func ResponseModifierMiddleware(modifier ResponseModifier, middlewareID string, priority int) Middleware {
	return newCustomMiddleware(nil, modifier, middlewareID, priority)
}

// RequestIDMiddleware add HeaderXRequestID to every request
func RequestIDMiddleware() Middleware {
	return RequestModifierMiddleware(func(req *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(req.Context(), "request_id", requestID)

		req = req.WithContext(ctx)
		req.Header.Set(HeaderXRequestID, requestID)
	}, "request-modifier-middleware", 100)
}

// ErrorHandlerMiddleware handles server errors
func ErrorHandlerMiddleware(handler ErrorHandler) Middleware {
	return newCustomMiddleware(nil, func(response *http.Response, serverError error) (*http.Response, error) {
		if serverError != nil {
			return response, handler(serverError)
		}

		return response, nil
	}, "error-handling-middleware", 3)
}

func (t *customMiddleware) ID() string {
	return t.middlewareID
}

func (t *customMiddleware) Priority() int {
	return t.priority
}

func (t *customMiddleware) Apply(next http.RoundTripper) http.RoundTripper {
	t.next = next

	return t
}

func (t *customMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.requestModifier != nil {
		t.requestModifier(req)
	}

	// Execute request
	resp, err := t.next.RoundTrip(req)
	// modify response
	if t.responseModifier != nil {
		return t.responseModifier(resp, err)
	}

	return resp, err
}
