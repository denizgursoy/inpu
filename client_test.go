package inpu

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/h2non/gock"
)

func (e *ClientSuite) Test_Client() {
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

	client := New().
		AcceptJson().
		ContentTypeJson().
		QueryBool("is_created", true).
		QueryString("foo", "bar").
		QueryFloat32("float", 1.2).
		QueryFloat64("float64", 2.2).
		QueryInt("int", 1)

	err := client.Get(testUrl).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	e.Require().NoError(err)

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

	err = client.Post(testUrl, testData).
		OnReply(StatusAnyExcept(http.StatusCreated), ReturnError(errors.New("unexpected status"))).
		Send()
	e.Require().NoError(err)
}

func (e *ClientSuite) Test_Client_Timeout() {
	// should get the headers and queries from the client
	gock.New(testUrl).
		Get("/").
		Map(func(req *http.Request) *http.Request {
			time.Sleep(300 * time.Millisecond)
			return req
		}).
		Reply(http.StatusOK)

	client := New().TimeOutIn(200 * time.Millisecond)

	err := client.Get(testUrl).
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	e.Require().ErrorIs(err, context.DeadlineExceeded)
}

func (e *ClientSuite) Test_Client_BasePath() {
	// should get the headers and queries from the client
	gock.New(testUrl).
		Get("^/people/1$").
		MatchParam("is_created", "^true$").
		MatchParam("foo", "^bar$").
		Reply(http.StatusOK)

	client := New().
		BasePath(testUrl).
		QueryBool("is_created", true)

	err := client.
		Get("/people/1").
		Query("foo", "bar").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	e.Require().NoError(err)
}

func (e *ClientSuite) Test_Client_Empty_BasePath() {
	// should get the headers and queries from the client
	gock.New("").
		Get("^/people/1$").
		MatchParam("is_created", "^true$").
		MatchParam("foo", "^bar$").
		Reply(http.StatusOK)

	client := New().
		QueryBool("is_created", true)

	err := client.
		Get("/people/1").
		Query("foo", "bar").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	e.Require().NoError(err)
}

func (e *ClientSuite) Test_Client_Empty_Uri() {
	// should get the headers and queries from the client
	gock.New(testUrl).
		Get("").
		MatchParam("is_created", "^true$").
		MatchParam("foo", "^bar$").
		Reply(http.StatusOK)

	client := New().
		BasePath(testUrl).
		QueryBool("is_created", true)

	err := client.
		Get("").
		Query("foo", "bar").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	e.Require().NoError(err)
}

func (e *ClientSuite) Test_Client_No_Duplicate_Slash() {
	// should get the headers and queries from the client
	gock.New(testUrl).
		Get("/").
		MatchParam("is_created", "^true$").
		MatchParam("foo", "^bar$").
		Reply(http.StatusOK)

	client := New().
		BasePath(testUrl+"/").
		QueryBool("is_created", true)

	err := client.
		Get("/").
		Query("foo", "bar").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	e.Require().NoError(err)
}

func (e *ClientSuite) Test_Client_No_Higher_Path_Than_Host() {
	// should get the headers and queries from the client
	gock.New(testUrl).
		Get("^/test$").
		MatchParam("is_created", "^true$").
		MatchParam("foo", "^bar$").
		Reply(http.StatusOK)

	client := New().
		BasePath(testUrl+"/people/1/subscription/23").
		QueryBool("is_created", true)

	err := client.
		Get("/../../../../../../test").
		Query("foo", "bar").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	e.Require().NoError(err)
}
