// Package api provides HTTP client utilities for Huntress API adapters.
package api

import (
	"context"
	"fmt"
	"net/http"
)

// HTTPClientAdapter wraps an http.Client for use in API adapters.
type HTTPClientAdapter struct {
	Client *http.Client
}

// Do executes the HTTP request using the underlying http.Client.
// The ctx parameter is currently unused.
func (a *HTTPClientAdapter) Do(_ context.Context, req *http.Request) (*http.Response, error) {
	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http client adapter do: %w", err)
	}
	return resp, nil
}
