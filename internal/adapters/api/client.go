package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/nickcampbell18/bishoujo-huntress/internal/domain/common"
	"github.com/nickcampbell18/bishoujo-huntress/internal/domain/errors"
	"golang.org/x/time/rate"
)

const (
	// DefaultBaseURL is the base URL for the Huntress API
	DefaultBaseURL = "https://api.huntress.io"
	// DefaultAPIVersion is the current version of the Huntress API
	DefaultAPIVersion = "v1"
	// DefaultUserAgent is the user agent used for API requests
	DefaultUserAgent = "Bishoujo-Huntress/1.0.0"
	// DefaultRateLimit is the default rate limit for the Huntress API (60 requests per minute)
	DefaultRateLimit = 60
)

// Client is an HTTP client for the Huntress API
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	apiKey     string
	apiSecret  string
	userAgent  string
	version    string
	limiter    *rate.Limiter
}

// ClientOption is a functional option for configuring the API client
type ClientOption func(*Client)

// NewClient creates a new Huntress API client with the provided options
func NewClient(options ...ClientOption) *Client {
	baseURL, _ := url.Parse(DefaultBaseURL)

	client := &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		version:    DefaultAPIVersion,
		userAgent:  DefaultUserAgent,
		limiter:    rate.NewLimiter(rate.Limit(DefaultRateLimit/60.0), 1), // 60 requests per minute
	}

	for _, option := range options {
		option(client)
	}

	return client
}

// WithBaseURL sets the base URL for the API client
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		if parsedURL, err := url.Parse(baseURL); err == nil {
			c.baseURL = parsedURL
		}
	}
}

// WithHTTPClient sets the HTTP client for the API client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithCredentials sets the API credentials for the API client
func WithCredentials(apiKey, apiSecret string) ClientOption {
	return func(c *Client) {
		c.apiKey = apiKey
		c.apiSecret = apiSecret
	}
}

// WithUserAgent sets the user agent for the API client
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithVersion sets the API version for the API client
func WithVersion(version string) ClientOption {
	return func(c *Client) {
		c.version = version
	}
}

// WithRateLimit sets the rate limit for the API client
func WithRateLimit(requestsPerMinute int) ClientOption {
	return func(c *Client) {
		c.limiter = rate.NewLimiter(rate.Limit(float64(requestsPerMinute)/60.0), 1)
	}
}

// newRequest creates a new HTTP request with the proper headers and authentication
func (c *Client) newRequest(ctx context.Context, method, endpoint string, body io.Reader) (*http.Request, error) {
	if c.apiKey == "" || c.apiSecret == "" {
		return nil, errors.ErrAuthMissing
	}

	// Build the full URL
	u := *c.baseURL
	u.Path = path.Join(c.version, endpoint)

	// Create the request
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrRequestCreation, err)
	}

	// Add headers
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add authentication
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.apiKey, c.apiSecret)))
	req.Header.Set("Authorization", "Basic "+auth)

	return req, nil
}

// do executes an HTTP request with rate limiting and error handling
func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	ctx := req.Context()

	// Apply rate limiting
	if err := c.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrRateLimitExceeded, err)
	}

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrRequestFailed, err)
	}
	defer resp.Body.Close()

	// Handle non-2xx responses
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.handleErrorResponse(resp)
	}

	// Parse the response if a target is provided
	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return nil, fmt.Errorf("%w: %v", errors.ErrResponseParsing, err)
		}
	}

	return resp, nil
}

// handleErrorResponse handles non-successful HTTP responses
func (c *Client) handleErrorResponse(resp *http.Response) (*http.Response, error) {
	var apiErr struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	// Try to decode the error response
	if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
		// If decoding fails, create a generic error
		return resp, fmt.Errorf("%w: %s", errors.ErrAPIResponse, resp.Status)
	}

	// Map status codes to domain errors
	switch resp.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return resp, fmt.Errorf("%w: %s - %s", errors.ErrAuthentication, apiErr.Error.Code, apiErr.Error.Message)
	case http.StatusNotFound:
		return resp, fmt.Errorf("%w: %s - %s", errors.ErrResourceNotFound, apiErr.Error.Code, apiErr.Error.Message)
	case http.StatusTooManyRequests:
		return resp, fmt.Errorf("%w: %s", errors.ErrRateLimitExceeded, apiErr.Error.Message)
	default:
		return resp, fmt.Errorf("%w: %s - %s", errors.ErrAPIResponse, apiErr.Error.Code, apiErr.Error.Message)
	}
}

// Get performs a GET request to the specified endpoint
func (c *Client) Get(ctx context.Context, endpoint string, params url.Values, v interface{}) error {
	if params != nil && len(params) > 0 {
		endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
	}

	req, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, v)
	return err
}

// Post performs a POST request to the specified endpoint
func (c *Client) Post(ctx context.Context, endpoint string, body interface{}, v interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("%w: %v", errors.ErrRequestCreation, err)
		}
		bodyReader = io.NopCloser(io.ByteReader(jsonBody))
	}

	req, err := c.newRequest(ctx, http.MethodPost, endpoint, bodyReader)
	if err != nil {
		return err
	}

	_, err = c.do(req, v)
	return err
}

// Put performs a PUT request to the specified endpoint
func (c *Client) Put(ctx context.Context, endpoint string, body interface{}, v interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("%w: %v", errors.ErrRequestCreation, err)
		}
		bodyReader = io.NopCloser(io.ByteReader(jsonBody))
	}

	req, err := c.newRequest(ctx, http.MethodPut, endpoint, bodyReader)
	if err != nil {
		return err
	}

	_, err = c.do(req, v)
	return err
}

// Delete performs a DELETE request to the specified endpoint
func (c *Client) Delete(ctx context.Context, endpoint string, v interface{}) error {
	req, err := c.newRequest(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	_, err = c.do(req, v)
	return err
}

// GetPaginated performs a GET request with pagination support
func (c *Client) GetPaginated(ctx context.Context, endpoint string, params url.Values, v interface{}) (*common.Pagination, error) {
	if params == nil {
		params = url.Values{}
	}

	if params.Get("page") == "" {
		params.Set("page", "1")
	}

	if params.Get("limit") == "" {
		params.Set("limit", "100")
	}

	endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
	req, err := c.newRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req, v)
	if err != nil {
		return nil, err
	}

	// Extract pagination information from headers
	pagination := &common.Pagination{
		CurrentPage:  1,
		ItemsPerPage: 100,
		TotalItems:   0,
		TotalPages:   1,
	}

	// Parse pagination headers if they exist
	if totalItems := resp.Header.Get("X-Total-Count"); totalItems != "" {
		fmt.Sscanf(totalItems, "%d", &pagination.TotalItems)
	}

	if totalPages := resp.Header.Get("X-Total-Pages"); totalPages != "" {
		fmt.Sscanf(totalPages, "%d", &pagination.TotalPages)
	}

	if currentPage := resp.Header.Get("X-Page"); currentPage != "" {
		fmt.Sscanf(currentPage, "%d", &pagination.CurrentPage)
	}

	if itemsPerPage := resp.Header.Get("X-Per-Page"); itemsPerPage != "" {
		fmt.Sscanf(itemsPerPage, "%d", &pagination.ItemsPerPage)
	}

	return pagination, nil
}
