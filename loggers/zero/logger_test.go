package zero

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/denizgursoy/inpu"
	"github.com/rs/zerolog"
)

func TestNewInpuZeroLogger(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`test`))
	}))
	defer server.Close()

	logx := zerolog.New(os.Stdout).With().Timestamp().Logger()
	zerolog.DefaultContextLogger = &logx
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	loggingMiddleware := inpu.LoggingMiddleware(inpu.LogLevelSimple)
	client := inpu.New().
		BasePath(server.URL).
		UseMiddlewares(loggingMiddleware)

	inpu.DefaultLogger = NewInpuZeroLogger()

	err := client.Post("/", nil).
		OnReply(inpu.StatusAnyExcept(http.StatusOK), inpu.ReturnError(errors.New("unexpected status"))).
		Send()
	if err != nil {
		t.FailNow()
	}
}
