package inpu

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httptest"
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

	client := New().
		TlsConfig(overidedTlsConfig).
		TlsConfig(expectedTlsConfig)

	client.Get("/")

	c.Require().Equal(client.tlsConfig, expectedTlsConfig)
	c.Require().Equal(client.userClient.Transport.(*http.Transport).TLSClientConfig, expectedTlsConfig)
}

func (c *ClientSuite) Test_Client_No_Redirect() {
	serverResponse := "got message from redirected server"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(serverResponse))
	}))

	defer server.Close()

	// Redirect server
	redirectServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, server.URL, http.StatusFound) // 302
	}))
	defer redirectServer.Close()

	redirectionError := errors.New("server has redirected")

	err := New().
		DisableRedirects().
		Get(redirectServer.URL).
		OnReply(StatusIsRedirection, ReturnError(redirectionError)).
		OnReply(StatusAny, ReturnDefaultError).
		Send()

	c.Require().ErrorIs(err, redirectionError)

	err = New().
		FollowRedirects(0). // should have same effect
		Get(redirectServer.URL).
		OnReply(StatusIsRedirection, ReturnError(redirectionError)).
		OnReply(StatusAny, ReturnDefaultError).
		Send()

	c.Require().ErrorIs(err, redirectionError)
}

func (c *ClientSuite) Test_Client_Redirect_Max_Count() {
	serverResponse := "got message from redirected server"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(serverResponse))
	}))

	defer server.Close()

	// Redirect server
	redirectServer1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, server.URL, http.StatusFound) // 302
	}))
	defer redirectServer1.Close()

	// Redirect server
	redirectServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectServer1.URL, http.StatusFound) // 302
	}))
	defer redirectServer.Close()

	redirectionError := errors.New("server has redirected")

	err := New().
		FollowRedirects(3).
		Get(redirectServer.URL).
		OnReply(StatusIsOk, DoNothing).
		OnReply(StatusIsRedirection, ReturnError(redirectionError)).
		OnReply(StatusAny, ReturnDefaultError).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Tls_verify_insecure() {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the protocol version
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	// c.Require().Equal("HTTP/1.1", r.Proto)
	// c.Require().EqualValues(1, r.ProtoMajor)
	// c.Require().EqualValues(1, r.ProtoMinor)
	err := New().
		DisableTLSVerification().
		Get(server.URL).
		OnReply(StatusAny, ReturnDefaultError).
		OnReply(StatusIsOk, DoNothing).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Tls_Get_Certificate_Error() {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	err := New().
		Get(server.URL).
		OnReply(StatusAny, ReturnDefaultError).
		OnReply(StatusIsOk, DoNothing).
		Send()

	c.Require().ErrorIs(err, ErrConnectionFailed)
	target := &tls.CertificateVerificationError{}
	c.Require().ErrorAs(err, &target)
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
