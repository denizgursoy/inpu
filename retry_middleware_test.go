package inpu

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
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

	client := New().UseMiddlewares(RetryMiddleware(2))

	err := client.Get(server.URL).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_No_Retry_On_CertificateVerificationError_Error() {
	c.T().Parallel()
	count := 0
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		count++

		return
	}))

	defer server.Close()

	client := New().UseMiddlewares(RetryMiddleware(2))

	err := client.Get(server.URL).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	var certificateVerificationError *tls.CertificateVerificationError
	c.Require().ErrorAs(err, &certificateVerificationError)
	c.Require().Zero(count)
}

func (c *ClientSuite) Test_Retry_On_429() {
	c.T().Parallel()
	count := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if count < 2 {
			w.WriteHeader(http.StatusTooManyRequests)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		count++

		return
	}))

	defer server.Close()

	client := New().UseMiddlewares(RetryMiddleware(2))

	err := client.Get(server.URL).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
	c.Require().Equal(3, count)
}

func (c *ClientSuite) Test_No_Try_On_The_Non_Retriable_Server_Errors() {
	c.T().Parallel()
	for _, serverError := range nonRetriableServerErrors {
		c.T().Run(strconv.Itoa(serverError), func(t *testing.T) {
			count := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
				w.WriteHeader(serverError)
				count++

				return
			}))

			defer server.Close()

			client := New().UseMiddlewares(RetryMiddleware(2))

			err := client.Get(server.URL).
				OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
				Send()

			c.Require().Error(err)
			c.Require().Equal(1, count)
		})
	}
}

func (c *ClientSuite) Test_UnsuccessfulRetryError() {
	c.T().Parallel()
	c.T().Log("should not panic because of nil response")
	client := New().UseMiddlewares(RetryMiddleware(2))

	err := client.Get("http://127.0.0.1:7777").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	var urlError *url.Error
	c.Require().ErrorAs(err, &urlError)
	c.Require().ErrorIs(err, ErrConnectionFailed)
}
