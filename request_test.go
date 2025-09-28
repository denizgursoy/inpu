package inpu

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	url = "https://x.com"
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
	gock.New(url).
		Get("/").
		MatchHeader(HeaderContentType, MimeTypeJson).
		MatchHeader(HeaderAccepts, MimeTypeJson).
		Reply(http.StatusOK)

	response, err := Get(url).
		AcceptJson().
		ContentTypeJson().
		Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}

func (e *ClientSuite) Test_Basic_Authentication() {
	gock.New(url).
		Get("/").
		BasicAuth(TestUserName, TestUserPassword).
		Reply(http.StatusOK)

	response, err := Get(url).
		AuthBasic(TestUserName, TestUserPassword).
		Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}

func (e *ClientSuite) Test_Token_Authentication() {
	token := "sdsds"
	gock.New(url).
		Get("/").
		MatchHeader(HeaderAuthorization, "Bearer "+token).
		Reply(http.StatusOK)

	response, err := Get(url).
		AuthToken(token).
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
