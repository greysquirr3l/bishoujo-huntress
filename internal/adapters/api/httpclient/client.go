// Package httpclient provides a reusable HTTP client for Huntress API adapters.
package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Retrier defines the interface for retry logic in HTTP clients.
type Retrier interface {
	MaxRetries() int
	IsRetryableStatusCode(int) bool
	CalculateBackoff(int) time.Duration
}

// Client is a reusable HTTP client with timeout, context, rate limiting, and retry support.
type Client struct {
	HTTPClient  *http.Client
	Timeout     time.Duration
	RateLimiter *RateLimiter
	Retrier     Retrier
}

// New creates a new HTTP client with the given timeout, rate limiter, and retrier.
// New creates a new HTTP client with the given timeout, rate limiter, and retrier.
// NOTE: Replace 'interface{}' with the actual Retrier type when available.
func New(timeout time.Duration, rateLimiter *RateLimiter, retrier Retrier) *Client {
	return &Client{
		HTTPClient:  &http.Client{Timeout: timeout},
		Timeout:     timeout,
		RateLimiter: rateLimiter,
		Retrier:     retrier,
	}
}

// Do executes an HTTP request with context, rate limiting, and retries.
func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Rate limiting
	if c.RateLimiter != nil {
		if err := c.RateLimiter.Wait(ctx); err != nil {
			return nil, &RateLimitError{Message: "API rate limit exceeded"}
		}
	}

	var resp *http.Response
	var err error
	maxRetries := 0
	if c.Retrier != nil {
		maxRetries = c.Retrier.MaxRetries()
	}
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if ctx.Err() != nil {
			// Always wrap with %w so errors.Is works
			return nil, fmt.Errorf("httpclient: context error: %w", ctx.Err())
		}
		resp, err = c.HTTPClient.Do(req.WithContext(ctx))
		if err == nil && resp != nil && (resp.StatusCode < 500 || resp.StatusCode > 599) && resp.StatusCode != 429 {
			return resp, nil
		}
		// Only retry on retryable status codes (5xx, 429)
		// If resp is nil, do not retry (network or context error)
		if resp == nil || c.Retrier == nil || !c.Retrier.IsRetryableStatusCode(resp.StatusCode) {
			break
		}
		delay := c.Retrier.CalculateBackoff(attempt)
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("httpclient: context error: %w", ctx.Err())
		case <-time.After(delay):
		}
	}
	if err != nil {
		return resp, fmt.Errorf("httpclient: http client do: %w", err)
	}
	return resp, nil
}

// DoJSON executes an HTTP request and decodes the JSON response into v.
func (c *Client) DoJSON(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.Do(ctx, req)
	if err != nil {
		return resp, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing response body: %v\n", err)
		}
	}()
	if v == nil {
		return resp, nil
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return resp, fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(body))
	}
	return resp, decodeJSON(resp.Body, v)
}

func decodeJSON(r io.Reader, v interface{}) error {
	dec := json.NewDecoder(r)
	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("httpclient: decode json: %w", err)
	}
	return nil
}

// RateLimitError is returned when the API rate limit is exceeded.
type RateLimitError struct {
	Message string
}

func (e *RateLimitError) Error() string {
	return e.Message
}

// Is allows errors.Is(err, &RateLimitError{}) to work for type checks.
func (e *RateLimitError) Is(target error) bool {
	_, ok := target.(*RateLimitError)
	return ok
}
