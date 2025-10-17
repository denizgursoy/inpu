package inpu

import (
	"errors"
	"net/http"
	"net/http/httptest"
)

func (c *ClientSuite) Test_RequestModifierMiddleware() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)

		url := request.URL.Query()
		c.Require().EqualValues(url.Get("foo1"), "bar1")
		c.Require().EqualValues(url.Get("foo2"), "bar2")
		c.Require().EqualValues(request.Header.Get("foo3"), "bar3")
		c.Require().EqualValues(request.Header.Get("foo"), "bar")
	}))

	defer server.Close()

	client := New().
		UseMiddlewares(RequestModifierMiddleware(func(request *http.Request) {
			request.Header.Add("foo", "bar")
			query := request.URL.Query()
			query.Add("foo1", "bar1")
			request.URL.RawQuery = query.Encode()
		}, "test-middleware", 99))

	err := client.
		Get(server.URL).
		QueryString("foo2", "bar2").
		Header("foo3", "bar3").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_RequestIDMiddleware() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)

		c.Require().NotEmpty(request.Header.Get(HeaderXRequestID))
	}))

	defer server.Close()

	err := New().
		UseMiddlewares(RequestIDMiddleware()).
		Get(server.URL).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_ErrorHandlerMiddleware() {
	httpError := errors.New("something happened")
	processedError := errors.New("error is processed")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		// Hijack connection and close it immediately
		hj, ok := w.(http.Hijacker)
		if !ok {
			c.Require().Fail("could not hijack")
		}

		conn, _, err := hj.Hijack()
		if err != nil {
			c.Require().NoError(err)
		}
		conn.Close() // Close connection abruptly
	}))
	defer server.Close()

	err := New().
		UseMiddlewares(ErrorHandlerMiddleware(func(err error) error {
			return errors.Join(processedError, httpError)
		})).
		Get(server.URL).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().ErrorIs(err, processedError)
}
