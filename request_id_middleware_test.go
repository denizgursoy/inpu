package inpu

import (
	"net/http"

	"github.com/h2non/gock"
)

func (e *ClientSuite) Test_RequestIDMiddleware() {
	gock.New(testUrl).
		Get("/").
		HeaderPresent(HeaderXRequestID).
		Reply(http.StatusOK)

	response, err := New().
		UseMiddlewares(RequestIDMiddleware()).
		Get(testUrl).
		Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}
