package inpu

import (
	"net/http"

	"github.com/h2non/gock"
)

func (e *ClientSuite) Test_Response_IsSuccess() {
	gock.New(testUrl).Post("/").Reply(http.StatusOK)
	response, err := Post(testUrl, nil).Send()
	e.Require().NoError(err)
	e.Require().True(response.IsSuccess())
	e.Require().True(response.Is(http.StatusOK))
	e.Require().True(response.Status() == http.StatusOK)
	e.Require().True(response.IsOneOf(http.StatusMovedPermanently, http.StatusOK))
}

func (e *ClientSuite) Test_Response_IsServerError() {
	gock.New(testUrl).Post("/").Reply(http.StatusInternalServerError)
	response, err := Post(testUrl, nil).Send()
	e.Require().NoError(err)
	e.Require().True(response.IsServerError())
	e.Require().True(response.Is(http.StatusInternalServerError))
	e.Require().True(response.Status() == http.StatusInternalServerError)
	e.Require().True(response.IsOneOf(http.StatusMovedPermanently, http.StatusInternalServerError))
}

func (e *ClientSuite) Test_Response_IsClientError() {
	gock.New(testUrl).Post("/").Reply(http.StatusBadRequest)
	response, err := Post(testUrl, nil).Send()
	e.Require().NoError(err)
	e.Require().True(response.IsClientError())
	e.Require().True(response.Is(http.StatusBadRequest))
	e.Require().True(response.Status() == http.StatusBadRequest)
	e.Require().True(response.IsOneOf(http.StatusMovedPermanently, http.StatusBadRequest))
}

func (e *ClientSuite) Test_Response_IsRedirection() {
	gock.New(testUrl).Post("/").Reply(http.StatusMultipleChoices)
	response, err := Post(testUrl, nil).Send()
	e.Require().NoError(err)
	e.Require().True(response.IsRedirection())
	e.Require().True(response.Is(http.StatusMultipleChoices))
	e.Require().True(response.Status() == http.StatusMultipleChoices)
	e.Require().True(response.IsOneOf(http.StatusMovedPermanently, http.StatusMultipleChoices))
}

func (e *ClientSuite) Test_Response_IsInformational() {
	gock.New(testUrl).Post("/").Reply(http.StatusContinue)
	response, err := Post(testUrl, nil).Send()
	e.Require().NoError(err)
	e.Require().True(response.IsInformational())
	e.Require().True(response.Is(http.StatusContinue))
	e.Require().True(response.Status() == http.StatusContinue)
	e.Require().True(response.IsOneOf(http.StatusBadRequest, http.StatusContinue))
}
