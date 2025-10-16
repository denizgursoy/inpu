package inpu

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testModel struct {
	Foo string `json:"foo" xml:"foo"`
}

var (
	TestUserName     = "test-user"
	TestUserPassword = "test-password"
	testData         = testModel{Foo: "bar"}
	testDataAsJson   = `{"foo":"bar"}`
	testDataAsXml    = `<testModel><foo>bar</foo></testModel>`
)

func (c *ClientSuite) Test_Headers() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		c.Require().Equal(MimeTypeJson, request.Header.Get(HeaderAccept))
		c.Require().Equal(MimeTypeJson, request.Header.Get(HeaderContentType))
	}))
	defer server.Close()

	err := Get(server.URL).
		AcceptJson().
		ContentTypeJson().
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Basic_Authentication() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		c.Require().Equal("Basic dGVzdC11c2VyOnRlc3QtcGFzc3dvcmQ=", request.Header.Get(HeaderAuthorization))
	}))
	defer server.Close()

	err := Get(server.URL).
		AuthBasic(TestUserName, TestUserPassword).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Token_Authentication() {
	token := "sdsds"
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		c.Require().Equal("Bearer "+token, request.Header.Get(HeaderAuthorization))
	}))
	defer server.Close()

	err := Get(server.URL).
		AuthToken(token).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Query_Parameters() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		c.Require().Equal("/?float=1.2000000476837158&float64=2.2&foo=bar+test+encoded&int=1&is_created=true", request.RequestURI)
	}))
	defer server.Close()

	err := Get(server.URL).
		QueryBool("is_created", true).
		Query("foo", "bar test encoded").
		QueryFloat32("float", 1.2).
		QueryFloat64("float64", 2.2).
		QueryInt("int", 1).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Multiple_Query_Parameters() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		c.Require().Equal("/?foo=bar1&foo=bar2", request.RequestURI)
	}))
	defer server.Close()

	err := Get(server.URL).
		Query("foo", "bar1").
		Query("foo", "bar2").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Body_Json_Marshal() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		all, err := io.ReadAll(request.Body)
		c.Require().NoError(err)
		c.Require().Equal(testDataAsJson, string(all))
	}))
	defer server.Close()

	err := Post(server.URL, BodyJson(testData)).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Body_Reader() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		all, err := io.ReadAll(request.Body)
		c.Require().NoError(err)
		c.Require().Equal(testDataAsJson, string(all))
	}))
	defer server.Close()

	err := Post(server.URL, BodyReader(bytes.NewReader([]byte(testDataAsJson)))).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Multiple_Chose_Correct_Reply_Behaviour() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	expectedError := errors.New("correct reply was executed")

	c.T().Run("should select StatusIs over StatusIsSuccess", func(t *testing.T) {
		err := Get(server.URL).
			Query("foo", "bar1").
			Query("foo", "bar2").
			OnReply(StatusIsSuccess, ReturnError(errors.New("unexpected status"))).
			OnReply(StatusIs(http.StatusOK), ReturnError(expectedError)).
			Send()

		c.Require().ErrorIs(err, expectedError)
	})
	c.T().Run("should select StatusIs over StatusIsOneOf", func(t *testing.T) {
		err := Get(server.URL).
			Query("foo", "bar1").
			Query("foo", "bar2").
			OnReply(StatusIsOneOf(http.StatusOK, http.StatusAccepted), ReturnError(errors.New("unexpected status"))).
			OnReply(StatusIs(http.StatusOK), ReturnError(expectedError)).
			Send()

		c.Require().ErrorIs(err, expectedError)
	})
	c.T().Run("should select StatusIsOneOf over StatusIsSuccess", func(t *testing.T) {
		err := Get(server.URL).
			Query("foo", "bar1").
			Query("foo", "bar2").
			OnReply(StatusIsSuccess, ReturnError(errors.New("unexpected status"))).
			OnReply(StatusIsOneOf(http.StatusOK, http.StatusAccepted), ReturnError(expectedError)).
			Send()

		c.Require().ErrorIs(err, expectedError)
	})
	c.T().Run("should select StatusIsOneOf over StatusAny ", func(t *testing.T) {
		err := Get(server.URL).
			Query("foo", "bar1").
			Query("foo", "bar2").
			OnReply(StatusAny, ReturnError(errors.New("unexpected status"))).
			OnReply(StatusIsOneOf(http.StatusOK, http.StatusAccepted), ReturnError(expectedError)).
			Send()

		c.Require().ErrorIs(err, expectedError)
	})
}
