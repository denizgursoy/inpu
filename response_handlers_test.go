package inpu

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/h2non/gock"
)

func (c *ClientSuite) Test_Response_UnmarshalJson() {
	gock.New(testUrl).Post("/").Reply(http.StatusOK).Body(bytes.NewReader([]byte(`{"foo":"bar"}`)))
	result := testModel{}
	req := Post(testUrl, nil).
		OnReply(StatusIs(http.StatusOK), UnmarshalJson(&result))

	err := req.Send()
	c.Require().NoError(err)

	expectedResponse := testModel{
		Foo: "bar",
	}
	c.Require().Equal(expectedResponse, result)
}

func (c *ClientSuite) Test_Response_No_Nil_Parameter() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"foo":"bar"}`))
	}))
	defer server.Close()

	err := Post(server.URL, nil).
		OnReply(StatusIs(http.StatusOK), UnmarshalJson(nil)).
		Send()

	c.Require().ErrorIs(err, ErrMarshalToNil)
}

func (c *ClientSuite) Test_Response_Parameter_Must_Be_Pointer() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"foo":"bar"}`))
	}))
	defer server.Close()

	result := testModel{}
	err := Post(server.URL, nil).
		OnReply(StatusIs(http.StatusOK), UnmarshalJson(result)).
		Send()
	c.Require().ErrorIs(err, ErrNotPointerParameter)
}

func (c *ClientSuite) Test_Response_ReturnDefaultError() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"foo":"bar"}`))
	}))
	defer server.Close()

	err := Post(server.URL, nil).
		OnReply(StatusAny, ReturnDefaultError).
		Send()

	c.Require().Error(err)
	c.Require().Equal(fmt.Sprintf("called [POST] -> %s and got 500", server.URL), err.Error())
}

func (c *ClientSuite) Test_Response_ReturnError() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"foo":"bar"}`))
	}))
	defer server.Close()

	expectedError := errors.New("something happened")
	actualError := Post(server.URL, nil).
		OnReply(StatusAny, ReturnError(expectedError)).
		Send()
	c.Require().ErrorIs(actualError, expectedError)
}

func (c *ClientSuite) Test_Response_Body_Closed() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"foo":"bar"}`))
	}))
	defer server.Close()

	var closedBody *spyBody
	client := New().UseMiddlewares(ResponseModifierMiddleware(func(response *http.Response, server error) (*http.Response, error) {
		// Wrap the body
		closedBody = &spyBody{
			ReadCloser: response.Body,
		}

		response.Body = closedBody
		return response, nil
	}, "test", 12))

	err := client.Post(server.URL, nil).
		OnReply(StatusAny, DoNothing).
		Send()

	c.Require().NoError(err)
	c.Require().True(closedBody.isClosed)
}

type spyBody struct {
	io.ReadCloser
	isClosed bool
}

func (s *spyBody) Close() error {
	s.isClosed = true
	return s.ReadCloser.Close()
}
