package inpu

import (
	"errors"
	"net/http"
	"slices"

	"github.com/h2non/gock"
)

func (c *ClientSuite) Test_Client_No_Duplicate_Middleware() {
	secondActiveMiddleWare := LoggingMiddleware(false, false)
	client := New().
		UseMiddlewares(
			LoggingMiddleware(true, false),
			secondActiveMiddleWare,
		)

	c.Require().Len(client.mws, 1)
	index := slices.IndexFunc(client.mws, func(m Middleware) bool {
		return m.ID() == secondActiveMiddleWare.ID()
	})
	c.Require().NotEqual(-1, index)
	c.Require().False(client.mws[index].(*loggingMiddleware).verbose)
}

func (c *ClientSuite) Test_Client_MiddlewareOrders() {
	// should get the headers and queries from the client
	gock.New(testUrl).
		Get("/").
		Reply(http.StatusOK)

	c.client.
		BasePath(testUrl).
		UseMiddlewares(
			LoggingMiddleware(true, false),
			RetryMiddleware(3),
			RequestIDMiddleware())

	err := c.client.Get("/").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}
