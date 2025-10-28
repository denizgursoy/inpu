package inpu

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"
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

	client := New().Use(RetryMiddleware(2))

	err := client.Get(server.URL).
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
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

	client := New().Use(RetryMiddleware(2))

	err := client.Get(server.URL).
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
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

	client := New().Use(RetryMiddleware(2))

	err := client.Get(server.URL).
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
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

			client := New().Use(RetryMiddleware(2))

			err := client.Get(server.URL).
				OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
				Send()

			c.Require().Error(err)
			c.Require().Equal(1, count)
		})
	}
}

func (c *ClientSuite) Test_UnsuccessfulRetryError() {
	c.T().Parallel()
	c.T().Log("should not panic because of nil response")
	client := New().Use(RetryMiddleware(2))

	err := client.Get("http://127.0.0.1:7777").
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
		Send()

	var urlError *url.Error
	c.Require().ErrorAs(err, &urlError)
	c.Require().ErrorIs(err, ErrConnectionFailed)
}

func (c *ClientSuite) Test_Wait_By_Retry_After_Value() {
	c.T().Parallel()
	count := 0
	var startTime time.Time
	var duration time.Duration
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		count++
		if count == 1 {
			startTime = time.Now()
			w.Header().Add(HeaderRetryAfter, "2")
			w.WriteHeader(http.StatusTooManyRequests)
		} else if count == 2 {
			duration = time.Since(startTime)
			w.WriteHeader(http.StatusOK)
		}
		return
	}))

	defer server.Close()

	client := New().Use(RetryMiddleware(2))

	err := client.
		Get(server.URL).
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
		Send()

	c.T().Logf("duration is %s", duration)
	c.Require().NoError(err)
	c.Require().Equal(2, count)
	c.Require().InDelta(duration, 2*time.Second, float64(time.Millisecond)*100)
}

func (c *ClientSuite) Test_Custom_Retry_Function() {
	c.T().Parallel()
	count := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		count++
		w.WriteHeader(http.StatusOK)

		return
	}))

	defer server.Close()

	client := New().Use(RetryMiddlewareWithConfig(RetryConfig{
		MaxRetries:        2,
		InitialBackoff:    500 * time.Millisecond,
		MaxBackoff:        30 * time.Second,
		BackoffMultiplier: 2,
		CustomRetryChecker: func(resp *http.Response, err error) bool {
			return resp.StatusCode == http.StatusOK
		},
	}))

	err := client.
		Get(server.URL).
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
	c.Require().Equal(3, count)
}
