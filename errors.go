package inpu

import (
	"errors"
	"fmt"
)

var (
	ErrRequestCreationFailed = errors.New("could not create the request")
	ErrInvalidBody           = errors.New("could not create the body")
	ErrConnectionFailed      = errors.New("connection failed")
	ErrCouldNotParseBaseUrl  = errors.New("invalid base path")
	ErrCouldNotParsePath     = errors.New("invalid path")
	ErrMarshalToNil          = errors.New("cannot unmarshal to nil")
	ErrNotPointerParameter   = errors.New("cannot marshal to non pointer type ")
)

type DefaultError struct {
	Method     string
	URL        string
	StatusCode int
}

func (d *DefaultError) Error() string {
	return fmt.Sprintf("called [%s] -> %s and got %d", d.Method, d.URL, d.StatusCode)
}
