package inpu

import (
	"errors"
	"net/http"
	"net/http/httptest"
)

func (c *ClientSuite) Test_OnOk_Shorthand() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	called := false
	err := New().Get(server.URL).
		OnOk(func(_ *http.Response) error {
			called = true
			return nil
		}).
		Send()

	c.Require().NoError(err)
	c.Require().True(called)
}

func (c *ClientSuite) Test_OnNotFound_Shorthand() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	expectedErr := errors.New("not found")
	err := New().Get(server.URL).
		OnNotFound(ThenReturnError(expectedErr)).
		Send()

	c.Require().ErrorIs(err, expectedErr)
}

func (c *ClientSuite) Test_OnAny_Shorthand() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	defer server.Close()

	called := false
	err := New().Get(server.URL).
		OnAny(func(_ *http.Response) error {
			called = true
			return nil
		}).
		Send()

	c.Require().NoError(err)
	c.Require().True(called)
}

func (c *ClientSuite) Test_OnAnyExcept_Shorthand() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	expectedErr := errors.New("unexpected status")
	err := New().Get(server.URL).
		OnAnyExcept(http.StatusOK, ThenReturnError(expectedErr)).
		Send()

	c.Require().ErrorIs(err, expectedErr)
}

func (c *ClientSuite) Test_OnAnyExcept_Shorthand_NoMatch() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := New().Get(server.URL).
		OnAnyExcept(http.StatusOK, ThenReturnError(errors.New("should not happen"))).
		Send()

	c.Require().NoError(err)
}

func (c *ClientSuite) Test_On_StatusCode_Shorthand() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	called := false
	err := New().Get(server.URL).
		On(http.StatusCreated, func(_ *http.Response) error {
			called = true
			return nil
		}).
		OnAny(ThenReturnDefaultError).
		Send()

	c.Require().NoError(err)
	c.Require().True(called)
}

func (c *ClientSuite) Test_OnOneOf_Shorthand() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	called := false
	err := New().Get(server.URL).
		OnOneOf(func(_ *http.Response) error {
			called = true
			return nil
		}, http.StatusOK, http.StatusCreated, http.StatusAccepted).
		OnAny(ThenReturnDefaultError).
		Send()

	c.Require().NoError(err)
	c.Require().True(called)
}

func (c *ClientSuite) Test_OnAnyExceptOneOf_Shorthand() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	expectedErr := errors.New("unexpected")
	err := New().Get(server.URL).
		OnAnyExceptOneOf(ThenReturnError(expectedErr), http.StatusOK, http.StatusCreated).
		Send()

	c.Require().ErrorIs(err, expectedErr)
}

func (c *ClientSuite) Test_OnSuccess_Shorthand() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	called := false
	err := New().Get(server.URL).
		OnSuccess(func(_ *http.Response) error {
			called = true
			return nil
		}).
		Send()

	c.Require().NoError(err)
	c.Require().True(called)
}

func (c *ClientSuite) Test_OnClientError_Shorthand() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	expectedErr := errors.New("client error")
	err := New().Get(server.URL).
		OnClientError(ThenReturnError(expectedErr)).
		Send()

	c.Require().ErrorIs(err, expectedErr)
}

func (c *ClientSuite) Test_OnServerError_Shorthand() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	expectedErr := errors.New("server error")
	err := New().Get(server.URL).
		OnServerError(ThenReturnError(expectedErr)).
		Send()

	c.Require().ErrorIs(err, expectedErr)
}

func (c *ClientSuite) Test_On_Priority_Ok_Over_Any() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	okCalled := false
	anyCalled := false
	err := New().Get(server.URL).
		OnAny(func(_ *http.Response) error {
			anyCalled = true
			return nil
		}).
		OnOk(func(_ *http.Response) error {
			okCalled = true
			return nil
		}).
		Send()

	c.Require().NoError(err)
	c.Require().True(okCalled, "OnOk should have been called")
	c.Require().False(anyCalled, "OnAny should not have been called when OnOk matched")
}

func (c *ClientSuite) Test_OnInternalServerError_Shorthand() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	expectedErr := errors.New("internal server error")
	err := New().Get(server.URL).
		OnInternalServerError(ThenReturnError(expectedErr)).
		Send()

	c.Require().ErrorIs(err, expectedErr)
}

func (c *ClientSuite) Test_OnUnauthorized_Shorthand() {
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	expectedErr := errors.New("unauthorized")
	err := New().Get(server.URL).
		OnUnauthorized(ThenReturnError(expectedErr)).
		Send()

	c.Require().ErrorIs(err, expectedErr)
}
