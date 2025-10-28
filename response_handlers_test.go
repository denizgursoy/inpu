package inpu

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

func (c *ClientSuite) Test_Response_UnmarshalJson() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"foo":"bar"}`))
	}))
	defer server.Close()

	result := testModel{}
	req := Post(server.URL, nil).
		OnReplyIf(StatusIs(http.StatusOK), ThenUnmarshalJsonTo(&result))

	err := req.Send()
	c.Require().NoError(err)

	expectedResponse := testModel{
		Foo: "bar",
	}
	c.Require().Equal(expectedResponse, result)
}

func (c *ClientSuite) Test_Response_No_Nil_Parameter() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"foo":"bar"}`))
	}))
	defer server.Close()

	err := Post(server.URL, nil).
		OnReplyIf(StatusIs(http.StatusOK), ThenUnmarshalJsonTo(nil)).
		Send()

	c.Require().ErrorIs(err, ErrMarshalToNil)
}

func (c *ClientSuite) Test_Response_Parameter_Must_Be_Pointer() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"foo":"bar"}`))
	}))
	defer server.Close()

	result := testModel{}
	err := Post(server.URL, nil).
		OnReplyIf(StatusIs(http.StatusOK), ThenUnmarshalJsonTo(result)).
		Send()
	c.Require().ErrorIs(err, ErrNotPointerParameter)
}

func (c *ClientSuite) Test_Response_ReturnDefaultError() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"foo":"bar"}`))
	}))
	defer server.Close()

	err := Post(server.URL, nil).
		OnReplyIf(StatusAny, ThenReturnDefaultError).
		Send()

	c.Require().Error(err)
	c.Require().Equal(fmt.Sprintf("called [POST] -> %s and got 500", server.URL), err.Error())
}

func (c *ClientSuite) Test_Response_ReturnError() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"foo":"bar"}`))
	}))
	defer server.Close()

	expectedError := errors.New("something happened")
	actualError := Post(server.URL, nil).
		OnReplyIf(StatusAny, ThenReturnError(expectedError)).
		Send()
	c.Require().ErrorIs(actualError, expectedError)
}

func (c *ClientSuite) Test_Response_Body_Closed() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"foo":"bar"}`))
	}))
	defer server.Close()

	var closedBody *spyBody
	client := New().Use(ResponseModifierMiddleware(func(response *http.Response, server error) (*http.Response, error) {
		// Wrap the body
		closedBody = &spyBody{
			ReadCloser: response.Body,
		}

		response.Body = closedBody
		return response, nil
	}, "test", 12))

	err := client.Post(server.URL, nil).
		OnReplyIf(StatusAny, ThenDoNothing).
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
