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
)

func (c *ClientSuite) Test_Client() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		// TODO check parameters here
		// should get the headers and queries from the client
		// gock.New(testUrl).
		// 	Get("/").
		// 	MatchHeader(HeaderContentType, MimeTypeJson).
		// 	MatchHeader(HeaderAccept, MimeTypeJson).
		// 	MatchParam("is_created", "true").
		// 	MatchParam("foo", "bar").
		// 	MatchParam("float", "1.2").
		// 	MatchParam("float64", "2.2").
		// 	MatchParam("int", "1").
		// 	Reply(http.StatusOK)
		if request.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusCreated)
		}
	}))
	defer server.Close()

	client := New().
		AcceptJson().
		ContentTypeJson().
		QueryBool("is_created", true).
		QueryString("foo", "bar").
		QueryFloat32("float", 1.2).
		QueryFloat64("float64", 2.2).
		QueryInt("int", 1)

	err := client.Get(server.URL).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)

	err = client.Post(server.URL, BodyJson(testData)).
		OnReply(StatusAnyExcept(http.StatusCreated), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Client_Timeout() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		time.Sleep(300 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := New().
		TimeOutIn(200*time.Millisecond).
		Get(server.URL).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().ErrorIs(err, context.DeadlineExceeded)
}

func (c *ClientSuite) Test_Client_BasePath() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		c.Require().Equal("/people/1", request.RequestURI)
	}))
	defer server.Close()

	err := New().
		BasePath(server.URL).
		Get("/people/1").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Client_Empty_BasePath() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		c.Require().Equal("/people/1?foo=bar&is_created=true", request.RequestURI)
	}))
	defer server.Close()

	err := New().
		BasePath(server.URL).
		QueryBool("is_created", true).
		Get("/people/1").
		Query("foo", "bar").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Client_Empty_Uri() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		c.Require().Equal("/?foo=bar&is_created=true", request.RequestURI)
	}))
	defer server.Close()

	err := New().
		BasePath(server.URL).
		QueryBool("is_created", true).
		Get("").
		Query("foo", "bar").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Client_No_Duplicate_Slash() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		c.Require().Equal("/?foo=bar&is_created=true", request.RequestURI)
	}))
	defer server.Close()

	err := New().
		BasePath(server.URL+"/").
		QueryBool("is_created", true).
		Get("/").
		Query("foo", "bar").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Client_No_Higher_Path_Than_Host() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		c.Require().Equal("/test?foo=bar&is_created=true", request.RequestURI)
	}))
	defer server.Close()

	err := New().
		BasePath(server.URL+"/people/1/subscription/23").
		QueryBool("is_created", true).
		Get("/../../../../../../test").
		Query("foo", "bar").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Client_No_Nil_Transport() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New()
	client.userClient.Transport = nil
	client.Get(server.URL)
	c.Require().NotNil(client.userClient.Transport)
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
	c.T().Parallel()
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
	c.T().Parallel()
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
	c.T().Parallel()
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
	c.T().Parallel()
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	err := New().
		Get(server.URL).
		Send()

	c.Require().ErrorIs(err, ErrConnectionFailed)
	target := &tls.CertificateVerificationError{}
	c.Require().ErrorAs(err, &target)
}

func (c *ClientSuite) Test_Client_Close() {
	c.T().Parallel()

	c.T().Log("will close the client to see requests got cancelled")
	client := New()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client.Close()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	err := client.Get(server.URL).Send()

	c.Require().ErrorIs(err, context.Canceled)
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
