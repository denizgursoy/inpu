package zap

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/denizgursoy/inpu"
	"go.uber.org/zap"
)

func TestNewInpuZapLogger(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`test`))
	}))
	defer server.Close()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	loggingMiddleware := inpu.LoggingMiddleware(true, false)
	client := inpu.New().
		BasePath(server.URL).
		UseMiddlewares(loggingMiddleware)

	inpu.DefaultLogger = NewInpuLoggerFromZapLogger(logger)

	err := client.Post("/", nil).
		OnReply(inpu.StatusAnyExcept(http.StatusOK), inpu.ReturnError(errors.New("unexpected status"))).
		Send()
	if err != nil {
		t.FailNow()
	}
}
