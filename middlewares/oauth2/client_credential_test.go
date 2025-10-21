package oauth2

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/denizgursoy/inpu"
	"golang.org/x/oauth2/clientcredentials"
)

func TestClientCredentialsMiddleware(t *testing.T) {
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "test-token-cc",
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
	}))
	defer tokenServer.Close()

	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token-cc" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer apiServer.Close()

	config := &clientcredentials.Config{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		TokenURL:     tokenServer.URL + "/token",
		Scopes:       []string{"read", "write"},
	}
	middleware := NewClientCredentialsMiddleware(config)

	client := inpu.New().
		BasePath(apiServer.URL).
		UseMiddlewares(middleware)

	err := client.Get("/api/data").Send()
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
}
