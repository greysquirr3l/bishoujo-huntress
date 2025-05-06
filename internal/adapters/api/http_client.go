// Package api provides HTTP client utilities for Huntress API adapters.
package api

import (
	"context"
	"net/http"
)

// HTTPClientAdapter wraps an http.Client for use in API adapters.
type HTTPClientAdapter struct {
	Client *http.Client
}

func (a *HTTPClientAdapter) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return a.Client.Do(req)
}
