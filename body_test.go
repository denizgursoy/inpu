package inpu

import (
	"errors"
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
	err := Post(testUrl, BodyFormData(data)).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	e.Require().NoError(err)
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
	err := Post(testUrl, BodyFormDataFromMap(data)).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	e.Require().NoError(err)
}

func (e *ClientSuite) Test_Body_String() {
	gock.New(testUrl).
		Post("/").
		BodyString("^foo$").
		Reply(http.StatusOK)

	err := Post(testUrl, BodyString("foo")).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	e.Require().NoError(err)
}

func (e *ClientSuite) Test_Body_Xml_Marshal() {
	gock.New(testUrl).
		Post("/").
		BodyString(testDataAsXml).
		Reply(http.StatusOK)

	err := Post(testUrl, BodyXml(testData)).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	e.Require().NoError(err)
}
