package inpu

import "errors"

var (
	ErrRequestCreationFailed = errors.New("could not create the request")
	ErrInvalidBody           = errors.New("could not create the body")
	ErrConnectionFailed      = errors.New("could not send the request")
	ErrCouldNotParseBaseUrl  = errors.New("invalid base path")
	ErrCouldNotParsePath     = errors.New("invalid path")
)
