package inpu

import (
	"net/http"
	"sync"
	"time"
)

var (
	once          sync.Once
	defaultClient *http.Client
)

func getDefaultClient() *http.Client {
	once.Do(func() {
		defaultClient = &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   10,
				MaxConnsPerHost:       100,
				Proxy:                 http.ProxyFromEnvironment,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		}
	})

	return defaultClient
}
