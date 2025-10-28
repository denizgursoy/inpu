package inpu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testCount = 5_000

func Benchmark_Client_Standard(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"foo":"bar"}`))
	}))
	defer server.Close()

	client := http.DefaultClient

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		body := generateRandomBody()
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, server.URL, bytes.NewReader(bodyBytes))
		for i := 0; i < testCount; i++ {
			headerName := fmt.Sprintf("X-Custom-%s", generateRandomString(8))
			headerValue := generateRandomString(20)
			req.Header.Add(headerName, headerValue)
		}

		query := req.URL.Query()
		for i := 0; i < testCount; i++ {
			paramName := fmt.Sprintf("param_%s", generateRandomString(8))
			paramValue := generateRandomString(15)
			query.Add(paramName, paramValue)
		}
		req.URL.RawQuery = query.Encode()

		resp, _ := client.Do(req)
		all, _ := io.ReadAll(resp.Body)
		parsedMap := make(map[string]any)
		json.Unmarshal(all, &parsedMap)

		resp.Body.Close()
	}
}

func Benchmark_Client_Inpu(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"foo":"bar"}`))
	}))
	defer server.Close()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := Post(server.URL, BodyJson(generateRandomBody()))
		for i := 0; i < testCount; i++ {
			headerName := fmt.Sprintf("X-Custom-%s", generateRandomString(8))
			headerValue := generateRandomString(20)
			req.Header(headerName, headerValue)
		}

		for i := 0; i < testCount; i++ {
			paramName := fmt.Sprintf("param_%s", generateRandomString(8))
			paramValue := generateRandomString(15)
			req.QueryString(paramName, paramValue)
		}
		parsedMap := make(map[string]any)

		err := req.
			OnReplyIf(StatusIsOk, ThenUnmarshalJsonTo(&parsedMap)).
			OnReplyIf(StatusAny, ThenReturnDefaultError).
			Send()
		if err != nil {
			return
		}
	}
}

// generateRandomString creates a random string of length n
func generateRandomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// generateRandomBody creates a map with random data
func generateRandomBody() map[string]interface{} {
	return map[string]interface{}{
		"user_id":   rand.Intn(10000),
		"email":     generateRandomString(10) + "@example.com",
		"username":  generateRandomString(12),
		"timestamp": "2025-10-22T10:30:00Z",
		"data": map[string]interface{}{
			"field1": generateRandomString(20),
			"field2": generateRandomString(15),
			"count":  rand.Intn(100),
		},
	}
}
