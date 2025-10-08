package inpu

import (
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
	gock.InterceptClient(getDefaultClient())
}

func (c *ClientSuite) TearDownTest() {
	gock.Off()
	gock.RestoreClient(getDefaultClient())
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
