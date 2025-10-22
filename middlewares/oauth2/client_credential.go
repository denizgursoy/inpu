package oauth2

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type ClientCredentialsMiddleware struct {
	tokenSource  oauth2.TokenSource
	mu           sync.RWMutex
	token        *oauth2.Token
	refreshMutex sync.Mutex
	next         http.RoundTripper
}

func NewClientCredentialsMiddleware(config clientcredentials.Config) *ClientCredentialsMiddleware {
	return &ClientCredentialsMiddleware{
		tokenSource: config.TokenSource(context.Background()),
	}
}

func (m *ClientCredentialsMiddleware) ID() string {
	return "oauth-client-secret-middleware"
}

func (m *ClientCredentialsMiddleware) Priority() int {
	return 4
}

func (m *ClientCredentialsMiddleware) Apply(next http.RoundTripper) http.RoundTripper {
	m.next = next

	return m
}

func (m *ClientCredentialsMiddleware) RoundTrip(request *http.Request) (*http.Response, error) {
	token, err := m.getValidToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth2 token: %w", err)
	}

	token.SetAuthHeader(request)

	return m.next.RoundTrip(request)
}

func (m *ClientCredentialsMiddleware) getValidToken() (*oauth2.Token, error) {
	m.mu.RLock()
	if m.token != nil && m.token.Valid() {
		defer m.mu.RUnlock()

		return m.token, nil
	}
	m.mu.RUnlock()

	m.refreshMutex.Lock()
	defer m.refreshMutex.Unlock()

	token, err := m.tokenSource.Token()
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.token = token

	return token, nil
}
