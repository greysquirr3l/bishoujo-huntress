// Package huntress provides a client for the Huntress API
package huntress

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/greysquirr3l/bishoujo-huntress/internal/infrastructure/logging"
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
	Logger      logging.Logger

	// Services for interacting with different API parts
	Account      AccountService
	Agent        AgentService
	Organization OrganizationService
	Incident     IncidentService
	Report       ReportService
	Billing      BillingService
	Webhook      WebhookService // <-- Added for webhook support
}

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
		Logger:      options.logger,
	}

	// Initialize services
	client.Account = &accountService{client: client}
	client.Agent = &agentService{client: client}
	client.Organization = &organizationService{client: client}
	client.Incident = &incidentService{client: client}
	client.Report = &reportService{client: client}
	client.Billing = &billingService{client: client}
	client.Webhook = NewWebhookService(client)

	return client
}

// NewRequest creates a new API request
func (c *Client) NewRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	url := c.baseURL + path
	var bodyReader io.Reader
	if body != nil {
		var ok bool
		bodyReader, ok = body.(io.Reader)
		if !ok {
			return nil, fmt.Errorf("body must implement io.Reader")
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+basicAuth(c.apiKey, c.apiSecret))
	req.Header.Set("User-Agent", c.userAgent)

	if c.Logger != nil {
		c.Logger.Debug("Creating new request", logging.String("method", method), logging.String("url", url))
	}

	return req, nil
}

// Do sends an API request and returns the response
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	if c.Logger != nil {
		c.Logger.Debug("Sending request", logging.String("method", req.Method), logging.String("url", req.URL.String()))
	}

	// Apply rate limiting if configured
	if c.rateLimiter != nil {
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limit error: %w", err)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if c.Logger != nil {
			c.Logger.Error("Request failed", logging.Error("error", err))
		}
		return nil, fmt.Errorf("request failed: %w", err)
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}

	if v == nil {
		return resp, nil
	}

	// Check for error responses
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		if c.Logger != nil {
			c.Logger.Warn("API error response", logging.Int("status", resp.StatusCode), logging.String("body", string(bodyBytes)))
		}
		return resp, fmt.Errorf("API error: status code %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Decode the response
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		if c.Logger != nil {
			c.Logger.Error("Error decoding response", logging.Error("error", err))
		}
		return resp, fmt.Errorf("error decoding response: %w", err)
	}
	return resp, nil
}

// basicAuth creates a basic auth header value from credentials
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
