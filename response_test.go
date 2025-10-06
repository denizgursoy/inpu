package inpu

import (
	"bytes"
	"net/http"

	"github.com/h2non/gock"
)

func (e *ClientSuite) Test_Response_IsSuccess() {
	e.Require().True(StatusIsSuccess(http.StatusOK))
}

func (e *ClientSuite) Test_Response_IsServerError() {
	e.Require().True(StatusIsServerError(http.StatusInternalServerError))
}

func (e *ClientSuite) Test_Response_IsClientError() {
	e.Require().True(StatusIsClientError(http.StatusBadRequest))
}

func (e *ClientSuite) Test_Response_IsRedirection() {
	e.Require().True(StatusIsRedirection(http.StatusMultipleChoices))
}

func (e *ClientSuite) Test_Response_IsInformational() {
	e.Require().True(StatusIsInformational(http.StatusContinue))
	e.Require().True(StatusIsOneOf(http.StatusBadRequest, http.StatusContinue)(http.StatusContinue))
}

func (e *ClientSuite) Test_Response_UnmarshalJson() {
	gock.New(testUrl).Post("/").Reply(http.StatusOK).Body(bytes.NewReader([]byte(`{"foo":"bar"}`)))
	result := testModel{}
	err := Post(testUrl, nil).
		OnReply(StatusIs(http.StatusOK), UnmarshalJson(&result)).
		Send()
	e.Require().NoError(err)

	expectedResponse := testModel{
		Foo: "bar",
	}
	e.Require().Equal(expectedResponse, result)
}

func (e *ClientSuite) Test_Response_No_Nil_Parameter() {
	gock.New(testUrl).Post("/").Reply(http.StatusOK).Body(bytes.NewReader([]byte(`{"foo":"bar"}`)))
	err := Post(testUrl, nil).
		OnReply(StatusIs(http.StatusOK), UnmarshalJson(nil)).
		Send()
	e.Require().ErrorIs(err, ErrMarshalToNil)
}

func (e *ClientSuite) Test_Response_Parameter_Must_Be_Pointer() {
	gock.New(testUrl).Post("/").Reply(http.StatusOK).Body(bytes.NewReader([]byte(`{"foo":"bar"}`)))
	result := testModel{}
	err := Post(testUrl, nil).
		OnReply(StatusIs(http.StatusOK), UnmarshalJson(result)).
		Send()
	e.Require().ErrorIs(err, ErrNotPointerParameter)
}
