package inpu

import (
	"net/http"

	"github.com/h2non/gock"
)

func (e *ClientSuite) Test_Body_String() {
	gock.New(testUrl).
		Post("/").
		BodyString("^foo$").
		Reply(http.StatusOK)

	response, err := Post(testUrl, BodyString("foo")).
		Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}
