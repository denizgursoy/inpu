package inpu

import (
	"errors"
	"net/http"

	"github.com/h2non/gock"
)

func (e *ClientSuite) Test_RequestModifierMiddleware() {
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

	response, err := client.
		Get(testUrl).
		Query("foo2", "bar2").
		Header("foo3", "bar3").
		Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}

func (e *ClientSuite) Test_IgnoreNilMiddleware() {
	gock.New(testUrl).
		Get("/").
		HeaderPresent(HeaderXRequestID).
		Reply(http.StatusOK)

	client := New().
		UseMiddlewares(nil, RequestIDMiddleware(), nil)
	response, err := client.
		Get(testUrl).
		Send()

	e.Require().NoError(err)
	e.Require().Len(client.mws, 1)
	e.Require().Equal(http.StatusOK, response.Status())
}

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
