package inpu

import (
	"bytes"
	"io"
	"net/http"
	"sync/atomic"

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

func (c *ClientSuite) Test_Response_UnmarshalJson_Body_Close() {
	var closeCalled atomic.Bool
	gock.New(testUrl).Post("/").Reply(http.StatusOK).Body(bytes.NewReader([]byte(`{"foo":"bar"}`)))
	result := testModel{}
	req := NewWithHttpClient(&http.Client{
		Transport: &spyTransport{
			base: http.DefaultTransport,
			onBodyClose: func() {
				closeCalled.Store(true)
			},
		},
	}).Post(testUrl, nil).
		OnReply(StatusIs(http.StatusOK), UnmarshalJson(&result))

	err := req.Send()
	c.Require().NoError(err)
	expectedResponse := testModel{
		Foo: "bar",
	}
	c.Require().Equal(expectedResponse, result)
	c.Require().True(closeCalled.Load())
}

func (c *ClientSuite) Test_Response_No_Nil_Parameter() {
	gock.New(testUrl).Post("/").Reply(http.StatusOK).Body(bytes.NewReader([]byte(`{"foo":"bar"}`)))
	err := Post(testUrl, nil).
		OnReply(StatusIs(http.StatusOK), UnmarshalJson(nil)).
		Send()
	c.Require().ErrorIs(err, ErrMarshalToNil)
}

func (c *ClientSuite) Test_Response_Parameter_Must_Be_Pointer() {
	gock.New(testUrl).Post("/").Reply(http.StatusOK).Body(bytes.NewReader([]byte(`{"foo":"bar"}`)))
	result := testModel{}
	err := Post(testUrl, nil).
		OnReply(StatusIs(http.StatusOK), UnmarshalJson(result)).
		Send()
	c.Require().ErrorIs(err, ErrNotPointerParameter)
}

type spyTransport struct {
	base        http.RoundTripper
	onBodyClose func()
}

func (s *spyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := s.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Wrap the body
	resp.Body = &spyBody{
		ReadCloser: resp.Body,
		onClose:    s.onBodyClose,
	}
	return resp, nil
}

type spyBody struct {
	io.ReadCloser
	onClose func()
}

func (s *spyBody) Close() error {
	if s.onClose != nil {
		s.onClose()
	}
	return s.ReadCloser.Close()
}
