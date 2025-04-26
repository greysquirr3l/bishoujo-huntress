package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/errors"
	"github.com/greysquirr3l/bishoujo-huntress/internal/ports/repository"
)

// HTTPClient defines the interface for an HTTP client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// createRequest creates an HTTP request with the provided parameters
func createRequest(ctx context.Context, method, url string, body []byte, headers map[string]string) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

// handleErrorResponse translates HTTP errors into domain errors
func handleErrorResponse(resp *http.Response) error {
	var apiError struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Details string `json:"details,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
		// If we can't decode the error, create a generic one
		return errors.NewAPIError(
			resp.StatusCode,
			"UNKNOWN_ERROR",
			fmt.Sprintf("HTTP error: %d", resp.StatusCode),
			"",
		)
	}

	// If the API returned an error code, use it
	if apiError.Code == "" {
		apiError.Code = "API_ERROR"
	}

	return errors.NewAPIError(
		resp.StatusCode,
		apiError.Code,
		apiError.Message,
		apiError.Details,
	)
}

// extractPagination extracts pagination information from response headers
func extractPagination(headers http.Header) repository.Pagination {
	pagination := repository.Pagination{
		Page:       1,
		PerPage:    20,
		TotalPages: 1,
		TotalItems: 0,
	}

	// Parse page number
	if page := headers.Get("X-Page"); page != "" {
		if val, err := strconv.Atoi(page); err == nil {
			pagination.Page = val
		}
	}

	// Parse per page
	if perPage := headers.Get("X-Per-Page"); perPage != "" {
		if val, err := strconv.Atoi(perPage); err == nil {
			pagination.PerPage = val
		}
	}

	// Parse total pages
	if totalPages := headers.Get("X-Total-Pages"); totalPages != "" {
		if val, err := strconv.Atoi(totalPages); err == nil {
			pagination.TotalPages = val
		}
	}

	// Parse total items
	if totalItems := headers.Get("X-Total-Count"); totalItems != "" {
		if val, err := strconv.Atoi(totalItems); err == nil {
			pagination.TotalItems = val
		}
	}

	return pagination
}

// parseTime parses a time string in RFC3339 format
func parseTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %w", err)
	}
	return t, nil
}
