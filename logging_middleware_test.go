package inpu

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

func (c *ClientSuite) Test_LoggingMiddleware_Level_info() {
	c.T().Log("should log only urls and durations")
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`ok`))
	}))
	defer server.Close()

	lgMiddleware := LoggingMiddleware(false, false)
	logger := newStringBufferLogger()
	ctx := ContextWithLogger(context.Background(), logger)

	client := New().
		BasePath(server.URL).
		UseMiddlewares(lgMiddleware).
		Header(HeaderAPISecret, "HeaderAPISecret").
		Header(HeaderAPIKey, "HeaderAPIKey").
		Header(HeaderAPIToken, "HeaderAPIToken").
		Header(HeaderCookie, "HeaderCookie")

	err := client.PostCtx(ctx, "/", BodyJson(testData)).
		ContentTypeJson().
		AuthToken("my-token").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	logs := logger.infoBuffer.String()
	c.Require().NoError(err)
	c.Require().NotEqual(len(logs), 0)
	c.Require().True(strings.Contains(logs, fmt.Sprintf("→ [POST] %s", server.URL)))
	c.Require().False(strings.Contains(logs, "Headers:"))
	c.Require().False(strings.Contains(logs, "Content-Type=application/json"))
	c.Require().False(strings.Contains(logs, "Authorization=XXXXXXXXXXXXXXX"))
	c.Require().False(strings.Contains(logs, "X-Api-Token=XXXXXXXXXXXXXX"))
	c.Require().False(strings.Contains(logs, "Cookie=XXXXXXXXXXXX"))
	c.Require().False(strings.Contains(logs, "X-Api-Secret=XXXXXXXXXXXXXXX"))
	c.Require().False(strings.Contains(logs, "X-Api-Key=XXXXXXXXXXXX"))
	c.Require().False(strings.Contains(logs, "Body: {\"foo\":\"bar\"}"))
	c.Require().True(strings.Contains(logs, fmt.Sprintf("← [POST] %s/ - Status: 200 - Duration:", server.URL)))
	c.Require().False(strings.Contains(logs, "Response Headers:"))
	c.Require().False(strings.Contains(logs, "Content-Length=2"))
	c.Require().False(strings.Contains(logs, "Content-Type=text/plain"))
	c.Require().False(strings.Contains(logs, "Response Body: ok"))
}

func (c *ClientSuite) Test_LoggingMiddleware() {
	c.T().Log("should log everything if it is verbose except the Auth headers")
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`ok`))
	}))
	defer server.Close()

	lgMiddleware := LoggingMiddleware(true, false)
	logger := newStringBufferLogger()
	ctx := ContextWithLogger(context.Background(), logger)

	client := New().
		BasePath(server.URL).
		UseMiddlewares(lgMiddleware).
		Header(HeaderAPISecret, "HeaderAPISecret").
		Header(HeaderAPIKey, "HeaderAPIKey").
		Header(HeaderAPIToken, "HeaderAPIToken").
		Header(HeaderCookie, "HeaderCookie")

	err := client.PostCtx(ctx, "/", BodyJson(testData)).
		ContentTypeJson().
		AuthToken("my-token").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	logs := logger.infoBuffer.String()
	c.Require().NoError(err)
	c.Require().True(strings.Contains(logs, fmt.Sprintf("→ [POST] %s/", server.URL)))
	c.Require().True(strings.Contains(logs, "Headers:"))
	c.Require().True(strings.Contains(logs, "Content-Type=application/json"))
	c.Require().True(strings.Contains(logs, "Authorization=XXXXXXXXXXXXXXX"))
	c.Require().True(strings.Contains(logs, "X-Api-Token=XXXXXXXXXXXXXX"))
	c.Require().True(strings.Contains(logs, "Cookie=XXXXXXXXXXXX"))
	c.Require().True(strings.Contains(logs, "X-Api-Secret=XXXXXXXXXXXXXXX"))
	c.Require().True(strings.Contains(logs, "X-Api-Key=XXXXXXXXXXXX"))
	c.Require().True(strings.Contains(logs, "Body: {\"foo\":\"bar\"}"))
	c.Require().True(strings.Contains(logs, fmt.Sprintf("← [POST] %s/ - Status: 200 - Duration:", server.URL)))
	c.Require().True(strings.Contains(logs, "Response Headers:"))
	c.Require().True(strings.Contains(logs, "Content-Length=2"))
	c.Require().True(strings.Contains(logs, "Content-Type=text/plain"))
	c.Require().True(strings.Contains(logs, "Response Body: ok"))
}

func (c *ClientSuite) Test_LoggingMiddleware_Disabled() {
	c.T().Log("should not log anything if logger is disabled")
	c.T().Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`ok`))
	}))
	defer server.Close()

	lgMiddleware := LoggingMiddleware(false, true)
	logger := newStringBufferLogger()
	ctx := ContextWithLogger(context.Background(), logger)

	client := New().
		BasePath(server.URL).
		UseMiddlewares(lgMiddleware).
		Header(HeaderAPISecret, "HeaderAPISecret").
		Header(HeaderAPIKey, "HeaderAPIKey").
		Header(HeaderAPIToken, "HeaderAPIToken").
		Header(HeaderCookie, "HeaderCookie")

	err := client.PostCtx(ctx, "/", BodyJson(testData)).
		ContentTypeJson().
		AuthToken("my-token").
		OnReply(StatusAnyExcept(http.StatusOK), ReturnError(errors.New("unexpected status"))).
		Send()

	c.Require().NoError(err)
	c.Require().Len(logger.infoBuffer.String(), 0)
	c.Require().Len(logger.warnBuffer.String(), 0)
	c.Require().Len(logger.errorBuffer.String(), 0)
	c.Require().Len(logger.debugBuffer.String(), 0)
}
