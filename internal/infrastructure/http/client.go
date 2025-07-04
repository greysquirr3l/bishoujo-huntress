// Package http provides the HTTP client and related utilities for the Huntress API client.
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/greysquirr3l/bishoujo-huntress/internal/infrastructure/http/retry"
	"github.com/greysquirr3l/bishoujo-huntress/internal/infrastructure/logging"
)

// APIError represents an error returned by the Huntress API
type APIError struct {
	StatusCode int
	Message    string
	RequestID  string
	Details    map[string]interface{}
	RawBody    []byte
}

func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("API error: %d - %s (Request ID: %s)", e.StatusCode, e.Message, e.RequestID)
	}
	return fmt.Sprintf("API error: %d - %s", e.StatusCode, e.Message)
}

// Client is an HTTP client for the Huntress API
type Client struct {
	BaseURL     *url.URL
	HTTPClient  *http.Client
	APIKey      string
	APISecret   string
	UserAgent   string
	RetryConfig *retry.Config
}

// NewClient creates a new HTTP client
func NewClient(baseURL string, apiKey string, apiSecret string, opts ...ClientOption) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	client := &Client{
		BaseURL: parsedURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		APIKey:      apiKey,
		APISecret:   apiSecret,
		UserAgent:   "Bishoujo-Huntress/0.1.0",
		RetryConfig: &retry.DefaultConfig,
	}

	// Apply all client options
	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithHTTPClient sets the HTTP client to use
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// WithUserAgent sets the default User-Agent header for requests
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.UserAgent = userAgent
	}
}

// WithPerRequestUserAgent allows setting a custom User-Agent per request via context or RequestOptions
func WithPerRequestUserAgent() ClientOption {
	// This is a no-op; actual per-request logic is handled in Do()
	return func(_ *Client) {}
}

// WithRetryConfig sets the retry configuration
func WithRetryConfig(retryConfig *retry.Config) ClientOption {
	return func(c *Client) {
		c.RetryConfig = retryConfig
	}
}

// RequestOptions represents options for a request
type RequestOptions struct {
	Headers map[string]string
	Query   url.Values
}

// Pagination represents pagination information from API responses
type Pagination struct {
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	TotalItems   int `json:"total_items"`
	ItemsPerPage int `json:"items_per_page"`
}

// Do performs an HTTP request
// TODO: Refactor this function to reduce cyclomatic complexity.
func (c *Client) Do(ctx context.Context, method, path string, body, result interface{}, opts *RequestOptions) (*http.Response, error) {
	// Create the request URL
	reqURL, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("invalid request path: %w", err)
	}

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)

	// Set Basic Auth
	if c.APIKey != "" && c.APISecret != "" {
		req.SetBasicAuth(c.APIKey, c.APISecret)
	}

	// Set additional headers
	if opts != nil && opts.Headers != nil {
		for k, v := range opts.Headers {
			req.Header.Set(k, v)
		}
	}

	// Set query parameters
	if opts != nil && opts.Query != nil {
		q := req.URL.Query()
		for k, values := range opts.Query {
			for _, v := range values {
				q.Add(k, v)
			}
		}
		req.URL.RawQuery = q.Encode()
	}

	// Use Retrier.Do to retry on both errors and retryable status codes
	retrier := retry.NewRetrier(*c.RetryConfig)
	var resp *http.Response
	resp, err = retrier.Do(ctx, func() (*http.Response, error) {
		return c.HTTPClient.Do(req)
	})
	if err != nil {
		if resp != nil {
			if closeErr := resp.Body.Close(); closeErr != nil {
				logging.GetLogger().Error("failed to close response body", logging.Field{Key: "error", Value: closeErr})
			}
		}
		return nil, fmt.Errorf("request failed after retries: %w", err)
	}

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		if closeErr := resp.Body.Close(); closeErr != nil {
			logging.GetLogger().Error("failed to close response body", logging.Field{Key: "error", Value: closeErr})
		}
		return resp, fmt.Errorf("error reading response body: %w", err)
	}
	if closeErr := resp.Body.Close(); closeErr != nil {
		logging.GetLogger().Error("failed to close response body", logging.Field{Key: "error", Value: closeErr})
	}

	// Check for error response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			RawBody:    respBody,
			RequestID:  resp.Header.Get("X-Request-Id"),
		}

		// Try to parse error message from response
		var errResp map[string]interface{}
		if err := json.Unmarshal(respBody, &errResp); err == nil {
			if msg, ok := errResp["message"].(string); ok {
				apiErr.Message = msg
			}
			if details, ok := errResp["details"].(map[string]interface{}); ok {
				apiErr.Details = details
			}
		} else {
			apiErr.Message = string(respBody)
		}

		return resp, apiErr
	}

	// Parse response body if result is provided
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return resp, fmt.Errorf("error parsing response body: %w", err)
		}
	}

	return resp, nil
}

// GetPagination extracts pagination information from an HTTP response
func GetPagination(resp *http.Response) (*Pagination, error) {
	if resp == nil {
		return nil, fmt.Errorf("response is nil")
	}

	// Extract pagination from headers if available
	pagination := &Pagination{}

	// If headers don't contain pagination info, return default values
	// This is a simplification - you'll need to adjust based on Huntress API's
	// actual pagination implementation
	pagination.CurrentPage = 1
	pagination.TotalPages = 1
	pagination.TotalItems = 0
	pagination.ItemsPerPage = 0

	return pagination, nil
}

// PaginatedResponse is a generic paginated response wrapper
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, path string, result interface{}, opts *RequestOptions) (*http.Response, error) {
	return c.Do(ctx, http.MethodGet, path, nil, result, opts)
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, path string, body, result interface{}, opts *RequestOptions) (*http.Response, error) {
	return c.Do(ctx, http.MethodPost, path, body, result, opts)
}

// Put performs a PUT request
func (c *Client) Put(ctx context.Context, path string, body, result interface{}, opts *RequestOptions) (*http.Response, error) {
	return c.Do(ctx, http.MethodPut, path, body, result, opts)
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, path string, result interface{}, opts *RequestOptions) (*http.Response, error) {
	return c.Do(ctx, http.MethodDelete, path, nil, result, opts)
}
