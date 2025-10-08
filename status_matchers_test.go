package inpu

import (
	"net/http"
)

func (c *ClientSuite) Test_Response_IsSuccess() {
	c.Require().False(StatusIsSuccess.Match(http.StatusEarlyHints))
	c.Require().True(StatusIsSuccess.Match(http.StatusOK))
	c.Require().True(StatusIsSuccess.Match(http.StatusCreated))
	c.Require().False(StatusIsSuccess.Match(http.StatusMultipleChoices))
}

func (c *ClientSuite) Test_Response_IsServerError() {
	c.Require().False(StatusIsServerError.Match(http.StatusUnavailableForLegalReasons))
	c.Require().True(StatusIsServerError.Match(http.StatusInternalServerError))
}

func (c *ClientSuite) Test_Response_IsClientError() {
	c.Require().False(StatusIsClientError.Match(http.StatusPermanentRedirect))
	c.Require().True(StatusIsClientError.Match(http.StatusBadRequest))
	c.Require().True(StatusIsClientError.Match(http.StatusUnavailableForLegalReasons))
	c.Require().False(StatusIsClientError.Match(http.StatusInternalServerError))
}

func (c *ClientSuite) Test_Response_IsRedirection() {
	c.Require().True(StatusIsRedirection.Match(http.StatusMultipleChoices))
	c.Require().False(StatusIsRedirection.Match(http.StatusBadRequest))
}

func (c *ClientSuite) Test_Response_IsInformational() {
	c.Require().True(StatusIsInformational.Match(http.StatusContinue))
	c.Require().False(StatusIsInformational.Match(http.StatusOK))
}
