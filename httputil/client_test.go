package httputil

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	assert.NotNil(t, client)
	assert.NotNil(t, client.client)
	assert.Equal(t, time.Second*30, client.client.Timeout)
}

func TestNewClient_WithOptions(t *testing.T) {
	headers := map[string]string{
		"Authorization": "Bearer token",
		"User-Agent":    "test-agent",
	}

	client := NewClient(
		WithTimeout(time.Second*10),
		WithBaseURL("https://api.example.com"),
		WithHeaders(headers),
	)

	assert.Equal(t, time.Second*10, client.client.Timeout)
	assert.Equal(t, "https://api.example.com", client.baseURL)
	assert.Equal(t, headers, client.headers)
}

func TestClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.Get(server.URL)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func TestClient_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 1}`))
	}))
	defer server.Close()

	client := NewClient()
	body := map[string]string{"name": "test"}
	resp, err := client.Post(server.URL, body)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()
}

func TestClient_WithBaseURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/users", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	resp, err := client.Get("/api/users")

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func TestClient_WithHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer token", r.Header.Get("Authorization"))
		assert.Equal(t, "test-agent", r.Header.Get("User-Agent"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	headers := map[string]string{
		"Authorization": "Bearer token",
		"User-Agent":    "test-agent",
	}
	client := NewClient(WithHeaders(headers))
	resp, err := client.Get(server.URL)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func TestClient_GetWithRetry(t *testing.T) {
	attemptCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.GetWithRetry(server.URL, 5)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 3, attemptCount)
	resp.Body.Close()
}

func TestClient_RequestWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 100)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	defer cancel()

	client := NewClient()
	_, err := client.RequestWithContext(ctx, "GET", server.URL, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestClient_DecodeJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "John", "age": 30}`))
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.Get(server.URL)
	require.NoError(t, err)

	var result struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	err = client.DecodeJSON(resp, &result)
	require.NoError(t, err)
	assert.Equal(t, "John", result.Name)
	assert.Equal(t, 30, result.Age)
}

func TestClient_DecodeJSON_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	client := NewClient()
	resp, err := client.Get(server.URL)
	require.NoError(t, err)

	var result map[string]interface{}
	err = client.DecodeJSON(resp, &result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP 404")
}

func TestBuildURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		url      string
		expected string
	}{
		{
			name:     "relative URL with base URL",
			baseURL:  "https://api.example.com",
			url:      "/users",
			expected: "https://api.example.com/users",
		},
		{
			name:     "relative URL with trailing slash in base URL",
			baseURL:  "https://api.example.com/",
			url:      "/users",
			expected: "https://api.example.com/users",
		},
		{
			name:     "absolute URL",
			baseURL:  "https://api.example.com",
			url:      "https://other.example.com/users",
			expected: "https://other.example.com/users",
		},
		{
			name:     "no base URL",
			baseURL:  "",
			url:      "https://api.example.com/users",
			expected: "https://api.example.com/users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{baseURL: tt.baseURL}
			result := client.buildURL(tt.url)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsRetryableStatusCode(t *testing.T) {
	client := &Client{}

	retryableCodes := []int{
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusTooManyRequests,
	}

	for _, code := range retryableCodes {
		assert.True(t, client.isRetryableStatusCode(code), "Status code %d should be retryable", code)
	}

	nonRetryableCodes := []int{
		http.StatusOK,
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusNotFound,
	}

	for _, code := range nonRetryableCodes {
		assert.False(t, client.isRetryableStatusCode(code), "Status code %d should not be retryable", code)
	}
}
