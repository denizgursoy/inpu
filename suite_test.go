package inpu

import (
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
)

type ClientSuite struct {
	suite.Suite
}

func TestClientService(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

func (c *ClientSuite) SetupTest() {
	gock.Observe(gock.DumpRequest)
}

func (c *ClientSuite) TearDownTest() {
	gock.Off()
	gock.RestoreClient(http.DefaultClient)
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
