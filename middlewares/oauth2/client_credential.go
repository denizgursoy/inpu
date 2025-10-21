package oauth2

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/denizgursoy/inpu"
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

type ClientCredentialsConfig struct {
	ClientID     string
	ClientSecret string
	TokenURL     string
	Scopes       []string
}

func NewClientCredentialsMiddleware(config ClientCredentialsConfig) *ClientCredentialsMiddleware {
	ccConfig := &clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     config.TokenURL,
		Scopes:       config.Scopes,
	}

	return &ClientCredentialsMiddleware{
		tokenSource: ccConfig.TokenSource(context.Background()),
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

	request.Header.Set(inpu.HeaderAuthorization, inpu.GetTokenHeaderValue(token.AccessToken))

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

	m.mu.RLock()
	if m.token != nil && m.token.Valid() {
		defer m.mu.RUnlock()
		return m.token, nil
	}
	m.mu.RUnlock()

	token, err := m.tokenSource.Token()
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	m.token = token
	m.mu.Unlock()

	return token, nil
}

func (m *ClientCredentialsMiddleware) GetToken() (*oauth2.Token, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.token == nil {
		return nil, fmt.Errorf("no token available")
	}
	return m.token, nil
}
