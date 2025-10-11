package inpu

import (
	"net/http"
	"sync"

	"github.com/hashicorp/go-cleanhttp"
)

var (
	once          sync.Once
	defaultClient *http.Client
)

func getDefaultClient() *http.Client {
	once.Do(func() {
		defaultClient = cleanhttp.DefaultClient()
	})

	return defaultClient
}
