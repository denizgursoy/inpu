package inpu

import (
	"errors"
	"net/http"
	"net/http/httptest"
)

func (c *ClientSuite) Test_RetryMiddleware() {
	c.T().Parallel()
	count := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if count < 2 {
			w.WriteHeader(http.StatusInternalServerError)
			count++

			return
		}
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	client := New().
		UseMiddlewares(RetryMiddleware(2))

	err := client.Get(server.URL).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}
