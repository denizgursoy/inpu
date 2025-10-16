package inpu

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
)

type ClientSuite struct {
	suite.Suite
}

func TestClientService(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
