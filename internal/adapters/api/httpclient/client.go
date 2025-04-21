package httpclient

import (
	"context"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// RateLimiter defines the interface for rate limiting HTTP requests
type RateLimiter interface {
	Wait(ctx context.Context) error
	Allow() bool
}

// DefaultRateLimiter implements a sliding window rate limiter for the Huntress API
// which has a limit of 60 requests per minute
type DefaultRateLimiter struct {
	limiter *rate.Limiter
}

// NewDefaultRateLimiter creates a new rate limiter that allows 60 requests per minute
func NewDefaultRateLimiter() *DefaultRateLimiter {
	// 60 requests per minute = 1 request per second
	return &DefaultRateLimiter{
		limiter: rate.NewLimiter(rate.Limit(1), 60), // 1 per second, burst of 60
	}
}

// Wait blocks until a token is available or the context is done
func (r *DefaultRateLimiter) Wait(ctx context.Context) error {
	return r.limiter.Wait(ctx)
}

// Allow returns true if a token is available immediately
func (r *DefaultRateLimiter) Allow() bool {
	return r.limiter.Allow()
}

// Client is an HTTP client for interacting with the Huntress API
type Client struct {
	httpClient  *http.Client
	baseURL     string
	credentials Credentials
	rateLimiter RateLimiter
	mu          sync.Mutex // protects concurrent access to client properties
}

// Option is a function that configures a Client
type Option func(*Client)

// WithHTTPClient sets the HTTP client to use
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets the base URL for API requests
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithCredentials sets the authentication credentials
func WithCredentials(credentials Credentials) Option {
	return func(c *Client) {
		c.credentials = credentials
	}
}

// WithRateLimiter sets a custom rate limiter
func WithRateLimiter(rateLimiter RateLimiter) Option {
	return func(c *Client) {
		c.rateLimiter = rateLimiter
	}
}

// NewClient creates a new API client with the given options
func NewClient(options ...Option) *Client {
	client := &Client{
		httpClient:  http.DefaultClient,
		baseURL:     "https://api.huntress.io",
		rateLimiter: NewDefaultRateLimiter(),
	}

	for _, option := range options {
		option(client)
	}

	// Set default timeout if none was provided
	if client.httpClient.Timeout == 0 {
		client.httpClient.Timeout = 30 * time.Second
	}

	return client
}

// DoRequest executes an HTTP request with rate limiting and authentication
func (c *Client) DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Apply authentication
	if c.credentials.IsConfigured() {
		c.credentials.ApplyToRequest(req)
	}

	// Wait for rate limiter
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}

	// Execute the request
	return c.httpClient.Do(req.WithContext(ctx))
}

// Do is a more generic version of DoRequest that handles request creation
func (c *Client) Do(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	// TODO: Implement request creation, body serialization, and call DoRequest
	// This is a placeholder for the full implementation
	return nil, nil
}
