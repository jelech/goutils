// Package http provides HTTP client utilities with built-in retry mechanisms.
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jelech/goutils/retry"
)

// Client represents an HTTP client with retry capabilities
type Client struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

// Option represents a configuration option for HTTP client
type Option func(*Client)

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.client.Timeout = timeout
	}
}

// WithBaseURL sets the base URL for all requests
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHeaders sets default headers for all requests
func WithHeaders(headers map[string]string) Option {
	return func(c *Client) {
		if c.headers == nil {
			c.headers = make(map[string]string)
		}
		for k, v := range headers {
			c.headers[k] = v
		}
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.client = client
	}
}

// NewClient creates a new HTTP client with the given options
func NewClient(options ...Option) *Client {
	client := &Client{
		client: &http.Client{
			Timeout: time.Second * 30,
		},
		headers: make(map[string]string),
	}

	for _, option := range options {
		option(client)
	}

	return client
}

// Get performs a GET request
func (c *Client) Get(url string) (*http.Response, error) {
	return c.Request("GET", url, nil)
}

// Post performs a POST request with JSON body
func (c *Client) Post(url string, body interface{}) (*http.Response, error) {
	return c.Request("POST", url, body)
}

// Put performs a PUT request with JSON body
func (c *Client) Put(url string, body interface{}) (*http.Response, error) {
	return c.Request("PUT", url, body)
}

// Delete performs a DELETE request
func (c *Client) Delete(url string) (*http.Response, error) {
	return c.Request("DELETE", url, nil)
}

// Request performs an HTTP request with the specified method, URL, and body
func (c *Client) Request(method, url string, body interface{}) (*http.Response, error) {
	fullURL := c.buildURL(url)

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	// Set content type for JSON body
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.client.Do(req)
}

// GetWithRetry performs a GET request with retry logic
func (c *Client) GetWithRetry(url string, maxAttempts int) (*http.Response, error) {
	return c.RequestWithRetry("GET", url, nil, maxAttempts)
}

// PostWithRetry performs a POST request with retry logic
func (c *Client) PostWithRetry(url string, body interface{}, maxAttempts int) (*http.Response, error) {
	return c.RequestWithRetry("POST", url, body, maxAttempts)
}

// RequestWithRetry performs an HTTP request with retry logic
func (c *Client) RequestWithRetry(method, url string, body interface{}, maxAttempts int) (*http.Response, error) {
	var response *http.Response
	var lastErr error

	err := retry.Do(func() error {
		resp, err := c.Request(method, url, body)
		if err != nil {
			lastErr = err
			return err
		}

		// Check if the response indicates a retryable error
		if c.isRetryableStatusCode(resp.StatusCode) {
			resp.Body.Close()
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
			return lastErr
		}

		response = resp
		return nil
	}, retry.WithMaxAttempts(maxAttempts), retry.WithRetryIf(func(err error) bool {
		// Retry on network errors and 5xx status codes
		return true
	}))

	if err != nil {
		return nil, lastErr
	}

	return response, nil
}

// RequestWithContext performs an HTTP request with context
func (c *Client) RequestWithContext(ctx context.Context, method, url string, body interface{}) (*http.Response, error) {
	fullURL := c.buildURL(url)

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	// Set content type for JSON body
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.client.Do(req)
}

// DecodeJSON decodes JSON response body into the provided interface
func (c *Client) DecodeJSON(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(v)
}

// buildURL builds the full URL by combining base URL and relative URL
func (c *Client) buildURL(url string) string {
	if c.baseURL == "" || isAbsoluteURL(url) {
		return url
	}

	baseURL := c.baseURL
	if baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}

	if url[0] != '/' {
		url = "/" + url
	}

	return baseURL + url
}

// isAbsoluteURL checks if the URL is absolute
func isAbsoluteURL(url string) bool {
	return len(url) > 7 && (url[:7] == "http://" || url[:8] == "https://")
}

// isRetryableStatusCode checks if the HTTP status code is retryable
func (c *Client) isRetryableStatusCode(statusCode int) bool {
	switch statusCode {
	case http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusTooManyRequests:
		return true
	default:
		return false
	}
}
