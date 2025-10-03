package inpu

import (
	"net/http"

	"github.com/h2non/gock"
)

func (e *ClientSuite) Test_Body_Xml_Marshal() {
	gock.New(testUrl).
		Post("/").
		BodyString(testDataAsXml).
		Reply(http.StatusOK)

	response, err := Post(testUrl, BodyXml(testData)).
		Send()

	e.Require().NoError(err)
	e.Require().Equal(http.StatusOK, response.Status())
}
