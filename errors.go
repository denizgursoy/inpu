package inpu

import "errors"

var (
	ErrRequestCreationFailed = errors.New("could not create the request")
	ErrConnectionFailed      = errors.New("could not send the request")
	ErrMarshalingFailed      = errors.New("could not marshall the body")
)
