package inpu

import "net/http"

type Middleware interface {
	ID() string
	Priority() int
	Apply(next http.RoundTripper) http.RoundTripper
}

type (
	RequestModifier  func(request *http.Request) (*http.Request, error)
	ResponseModifier func(response *http.Response, server error) (*http.Response, error)
	ErrorHandler     func(serverError error) error
)
