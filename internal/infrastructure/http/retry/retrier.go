package retry

import (
	"context"
	"math"
	"net/http"
	"time"
)

// RetryConfig defines the configuration for the retry mechanism
type RetryConfig struct {
	// MaxRetries is the maximum number of retries
	MaxRetries int
	// MaxDelay is the maximum delay between retries
	MaxDelay time.Duration
	// BaseDelay is the base delay for the exponential backoff
	BaseDelay time.Duration
	// RetryableStatusCodes is a list of HTTP status codes that should trigger a retry
	RetryableStatusCodes []int
}

// DefaultRetryConfig provides sensible defaults for the retry configuration
var DefaultRetryConfig = RetryConfig{
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
	config RetryConfig
}

// NewRetrier creates a new Retrier with the given configuration
func NewRetrier(config RetryConfig) *Retrier {
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
			return nil, ctx.Err()
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
			return resp, ctx.Err()
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

// calculateBackoff computes the backoff delay for a given retry attempt
// using exponential backoff with jitter
func (r *Retrier) calculateBackoff(attempt int) time.Duration {
	// Calculate exponential backoff: baseDelay * 2^attempt
	backoff := float64(r.config.BaseDelay) * math.Pow(2, float64(attempt))

	// Add jitter: random value between 0 and backoff/2
	jitter := time.Duration(backoff / 2 * (1 - (0.5 * math.Sin(float64(attempt)))))

	// Calculate final delay
	delay := time.Duration(backoff) + jitter

	// Ensure delay doesn't exceed the maximum
	if delay > r.config.MaxDelay {
		delay = r.config.MaxDelay
	}

	return delay
}
