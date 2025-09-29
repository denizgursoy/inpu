package inpu

import (
	"net/http"

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

	response, err := client.Get(testUrl).Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())

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

	response, err = client.Post(testUrl, testData).Send()
	e.Require().NoError(err)
	e.Require().Equal(http.StatusCreated, response.Status())
}
