package inpu

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	testUrl = "https://x.com"
)

var (
	TestUserName     = "test-user"
	TestUserPassword = "test-password"
)

type ClientSuite struct {
	suite.Suite
	controller *gomock.Controller
}

func (e *ClientSuite) Test_Headers() {
	gock.New(testUrl).
		Get("/").
		MatchHeader(HeaderContentType, MimeTypeJson).
		MatchHeader(HeaderAccepts, MimeTypeJson).
		Reply(http.StatusOK)

	response, err := Get(testUrl).
		AcceptJson().
		ContentTypeJson().
		Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}

func (e *ClientSuite) Test_Basic_Authentication() {
	gock.New(testUrl).
		Get("/").
		BasicAuth(TestUserName, TestUserPassword).
		Reply(http.StatusOK)

	response, err := Get(testUrl).
		AuthBasic(TestUserName, TestUserPassword).
		Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}

func (e *ClientSuite) Test_Token_Authentication() {
	token := "sdsds"
	gock.New(testUrl).
		Get("/").
		MatchHeader(HeaderAuthorization, "Bearer "+token).
		Reply(http.StatusOK)

	response, err := Get(testUrl).
		AuthToken(token).
		Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
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

	response, err := Get(testUrl).
		QueryBool("is_created", true).
		QueryString("foo", "bar").
		QueryFloat32("float", 1.2).
		QueryFloat64("float64", 2.2).
		QueryInt("int", 1).
		Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}

func (e *ClientSuite) Test_Multiple_Query_Parameters() {
	// TODO test is wrong
	gock.New(testUrl).
		Get("/").
		MatchParam("foo", "bar1").
		MatchParam("foo", "bar2").
		Reply(http.StatusOK)

	response, err := Get(testUrl).
		QueryString("foo", "bar1").
		QueryString("foo", "bar2").
		Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}

func TestClientService(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

func (e *ClientSuite) TearDownTest() {
	gock.Off()
	gock.RestoreClient(http.DefaultClient)
}
