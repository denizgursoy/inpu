package inpu

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// RequestIDMiddleware add HeaderXRequestID to every request
func RequestIDMiddleware() Middleware {
	return RequestModifierMiddleware(func(req *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(req.Context(), "request_id", requestID)

		req = req.WithContext(ctx)
		req.Header.Set(HeaderXRequestID, requestID)
	})
}
