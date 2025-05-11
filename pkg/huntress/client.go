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

	"github.com/greysquirr3l/bishoujo-huntress/internal/infrastructure/logging"
)

// BulkService provides bulk actions for agents and organizations.
type BulkService interface {
	// BulkAgentAction performs a bulk action on agents.
	BulkAgentAction(ctx context.Context, action string, agentIDs []string, payload interface{}) (map[string]interface{}, error)
	// BulkOrgAction performs a bulk action on organizations.
	BulkOrgAction(ctx context.Context, action string, orgIDs []string, payload interface{}) (map[string]interface{}, error)
}

// AuditLogService provides access to audit logs (typed).
type AuditLogService interface {
	List(ctx context.Context, params *AuditLogListParams) ([]*AuditLog, *Pagination, error)
	Get(ctx context.Context, id string) (*AuditLog, error)
}

// IntegrationService provides access to integrations.
type IntegrationService interface {
	List(ctx context.Context, params map[string]string) ([]map[string]interface{}, error)
	Get(ctx context.Context, id string) (map[string]interface{}, error)
	Create(ctx context.Context, integration map[string]interface{}) (map[string]interface{}, error)
	Update(ctx context.Context, id string, integration map[string]interface{}) (map[string]interface{}, error)
	Delete(ctx context.Context, id string) error
}

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

	cache *Cache // Optional: in-memory cache for GET requests

	// Services for interacting with different API parts
	Account      AccountService
	Agent        AgentService
	Organization OrganizationService
	Incident     IncidentService
	Report       ReportService
	Billing      BillingService
	Webhook      WebhookService
	AuditLog     AuditLogService
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

	// Enable response caching for GET requests if requested
	if options.cacheTTL > 0 {
		client.cache = NewCache(options.cacheTTL)
	}

	// Initialize services

	client.Account = &accountService{client: client}
	client.Agent = &agentService{client: client}
	client.Organization = &organizationService{client: client}
	client.Incident = &incidentService{client: client}
	client.Report = &reportService{client: client}
	client.Billing = &billingService{client: client}
	client.Webhook = NewWebhookService(client)

	// Wire up audit log service
	auditlogRepo := newInternalAuditLogRepo(client)
	client.AuditLog = &auditLogService{repo: auditlogRepo}

	return client
}

// NewRequest creates a new API request
func (c *Client) NewRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	url := c.baseURL + path
	var bodyReader io.Reader
	if body != nil {
		if rdr, ok := body.(io.Reader); ok {
			bodyReader = rdr
		} else {
			jsonBody, err := json.Marshal(body)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal request body: %w", err)
			}
			bodyReader = bytes.NewReader(jsonBody)
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	const jsonMime = "application/json"
	req.Header.Set("Content-Type", jsonMime)
	req.Header.Set("Accept", jsonMime)
	req.Header.Set("Authorization", "Basic "+basicAuth(c.apiKey, c.apiSecret))
	req.Header.Set("User-Agent", c.userAgent)

	if c.Logger != nil {
		c.Logger.Debug("Creating new request", logging.String("method", method), logging.String("url", url))
	}

	return req, nil
}

// Do sends an API request and returns the response.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	// Response caching for GET requests
	if c.cache != nil && req.Method == http.MethodGet && v != nil {
		key := CacheKey(req)
		if cached := c.cache.Get(key); cached != nil {
			if err := json.Unmarshal(cached, v); err == nil {
				if c.Logger != nil {
					c.Logger.Debug("Cache hit", logging.String("url", req.URL.String()))
				}
				return nil, fmt.Errorf("response served from cache") // No HTTP response, but data is filled
			}
		}
	}
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
		errClose := resp.Body.Close()
		if errClose != nil {
			return nil, fmt.Errorf("client doJSON: error closing response body: %w", errClose)
		}
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

	// Read and cache the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		if c.Logger != nil {
			c.Logger.Error("Error reading response body", logging.Error("error", err))
		}
		return resp, fmt.Errorf("error reading response body: %w", err)
	}
	if err := json.Unmarshal(bodyBytes, v); err != nil {
		if c.Logger != nil {
			c.Logger.Error("Error decoding response", logging.Error("error", err))
		}
		return resp, fmt.Errorf("error decoding response: %w", err)
	}
	// Store in cache
	if c.cache != nil && req.Method == http.MethodGet {
		key := CacheKey(req)
		c.cache.Set(key, bodyBytes)
	}
	return resp, nil
}

// basicAuth creates a basic auth header value from credentials
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
