package inpu

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/h2non/gock"
)

func (c *ClientSuite) Test_Client() {
	// should get the headers and queries from the client
	gock.New(testUrl).
		Get("/").
		MatchHeader(HeaderContentType, MimeTypeJson).
		MatchHeader(HeaderAccept, MimeTypeJson).
		MatchParam("is_created", "true").
		MatchParam("foo", "bar").
		MatchParam("float", "1.2").
		MatchParam("float64", "2.2").
		MatchParam("int", "1").
		Reply(http.StatusOK)

	c.client.
		AcceptJson().
		ContentTypeJson().
		QueryBool("is_created", true).
		QueryString("foo", "bar").
		QueryFloat32("float", 1.2).
		QueryFloat64("float64", 2.2).
		QueryInt("int", 1)

	err := c.client.Get(testUrl).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)

	// TODO check the mock after changing the post path
	gock.New(testUrl).
		Post("/").
		MatchHeader(HeaderContentType, MimeTypeJson).
		MatchHeader(HeaderAccept, MimeTypeJson).
		MatchParam("is_created", "true").
		MatchParam("foo", "bar").
		MatchParam("float", "1.2").
		MatchParam("float64", "2.2").
		MatchParam("int", "1").
		BodyString(testDataAsJson).
		Reply(http.StatusCreated)

	err = c.client.Post(testUrl, BodyJson(testData)).
		OnReply(StatusAnyExcept(http.StatusCreated), ReturnError(errors.New("unexpected status"))).
		Send()
	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Client_Timeout() {
	// should get the headers and queries from the client
	gock.New(testUrl).
		Get("/").
		Map(func(req *http.Request) *http.Request {
			time.Sleep(300 * time.Millisecond)
			return req
		}).
		Reply(http.StatusOK)

	c.client.TimeOutIn(200 * time.Millisecond)

	err := c.client.Get(testUrl).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().ErrorIs(err, context.DeadlineExceeded)
}

func (c *ClientSuite) Test_Client_BasePath() {
	// should get the headers and queries from the client
	gock.New(testUrl).
		Get("^/people/1$").
		MatchParam("is_created", "^true$").
		MatchParam("foo", "^bar$").
		Reply(http.StatusOK)

	c.client.
		BasePath(testUrl).
		QueryBool("is_created", true)

	err := c.client.
		Get("/people/1").
		Query("foo", "bar").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Client_Empty_BasePath() {
	// should get the headers and queries from the client
	gock.New("").
		Get("^/people/1$").
		MatchParam("is_created", "^true$").
		MatchParam("foo", "^bar$").
		Reply(http.StatusOK)

	c.client.
		QueryBool("is_created", true)

	err := c.client.
		Get("/people/1").
		Query("foo", "bar").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Client_Empty_Uri() {
	// should get the headers and queries from the client
	gock.New(testUrl).
		Get("").
		MatchParam("is_created", "^true$").
		MatchParam("foo", "^bar$").
		Reply(http.StatusOK)

	c.client.
		BasePath(testUrl).
		QueryBool("is_created", true)

	err := c.client.
		Get("").
		Query("foo", "bar").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Client_No_Duplicate_Slash() {
	// should get the headers and queries from the client
	gock.New(testUrl).
		Get("/").
		MatchParam("is_created", "^true$").
		MatchParam("foo", "^bar$").
		Reply(http.StatusOK)

	client := c.client.
		BasePath(testUrl+"/").
		QueryBool("is_created", true)

	err := client.
		Get("/").
		Query("foo", "bar").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Client_No_Higher_Path_Than_Host() {
	// should get the headers and queries from the client
	gock.New(testUrl).
		Get("^/test$").
		MatchParam("is_created", "^true$").
		MatchParam("foo", "^bar$").
		Reply(http.StatusOK)

	c.client.
		BasePath(testUrl+"/people/1/subscription/23").
		QueryBool("is_created", true)

	err := c.client.
		Get("/../../../../../../test").
		Query("foo", "bar").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Client_No_Nil_Transport() {
	gock.New(testUrl).Get("/").Reply(http.StatusOK)
	c.client.userClient.Transport = nil
	c.client.Get(testUrl)
	c.Require().NotNil(c.client.userClient.Transport)
}

func (c *ClientSuite) Test_Client_Use_The_Last_Provided_Tls_Config() {
	overidedTlsConfig := &tls.Config{}
	expectedTlsConfig := &tls.Config{
		ServerName: "expected",
	}

	client := NewWithHttpClient(&http.Client{}).
		TlsConfig(overidedTlsConfig).
		TlsConfig(expectedTlsConfig)

	client.Get("/")

	c.Require().Equal(client.tlsConfig, expectedTlsConfig)
	c.Require().Equal(client.userClient.Transport.(*http.Transport).TLSClientConfig, expectedTlsConfig)
}

// Measure allocations
func Benchmark_QueryBuild(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	get := Get("https://jsonplaceholder.typicode.com/todos")
	for i := 0; i < b.N; i++ {
		itoa := strconv.Itoa(i)
		get.Query("foo"+itoa, "% &bar").QueryInt("foo"+itoa, 1)
	}
	b.StopTimer()
}
