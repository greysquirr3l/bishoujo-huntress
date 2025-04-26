// Package retry provides retry logic for HTTP requests in the Huntress API client.
package retry

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"
)

// Config defines the configuration for the retry mechanism
type Config struct {
	// MaxRetries is the maximum number of retries
	MaxRetries int
	// MaxDelay is the maximum delay between retries
	MaxDelay time.Duration
	// BaseDelay is the base delay for the exponential backoff
	BaseDelay time.Duration
	// RetryableStatusCodes is a list of HTTP status codes that should trigger a retry
	RetryableStatusCodes []int
}

// DefaultConfig provides sensible defaults for the retry configuration
var DefaultConfig = Config{
	MaxRetries: 3,
	MaxDelay:   30 * time.Second,
	BaseDelay:  500 * time.Millisecond,
	RetryableStatusCodes: []int{
		http.StatusTooManyRequests,     // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout,      // 504
	},
}

// Retrier provides retry functionality for HTTP requests
type Retrier struct {
	config Config
}

// NewRetrier creates a new Retrier with the given configuration
func NewRetrier(config Config) *Retrier {
	return &Retrier{
		config: config,
	}
}

// Do executes the given function with retry logic
// It will retry the function if it returns an error or if the response status code
// is in the list of retryable status codes
func (r *Retrier) Do(ctx context.Context, fn func() (*http.Response, error)) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= r.config.MaxRetries; attempt++ {
		// Check if context is done before making the request
		if ctx.Err() != nil {
			return nil, fmt.Errorf("context error: %w", ctx.Err())
		}

		resp, err = fn()

		// Don't retry if there's no error and the status code is not retryable
		if err == nil && !r.isRetryableStatusCode(resp.StatusCode) {
			return resp, nil
		}

		// Last attempt, return the error
		if attempt == r.config.MaxRetries {
			return resp, err
		}

		// Calculate backoff delay using exponential backoff with jitter
		delay := r.calculateBackoff(attempt)

		// Create a timer for the backoff delay
		timer := time.NewTimer(delay)

		// Wait for either the timer to expire or the context to be canceled
		select {
		case <-ctx.Done():
			timer.Stop()
			return resp, fmt.Errorf("context error: %w", ctx.Err())
		case <-timer.C:
			// Continue to the next attempt
		}
	}

	return resp, err
}

// isRetryableStatusCode checks if the given status code should trigger a retry
func (r *Retrier) isRetryableStatusCode(statusCode int) bool {
	for _, code := range r.config.RetryableStatusCodes {
		if statusCode == code {
			return true
		}
	}
	return false
}

// ExecuteWithRetry runs fn with retry logic based on the provided config.
// It expects fn to return an error; if the error is nil, it returns nil.
func ExecuteWithRetry(ctx context.Context, fn func() error, cfg *Config) error {
	var err error
	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		if ctx.Err() != nil {
			return fmt.Errorf("context error: %w", ctx.Err())
		}
		err = fn()
		if err == nil {
			return nil
		}
		delay := calculateBackoff(cfg, attempt)
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return fmt.Errorf("context error: %w", ctx.Err())
		case <-timer.C:
		}
	}
	return err
}

// calculateBackoff computes the backoff delay for a given retry attempt (stateless version for ExecuteWithRetry)
func calculateBackoff(cfg *Config, attempt int) time.Duration {
	backoff := float64(cfg.BaseDelay) * math.Pow(2, float64(attempt))
	jitter := time.Duration(backoff / 2 * (1 - (0.5 * math.Sin(float64(attempt)))))
	delay := time.Duration(backoff) + jitter
	if delay > cfg.MaxDelay {
		delay = cfg.MaxDelay
	}
	return delay
}

// calculateBackoff computes the backoff delay for a given retry attempt
func (r *Retrier) calculateBackoff(attempt int) time.Duration {
	backoff := float64(r.config.BaseDelay) * math.Pow(2, float64(attempt))
	jitter := time.Duration(backoff / 2 * (1 - (0.5 * math.Sin(float64(attempt)))))
	delay := time.Duration(backoff) + jitter
	if delay > r.config.MaxDelay {
		delay = r.config.MaxDelay
	}
	return delay
}
