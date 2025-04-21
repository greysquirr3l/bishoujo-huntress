package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	// DefaultBaseURL is the default Huntress API URL
	DefaultBaseURL = "https://api.huntress.io"
	// RateLimitPerMinute defines the maximum requests per minute as per API docs
	RateLimitPerMinute = 60
)

// HTTPClient wraps the standard http.Client with additional functionality
// for interacting with the Huntress API
type HTTPClient struct {
	client      *http.Client
	baseURL     string
	apiKey      string
	apiSecret   string
	userAgent   string
	rateLimiter *rate.Limiter
	retrier     *Retrier
	mu          sync.Mutex
}

// HTTPClientOption is a function that configures an HTTPClient
type HTTPClientOption func(*HTTPClient)

// NewHTTPClient creates a new HTTP client for interacting with the Huntress API
func NewHTTPClient(options ...HTTPClientOption) *HTTPClient {
	c := &HTTPClient{
		client:      &http.Client{Timeout: 30 * time.Second},
		baseURL:     DefaultBaseURL,
		rateLimiter: rate.NewLimiter(rate.Limit(RateLimitPerMinute/60.0), 1), // 1 request per second on average
		retrier:     NewRetrier(DefaultRetryConfig()),
	}

	for _, option := range options {
		option(c)
	}

	return c
}

// WithBaseURL sets the base URL for the client
func WithBaseURL(baseURL string) HTTPClientOption {
	return func(c *HTTPClient) {
		c.baseURL = baseURL
	}
}

// WithCredentials sets the API credentials
func WithCredentials(apiKey, apiSecret string) HTTPClientOption {
	return func(c *HTTPClient) {
		c.apiKey = apiKey
		c.apiSecret = apiSecret
	}
}

// WithHTTPClient sets the underlying HTTP client
func WithHTTPClient(client *http.Client) HTTPClientOption {
	return func(c *HTTPClient) {
		c.client = client
	}
}

// WithUserAgent sets the User-Agent header for requests
func WithUserAgent(userAgent string) HTTPClientOption {
	return func(c *HTTPClient) {
		c.userAgent = userAgent
	}
}

// WithRateLimit sets a custom rate limit for the client
func WithRateLimit(requestsPerMinute int) HTTPClientOption {
	return func(c *HTTPClient) {
		if requestsPerMinute <= 0 {
			requestsPerMinute = RateLimitPerMinute
		}
		// Convert requests per minute to requests per second with burst capacity
		c.rateLimiter = rate.NewLimiter(rate.Limit(float64(requestsPerMinute)/60.0), 5)
	}
}

// WithRetryConfig sets the retry configuration for failed requests
func WithRetryConfig(config RetryConfig) HTTPClientOption {
	return func(c *HTTPClient) {
		c.retrier = NewRetrier(config)
	}
}

// Do executes an HTTP request with authentication, rate limiting,
// and retrying capabilities
func (c *HTTPClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Apply authentication
	c.applyAuth(req)

	// Apply user agent if set
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	// Wait for rate limit token (respect context cancellation)
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter wait: %w", err)
	}

	// Execute request with retry logic
	resp, err := c.retrier.Do(ctx, func(ctx context.Context) (*http.Response, error) {
		// Create a clone of the request with the updated context
		reqWithCtx := req.Clone(ctx)
		return c.client.Do(reqWithCtx)
	})

	if err != nil {
		return nil, fmt.Errorf("request execution failed: %w", err)
	}

	return resp, nil
}

// applyAuth adds the Basic Auth credentials to the request
func (c *HTTPClient) applyAuth(req *http.Request) {
	// Base64 encode API key and secret for Basic Auth
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.apiKey, c.apiSecret)))
	req.Header.Set("Authorization", "Basic "+auth)
}
