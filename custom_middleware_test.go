package inpu

import (
	"errors"
	"net/http"

	"github.com/h2non/gock"
)

func (c *ClientSuite) Test_RequestModifierMiddleware() {
	httpClient := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(httpClient)

	gock.New(testUrl).
		Get("/").
		MatchHeader("foo", "bar").
		MatchHeader("foo3", "bar3").
		MatchParam("foo1", "bar1").
		MatchParam("foo2", "bar2").
		Reply(http.StatusOK)

	client := NewWithHttpClient(httpClient).
		UseMiddlewares(RequestModifierMiddleware(func(request *http.Request) {
			request.Header.Add("foo", "bar")
			query := request.URL.Query()
			query.Add("foo1", "bar1")
			request.URL.RawQuery = query.Encode()
		}, "test-middleware", 99))

	err := client.
		Get(testUrl).
		Query("foo2", "bar2").
		Header("foo3", "bar3").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_RequestIDMiddleware() {
	gock.New(testUrl).
		Get("/").
		HeaderPresent(HeaderXRequestID).
		Reply(http.StatusOK)

	err := c.client.
		UseMiddlewares(RequestIDMiddleware()).
		Get(testUrl).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_ErrorHandlerMiddleware() {
	httpError := errors.New("something happened")
	processedError := errors.New("error is processed")
	gock.New(testUrl).
		Get("/").
		ReplyError(httpError)

	err := c.client.
		UseMiddlewares(ErrorHandlerMiddleware(func(err error) error {
			return errors.Join(processedError, httpError)
		})).
		Get(testUrl).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().ErrorIs(err, processedError)
}
