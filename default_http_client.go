package inpu

import (
	"net"
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
		defaultClient = getNewClient()
	})

	return defaultClient
}

func getNewClient() *http.Client {
	newClient := &http.Client{
		Transport: getDefaultTransport(),
	}

	return newClient
}

func getDefaultTransport() *http.Transport {
	return &http.Transport{
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     100,
		Proxy:               http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}
