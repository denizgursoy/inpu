package inpu

import (
	"errors"

	"github.com/h2non/gock"
)

func (e *ClientSuite) Test_ErrorHandlerMiddleware() {
	httpError := errors.New("something happened")
	processedError := errors.New("error is processed")
	gock.New(testUrl).
		Get("/").
		ReplyError(httpError)

	response, err := New().
		UseMiddlewares(ErrorHandlerMiddleware(func(err error) error {
			return errors.Join(processedError, httpError)
		})).
		Get(testUrl).
		Send()

	e.Require().ErrorIs(err, processedError)
	e.Require().Nil(response)
}
