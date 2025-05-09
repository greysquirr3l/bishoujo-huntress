// Package api provides the main API client adapter for Huntress.
package api

import (
	"context"
	"fmt"
	"net/http"
)

// Client defines the interface for making API requests.
type Client interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

// DefaultClient is a basic implementation of Client using http.Client.
type DefaultClient struct {
	HTTPClient *http.Client
}

// Do executes the HTTP request using the underlying http.Client.
// The ctx parameter is currently unused.
func (c *DefaultClient) Do(_ context.Context, req *http.Request) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client do: %w", err)
	}
	return resp, nil
}
