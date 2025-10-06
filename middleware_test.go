package inpu

import (
	"errors"
	"net/http"

	"github.com/h2non/gock"
)

func (e *ClientSuite) Test_Client_No_Duplicate_Middleware() {
	secondActiveMiddleWare := LoggingMiddleware(false)
	client := New().
		UseMiddlewares(
			LoggingMiddleware(true),
			secondActiveMiddleWare,
		)

	e.Require().Len(client.mws, 1)
	e.Require().False(client.mws[secondActiveMiddleWare.ID()].(*loggingMiddleware).verbose)
}

func (e *ClientSuite) Test_Client_MiddlewareOrders() {
	// should get the headers and queries from the client
	gock.New(testUrl).
		Get("/").
		Reply(http.StatusOK)

	client := New().
		BasePath(testUrl).
		UseMiddlewares(
			LoggingMiddleware(true),
			RetryMiddleware(3),
			RequestIDMiddleware())

	err := client.Get("/").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	e.Require().NoError(err)
}
