package couchdb

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// CustomHTTPClient represents an HTTP client with configurable settings.
// It allows making HTTP requests with options for timeout and retries.
type CustomHTTPClient struct {
	baseURL    string        // Base URL for the HTTP client
	client     *http.Client  // HTTP client for making requests
	maxRetries int           // Maximum number of retries for failed requests
	retryWait  time.Duration // Duration to wait between retries
	timeout    time.Duration // Timeout for each HTTP request
}

// NewCustomHTTPClient creates a new CustomHTTPClient with the specified base URL and configuration options.
// It returns a pointer to the created CustomHTTPClient instance.
func NewCustomHTTPClient(baseURL string, maxRetries int, retryWait, timeout time.Duration) *CustomHTTPClient {
	return &CustomHTTPClient{
		baseURL:    baseURL,
		client:     &http.Client{},
		maxRetries: maxRetries,
		retryWait:  retryWait,
		timeout:    timeout,
	}
}

// makeRequest makes an HTTP request with the provided method, endpoint, and body.
// It handles retries according to the configured settings.
// The function returns the response status code, body, and any error encountered.
func (c *CustomHTTPClient) makeRequest(ctx context.Context, method, endpoint string, body interface{}) (int, []byte, error) {
	url := c.baseURL + endpoint

	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			return 0, nil, err
		}
	}

	var respBody []byte
	var respCode int
	for i := 0; i < c.maxRetries; i++ {
		req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
		if err != nil {
			return 0, nil, err
		}

		req.Header.Set("Content-Type", "application/json")

		ctx, cancel := context.WithTimeout(req.Context(), c.timeout)
		defer cancel()
		req = req.WithContext(ctx)

		resp, err := c.client.Do(req)
		if err != nil {
			if i == c.maxRetries-1 {
				return 0, nil, err
			}
			time.Sleep(c.retryWait)
			continue
		}
		defer resp.Body.Close()

		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return 0, nil, err
		}

		respCode = resp.StatusCode

		if respCode < 500 {
			break
		}
		time.Sleep(c.retryWait)
	}
	return respCode, respBody, nil
}

// Get sends a GET request to the specified endpoint with optional request body.
// It returns the response status code, body, and any error encountered.
func (c *CustomHTTPClient) Get(ctx context.Context, endpoint string) (int, []byte, error) {
	return c.makeRequest(ctx, "GET", endpoint, nil)
}

// Post sends a POST request to the specified endpoint with the provided body.
// It returns the response status code, body, and any error encountered.
func (c *CustomHTTPClient) Post(ctx context.Context, endpoint string, body interface{}) (int, []byte, error) {
	return c.makeRequest(ctx, "POST", endpoint, body)
}

// Put sends a PUT request to the specified endpoint with the provided body.
// It returns the response status code, body, and any error encountered.
func (c *CustomHTTPClient) Put(ctx context.Context, endpoint string, body interface{}) (int, []byte, error) {
	return c.makeRequest(ctx, "PUT", endpoint, body)
}

// Delete sends a DELETE request to the specified endpoint.
// It returns the response status code, body, and any error encountered.
func (c *CustomHTTPClient) Delete(ctx context.Context, endpoint string) (int, []byte, error) {
	return c.makeRequest(ctx, "DELETE", endpoint, nil)
}

// Head sends a HEAD request to the specified endpoint.
// It returns the response status code, body, and any error encountered.
func (c *CustomHTTPClient) Head(ctx context.Context, endpoint string) (int, []byte, error) {
	return c.makeRequest(ctx, "HEAD", endpoint, nil)
}
