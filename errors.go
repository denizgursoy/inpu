package inpu

import "errors"

var (
	ErrRequestCreationFailed = errors.New("could not create the request")
	ErrConnectionFailed      = errors.New("could not send the request")
)
