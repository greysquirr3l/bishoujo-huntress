// pkg/huntress/options.go
package huntress

import (
	"net/http"
	"time"
)

// Option configures a Client instance
type Option func(*Client)

// WithBaseURL configures the API base URL for the client
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient configures a custom HTTP client for the client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithCredentials configures the API credentials for the client
func WithCredentials(apiKey, apiSecret string) Option {
	return func(c *Client) {
		c.credentials = &apiCredentials{
			APIKey:    apiKey,
			APISecret: apiSecret,
		}
	}
}

// WithUserAgent configures the user agent string for the client
func WithUserAgent(userAgent string) Option {
	return func(c *Client) {
		c.userAgent = userAgent
	}
}

// WithTimeout configures the timeout for HTTP requests
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = timeout
	}
}

// WithRetryConfig configures the retry behavior for the client
func WithRetryConfig(maxRetries int, minWait, maxWait time.Duration) Option {
	return func(c *Client) {
		c.retryConfig = &retryConfig{
			MaxRetries:   maxRetries,
			RetryWaitMin: minWait,
			RetryWaitMax: maxWait,
		}
	}
}

// WithRateLimit configures the rate limit (requests per minute) for the client
func WithRateLimit(requestsPerMinute int) Option {
	return func(c *Client) {
		if requestsPerMinute <= 0 {
			requestsPerMinute = defaultRequestsPerMinute
		}
		c.rateLimiter = newRateLimiter(requestsPerMinute)
	}
}

// WithDebug enables debug mode for the client
func WithDebug(enabled bool) Option {
	return func(c *Client) {
		c.debug = enabled
	}
}
