package inpu

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/h2non/gock"
)

type testModel struct {
	Foo string `json:"foo" xml:"foo"`
}

var testUrl = "https://x.com"

var (
	TestUserName     = "test-user"
	TestUserPassword = "test-password"
	testData         = testModel{Foo: "bar"}
	testDataAsJson   = `{"foo":"bar"}`
	testDataAsXml    = `<testModel><foo>bar</foo></testModel>`
)

func (e *ClientSuite) Test_Headers() {
	gock.New(testUrl).
		Get("/").
		MatchHeader(HeaderContentType, MimeTypeJson).
		MatchHeader(HeaderAccept, MimeTypeJson).
		Reply(http.StatusOK)

	err := Get(testUrl).
		AcceptJson().
		ContentTypeJson().
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	e.Require().NoError(err)
}

func (e *ClientSuite) Test_Basic_Authentication() {
	gock.New(testUrl).
		Get("/").
		BasicAuth(TestUserName, TestUserPassword).
		Reply(http.StatusOK)

	err := Get(testUrl).
		AuthBasic(TestUserName, TestUserPassword).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	e.Require().NoError(err)
}

func (e *ClientSuite) Test_Token_Authentication() {
	token := "sdsds"
	gock.New(testUrl).
		Get("/").
		MatchHeader(HeaderAuthorization, "Bearer "+token).
		Reply(http.StatusOK)

	err := Get(testUrl).
		AuthToken(token).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	e.Require().NoError(err)
}

func (e *ClientSuite) Test_Query_Parameters() {
	gock.New(testUrl).
		Get("/").
		MatchParam("is_created", "true").
		MatchParam("foo", "bar").
		MatchParam("float", "1.2").
		MatchParam("float64", "2.2").
		MatchParam("int", "1").
		Reply(http.StatusOK)

	err := Get(testUrl).
		QueryBool("is_created", true).
		Query("foo", "bar").
		QueryFloat32("float", 1.2).
		QueryFloat64("float64", 2.2).
		QueryInt("int", 1).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	e.Require().NoError(err)
}

func (e *ClientSuite) Test_Multiple_Query_Parameters() {
	// TODO test is wrong
	gock.New(testUrl).
		Get("/").
		MatchParam("foo", "bar1").
		MatchParam("foo", "bar2").
		Reply(http.StatusOK)

	err := Get(testUrl).
		Query("foo", "bar1").
		Query("foo", "bar2").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	e.Require().NoError(err)

}

func (e *ClientSuite) Test_Body_Json_Marshal() {
	gock.New(testUrl).
		Post("/").
		BodyString(testDataAsJson).
		Reply(http.StatusOK)

	err := Post(testUrl, testData).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	e.Require().NoError(err)
}

func (e *ClientSuite) Test_Body_Reader() {
	gock.New(testUrl).
		Post("/").
		BodyString(testDataAsJson).
		Reply(http.StatusOK)

	err := Post(testUrl, bytes.NewReader([]byte(testDataAsJson))).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	e.Require().NoError(err)
}
