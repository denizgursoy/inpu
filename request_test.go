package inpu

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ClientSuite struct {
	suite.Suite
	controller *gomock.Controller
}

func (e *ClientSuite) SetupSuite() {

}

type TestMode struct {
}

func (e *ClientSuite) Test_HappyCase() {
	response, err := Get("https://x.com").
		AcceptJson().
		ContentTypeJson().
		QueryBool("triggerered", false).
		QueryString("test", "").
		Header("a", "b").
		AuthBasic("sds", "asdsds").
		FollowRedirect().
		TimeOutIn(time.Second * 10).
		Send()
	if err != nil {
		return
	}

	if response.Is(http.StatusOK) {
		response.ParseJson(nil)
	} else if response.IsOneOf(http.StatusBadRequest, http.StatusBadGateway) {
		response.ParseJson(nil)
	} else if response.IsServerError() {
		response.ParseJson(nil)
	}
}

func (e *ClientSuite) SetupTest() {
	e.controller = gomock.NewController(e.T())
}

func TestClientService(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

func (e *ClientSuite) TearDownTest() {
	e.controller.Finish()
}

func (e *ClientSuite) TearDownSuite() {
}
