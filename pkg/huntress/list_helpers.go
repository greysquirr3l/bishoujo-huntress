// Package huntress provides a client for the Huntress API
package huntress

import (
	"context"
	"fmt"
	"net/http"
	"os"
)

// listResource is a generic helper for paginated GET endpoints.
// It takes a path, params, a pointer to a slice for results, and a function to create a new request.
func listResource[T any](ctx context.Context, client *Client, path string, params interface{}, result *[]T) (*Pagination, error) {
	if params != nil {
		query, err := addQueryParams(path, params)
		if err != nil {
			return nil, fmt.Errorf("failed to add query params in List: %w", err)
		}
		path = query
	}

	req, err := client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for List: %w", err)
	}

	resp, err := client.Do(ctx, req, result)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for List: %w", err)
	}
	if resp != nil {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing response body: %v\n", err)
		}
	}

	pagination := extractPagination(resp)
	return pagination, nil
}
