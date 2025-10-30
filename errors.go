package inpu

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrRequestCreationFailed = errors.New("could not create the request")
	ErrInvalidBody           = errors.New("could not create the body")
	ErrConnectionFailed      = errors.New("connection failed")
	ErrCouldNotParseBaseUrl  = errors.New("invalid base path")
	ErrCouldNotParsePath     = errors.New("invalid path")
	ErrMarshalToNil          = errors.New("cannot unmarshal to nil")
	ErrNotPointerParameter   = errors.New("cannot marshal to non pointer type")
	ErrPanickedDuringTheCall = errors.New("panicked on send")
)

type DefaultError struct {
	res *http.Response
}

func (d *DefaultError) Error() string {
	return fmt.Sprintf("called [%s] -> %s and got %d",
		d.res.Request.Method, d.res.Request.URL.Redacted(), d.res.StatusCode)
}
