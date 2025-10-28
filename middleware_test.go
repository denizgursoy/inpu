package inpu

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"slices"
)

func (c *ClientSuite) Test_Client_No_Duplicate_Middleware() {
	secondActiveMiddleWare := LoggingMiddleware(true, true)
	client := New().
		UseMiddleware(
			LoggingMiddleware(true, false),
			secondActiveMiddleWare,
		)

	c.Require().Len(client.mws, 1)
	index := slices.IndexFunc(client.mws, func(m Middleware) bool {
		return m.ID() == secondActiveMiddleWare.ID()
	})
	c.Require().NotEqual(-1, index)
	c.Require().Equal(client.mws[index].(*loggingMiddleware).disabled, true)
}

func (c *ClientSuite) Test_Client_MiddlewareOrders() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`test`))
	}))
	defer server.Close()

	loggingMiddleware := LoggingMiddleware(true, false)
	retryMiddleware := RetryMiddleware(3)
	requestIDMiddleware := RequestIDMiddleware()
	client := New().
		BasePath(server.URL).
		UseMiddleware(
			loggingMiddleware,
			retryMiddleware,
			requestIDMiddleware)

	err := client.Post("/", BodyJson(testData)).
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
	c.Equal(client.mws[0], loggingMiddleware)
	c.Equal(client.mws[1], retryMiddleware)
	c.Equal(client.mws[2], requestIDMiddleware)
}

func (c *ClientSuite) Test_IgnoreNilMiddleware() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		c.Require().NotEmpty(request.Header.Get(HeaderXRequestID))
	}))
	defer server.Close()

	client := New().
		UseMiddleware(nil, RequestIDMiddleware(), nil)

	err := client.
		Get(server.URL).
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
	c.Require().Len(client.mws, 1)
}
