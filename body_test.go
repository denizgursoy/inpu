package inpu

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
)

func (c *ClientSuite) Test_Body_BodyFormDataFromUrl() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		all, err := io.ReadAll(r.Body)
		c.Require().NoError(err)
		c.Require().Equal("email=user%40example.com&email=user2%40example.com&foo=bar&foo1=bar1", string(all))
	}))
	defer server.Close()

	data := map[string][]string{
		"email": {"user@example.com", "user2@example.com"},
		"foo":   {"bar"},
		"foo1":  {"bar1"},
	}
	err := Post(server.URL, BodyFormData(data)).
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Body_BodyFormDataFromMap() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		all, err := io.ReadAll(r.Body)
		c.Require().NoError(err)
		c.Require().Equal("email=user%40example.com&foo=bar&foo1=bar1", string(all))
	}))
	defer server.Close()

	data := map[string]string{
		"email": "user@example.com",
		"foo":   "bar",
		"foo1":  "bar1",
	}
	err := Post(server.URL, BodyFormDataFromMap(data)).
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Body_String() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		all, err := io.ReadAll(r.Body)
		c.Require().NoError(err)
		c.Require().Equal("foo", string(all))
	}))
	defer server.Close()

	err := Post(server.URL, BodyString("foo")).
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_Body_Xml_Marshal() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		all, err := io.ReadAll(r.Body)
		c.Require().NoError(err)
		c.Require().Equal(testDataAsXml, string(all))
	}))
	defer server.Close()

	err := Post(server.URL, BodyXml(testData)).
		OnReplyIf(StatusAnyExcept(http.StatusOK), ThenReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
}
