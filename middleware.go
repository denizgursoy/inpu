package inpu

import "net/http"

type Middleware func(http.RoundTripper) http.RoundTripper

type (
	RequestModifier func(*http.Request)
	ErrorHandler    func(error) error
)
