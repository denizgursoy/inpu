package inpu

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
)

type ClientSuite struct {
	suite.Suite
}

func TestClientService(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

func (e *ClientSuite) TearDownTest() {
	gock.Off()
	gock.RestoreClient(http.DefaultClient)
}
