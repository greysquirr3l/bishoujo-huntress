// Package huntress provides a client for the Huntress API
package huntress

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// RateLimiter defines the interface for rate limiting API requests
type RateLimiter interface {
	Wait(ctx context.Context) error
	Reserve() (bool, time.Duration)
}

// Client is the main client for the Huntress API
type Client struct {
	httpClient  *http.Client
	baseURL     string
	apiKey      string
	apiSecret   string
	userAgent   string
	apiVersion  string
	rateLimiter RateLimiter
	retryConfig *retryConfig
	debug       bool

	// Services
	Account      AccountService
	Agent        AgentService
	Organization OrganizationService
	Incident     IncidentService
	Report       ReportService
	Billing      BillingService
}

// clientOptions holds the options for creating a new client
type clientOptions struct {
	httpClient  *http.Client
	baseURL     string
	apiKey      string
	apiSecret   string
	userAgent   string
	apiVersion  string
	timeout     time.Duration
	rateLimiter RateLimiter
	retryConfig *retryConfig
	debug       bool
}

// retryConfig defines retry behavior for API requests
type retryConfig struct {
	MaxRetries   int
	RetryWaitMin time.Duration
	RetryWaitMax time.Duration
}

// Option is a function that configures client options
type Option func(*clientOptions)

// New creates a new Huntress API client
func New(opts ...Option) *Client {
	options := &clientOptions{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:    "https://api.huntress.io/v1",
		userAgent:  "bishoujo-huntress/1.0.0",
		apiVersion: "v1",
	}

	// Apply options
	for _, opt := range opts {
		opt(options)
	}

	// Create the HTTP client
	if options.httpClient == nil {
		options.httpClient = &http.Client{
			Timeout: options.timeout,
		}
	}

	client := &Client{
		httpClient:  options.httpClient,
		baseURL:     options.baseURL,
		apiKey:      options.apiKey,
		apiSecret:   options.apiSecret,
		userAgent:   options.userAgent,
		apiVersion:  options.apiVersion,
		rateLimiter: options.rateLimiter,
	}

	// Initialize services
	client.Account = &accountService{client: client}
	client.Agent = &agentService{client: client}
	client.Organization = &organizationService{client: client}
	client.Incident = &incidentService{client: client}
	client.Report = &reportService{client: client}
	client.Billing = &billingService{client: client}

	return client
}

// WithCredentials sets the API credentials for authentication
func WithCredentials(apiKey, apiSecret string) Option {
	return func(o *clientOptions) {
		o.apiKey = apiKey
		o.apiSecret = apiSecret
	}
}

// WithBaseURL sets the API base URL
func WithBaseURL(baseURL string) Option {
	return func(o *clientOptions) {
		o.baseURL = baseURL
	}
}

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) Option {
	return func(o *clientOptions) {
		o.timeout = timeout
		if o.httpClient != nil {
			o.httpClient.Timeout = timeout
		}
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) Option {
	return func(o *clientOptions) {
		o.httpClient = client
	}
}

// WithUserAgent sets the user agent string
func WithUserAgent(userAgent string) Option {
	return func(o *clientOptions) {
		o.userAgent = userAgent
	}
}

// WithRateLimiter sets a rate limiter for API requests
func WithRateLimiter(limiter RateLimiter) Option {
	return func(o *clientOptions) {
		o.rateLimiter = limiter
	}
}

// WithRetryConfig configures the retry behavior for the client
func WithRetryConfig(maxRetries int, minWait, maxWait time.Duration) Option {
	return func(o *clientOptions) {
		o.retryConfig = &retryConfig{
			MaxRetries:   maxRetries,
			RetryWaitMin: minWait,
			RetryWaitMax: maxWait,
		}
	}
}

// NewRequest creates a new API request
func (c *Client) NewRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	url := c.baseURL + path
	var bodyReader io.Reader

	if body != nil {
		bodyData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+basicAuth(c.apiKey, c.apiSecret))
	req.Header.Set("User-Agent", c.userAgent)

	return req, nil
}

// Do sends an API request and returns the response
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	// Apply rate limiting if configured
	if c.rateLimiter != nil {
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limit error: %w", err)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if v == nil {
		return resp, nil
	}

	// Check for error responses
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return resp, fmt.Errorf("API error: status code %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Decode the response
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return resp, fmt.Errorf("error decoding response: %w", err)
	}

	return resp, nil
}

// basicAuth creates a basic auth header value from credentials
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
