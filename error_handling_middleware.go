package inpu

import "net/http"

// errorHandlingTransport modifies server errors
type errorHandlingTransport struct {
	next    http.RoundTripper
	handler ErrorHandler
}

// ErrorHandlerMiddleware handles server errors
func ErrorHandlerMiddleware(handler ErrorHandler) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &errorHandlingTransport{
			next:    next,
			handler: handler,
		}
	}
}

func (t *errorHandlingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.next.RoundTrip(req)
	if err != nil {
		return resp, t.handler(err)
	}

	return resp, nil
}
