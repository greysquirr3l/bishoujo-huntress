// Package api provides the main API client adapter for Huntress.
package api

import (
	"context"
	"net/http"
)

// APIClient defines the interface for making API requests.
type APIClient interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

// DefaultAPIClient is a basic implementation of APIClient using http.Client.
type DefaultAPIClient struct {
	HTTPClient *http.Client
}

func (c *DefaultAPIClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return c.HTTPClient.Do(req)
}
