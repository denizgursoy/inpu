package inpu

import (
	"net/http"

	"github.com/h2non/gock"
)

func (e *ClientSuite) Test_Body_BodyFormDataFromUrl() {
	gock.New(testUrl).
		Post("/").
		BodyString("^email=user%40example.com&email=user2%40example.com&foo=bar&foo1=bar1$").
		Reply(http.StatusOK)

	data := map[string][]string{
		"email": {"user@example.com", "user2@example.com"},
		"foo":   {"bar"},
		"foo1":  {"bar1"},
	}
	response, err := Post(testUrl, BodyFormData(data)).
		Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}

func (e *ClientSuite) Test_Body_BodyFormDataFromMap() {
	gock.New(testUrl).
		Post("/").
		BodyString("^email=user%40example.com&foo=bar&foo1=bar1$").
		Reply(http.StatusOK)

	data := map[string]string{
		"email": "user@example.com",
		"foo":   "bar",
		"foo1":  "bar1",
	}
	response, err := Post(testUrl, BodyFormDataFromMap(data)).
		Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}
