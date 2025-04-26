package huntress

import (
	"net/http"
	"time"

	"github.com/greysquirr3l/bishoujo-huntress/internal/infrastructure/logging"
)

// clientOptions holds configuration settings for the Huntress client
// logger is an optional structured logger for the client. If nil, logging is disabled.
type clientOptions struct {
	baseURL     string
	apiKey      string
	apiSecret   string
	httpClient  *http.Client
	userAgent   string
	apiVersion  string
	timeout     time.Duration
	rateLimiter RateLimiter
	retryConfig *retryConfig
	debug       bool
	// logger is an optional structured logger for the client. If nil, logging is disabled.
	logger logging.Logger
}

// retryConfig defines retry behavior for API requests
type retryConfig struct {
	MaxRetries   int
	RetryWaitMin time.Duration
	RetryWaitMax time.Duration
}

// Option is a function that configures client options
type Option func(*clientOptions)

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

// WithDebug enables debug logging
func WithDebug(debug bool) Option {
	return func(o *clientOptions) {
		o.debug = debug
	}
}

// WithAPIVersion sets the API version
func WithAPIVersion(version string) Option {
	return func(o *clientOptions) {
		o.apiVersion = version
	}
}

// WithLogger sets a structured logger for the Huntress client. If not set, logging is disabled.
func WithLogger(logger logging.Logger) Option {
	return func(o *clientOptions) {
		o.logger = logger
	}
}
