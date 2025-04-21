// Package huntress provides a Go client for the Huntress API.
package huntress

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

const (
	// Default API URL
	defaultBaseURL = "https://api.huntress.io/v1"
	// Default user agent
	defaultUserAgent = "bishoujo-huntress/1.0.0"
	// Default requests per minute (API rate limit)
	defaultRequestsPerMinute = 60
	// Default max retries
	defaultMaxRetries = 3
	// Default min retry wait
	defaultRetryWaitMin = 100 * time.Millisecond
	// Default max retry wait
	defaultRetryWaitMax = 2 * time.Second
)

// Client is the Huntress API client
type Client struct {
	// HTTP client used to communicate with the API
	httpClient *http.Client

	// Base URL for API requests
	baseURL string

	// User agent for API requests
	userAgent string

	// Authentication credentials
	credentials *apiCredentials

	// Retry configuration
	retryConfig *retryConfig

	// Rate limiter to respect API limits
	rateLimiter *rate.Limiter

	// Debug mode
	debug bool

	// Services
	Account      AccountService
	Organization OrganizationService
	Agent        AgentService
	Incident     IncidentService
	Report       ReportService
	Billing      BillingService
}

// apiCredentials stores API authentication credentials
type apiCredentials struct {
	APIKey    string
	APISecret string
}

// retryConfig stores retry configuration
type retryConfig struct {
	MaxRetries   int
	RetryWaitMin time.Duration
	RetryWaitMax time.Duration
}

// rateLimiter maintains the API rate limit
type rateLimiter struct {
	limiter *rate.Limiter
}

// New creates a new Huntress API client
func New(options ...Option) *Client {
	client := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:   defaultBaseURL,
		userAgent: defaultUserAgent,
		retryConfig: &retryConfig{
			MaxRetries:   defaultMaxRetries,
			RetryWaitMin: defaultRetryWaitMin,
			RetryWaitMax: defaultRetryWaitMax,
		},
	}

	// Apply options
	for _, option := range options {
		option(client)
	}

	// Create rate limiter if not provided
	if client.rateLimiter == nil {
		client.rateLimiter = rate.NewLimiter(rate.Limit(defaultRequestsPerMinute/60.0), 1)
	}

	// Initialize services
	client.Account = &accountService{client: client}
	client.Organization = &organizationService{client: client}
	client.Agent = &agentService{client: client}
	client.Incident = &incidentService{client: client}
	client.Report = &reportService{client: client}
	client.Billing = &billingService{client: client}

	return client
}

// newRateLimiter creates a new rate limiter for the given requests per minute
func newRateLimiter(requestsPerMinute int) *rate.Limiter {
	return rate.NewLimiter(rate.Limit(float64(requestsPerMinute)/60.0), 1)
}

// NewRequest creates a new API request
func (c *Client) NewRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	// Construct the full URL
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, err
	}

	// Create request body if needed
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	// Set authentication if credentials are provided
	if c.credentials != nil {
		auth := base64.StdEncoding.EncodeToString(
			[]byte(fmt.Sprintf("%s:%s", c.credentials.APIKey, c.credentials.APISecret)),
		)
		req.Header.Set("Authorization", "Basic "+auth)
	}

	return req, nil
}

// Do performs an API request and decodes the response into v
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	// Apply rate limiting
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}

	var resp *http.Response
	var err error
	var attemptNum int

	// Keep retrying until we get a success response, reach max retries, or hit a non-retryable error
	for attemptNum = 0; attemptNum <= c.retryConfig.MaxRetries; attemptNum++ {
		// Only wait if this is a retry
		if attemptNum > 0 {
			waitTime := c.getBackoffDuration(attemptNum)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(waitTime):
				// Continue with retry
			}

			// Create a new request for each retry to ensure a fresh request
			// This is important especially for requests with a body
			if req.Body != nil {
				// If we have a body, we need to recreate the request
				// This is necessary because we can't reuse a request with a body
				newReq, err := c.NewRequest(ctx, req.Method, strings.TrimPrefix(req.URL.String(), c.baseURL), nil)
				if err != nil {
					return nil, err
				}
				// Copy headers from the original request
				newReq.Header = req.Header
				req = newReq
			}
		}

		// Execute the request
		resp, err = c.httpClient.Do(req)

		// If there was a network-level error, retry
		if err != nil {
			continue
		}

		// Check if we need to retry based on response status
		if !c.shouldRetry(resp) {
			break
		}

		// Close the response body to reuse the connection
		if resp.Body != nil {
			resp.Body.Close()
		}
	}

	// If all retries failed with network errors
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Check for API errors
	if err := c.checkResponse(resp); err != nil {
		return resp, err
	}

	// If v is provided, decode the response into it
	if v != nil && resp.StatusCode != http.StatusNoContent {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
}

// shouldRetry returns true if the request should be retried based on the response
func (c *Client) shouldRetry(resp *http.Response) bool {
	// Retry on rate limit errors (429)
	// Retry on server errors (5xx)
	// Retry on temporary issues (408, 409, 423, 425, 500, 502, 503, 504)
	switch resp.StatusCode {
	case http.StatusTooManyRequests, http.StatusRequestTimeout, http.StatusConflict,
		http.StatusLocked, http.StatusTooEarly, http.StatusInternalServerError,
		http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return resp.StatusCode >= 500
	}
}

// getBackoffDuration calculates the backoff duration for a retry attempt
func (c *Client) getBackoffDuration(attemptNum int) time.Duration {
	minTime := float64(c.retryConfig.RetryWaitMin)
	maxTime := float64(c.retryConfig.RetryWaitMax)

	// Exponential backoff with jitter
	// Formula: min(maxTime, minTime * 2^attempt) + random jitter
	timeMultiplier := math.Pow(2, float64(attemptNum))
	delay := minTime * timeMultiplier

	// Cap the delay to max wait time
	if delay > maxTime {
		delay = maxTime
	}

	// Add jitter (random variance up to 25% of the delay)
	jitter := rand.Float64() * delay * 0.25
	delay += jitter

	return time.Duration(delay)
}

// checkResponse checks the API response for errors
func (c *Client) checkResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	errorResponse := &ErrorResponse{
		StatusCode: resp.StatusCode,
	}

	// Try to decode the error response
	data, err := io.ReadAll(resp.Body)
	if err == nil && len(data) > 0 {
		// Try to unmarshal into ErrorResponse
		if err := json.Unmarshal(data, errorResponse); err != nil {
			// If unmarshal fails, set a generic message
			errorResponse.Message = fmt.Sprintf("failed to parse error response: %s", string(data))
		}
	}

	// If no specific message was decoded, use the status text
	if errorResponse.Message == "" {
		errorResponse.Message = resp.Status
	}

	// Set the request ID if available
	errorResponse.RequestID = resp.Header.Get("X-Request-ID")

	return errorResponse
}
