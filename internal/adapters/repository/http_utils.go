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

	"github.com/greysquirr3l/bishoujo-huntress/internal/ports/repository"
)

// HTTPClient defines the interface for an HTTP client
// We use an interface instead of a concrete type to make testing easier
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// createRequest creates a new HTTP request with authentication and standard headers
func createRequest(ctx context.Context, method, url string, body []byte, headers map[string]string) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set standard headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "bishoujo-huntress-client/1.0")

	// Add auth and custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

// handleErrorResponse extracts error information from an HTTP response
func handleErrorResponse(resp *http.Response) error {
	var errorResponse struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Details any    `json:"details,omitempty"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read error response: %w (status code: %d)", err, resp.StatusCode)
	}

	// Reset the body for potential future readers
	resp.Body = io.NopCloser(bytes.NewReader(body))

	// Try to parse the error response as JSON
	if err := json.Unmarshal(body, &errorResponse); err != nil {
		// If we can't parse it as JSON, return a generic error with the status code
		return fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	// If we successfully parsed the error, return a more detailed error
	if errorResponse.Message != "" {
		return fmt.Errorf("API error %d (%s): %s", resp.StatusCode, errorResponse.Code, errorResponse.Message)
	}

	// Fallback for unexpected error format
	return fmt.Errorf("request failed with status code %d", resp.StatusCode)
}

// extractPagination extracts pagination information from HTTP headers
func extractPagination(headers http.Header) repository.Pagination {
	pagination := repository.Pagination{
		Page:  1,
		Limit: 10, // Default values
	}

	// Parse page
	if page := headers.Get("X-Page"); page != "" {
		if pageNum, err := strconv.Atoi(page); err == nil && pageNum > 0 {
			pagination.Page = pageNum
		}
	}

	// Parse limit (per page)
	if limit := headers.Get("X-Per-Page"); limit != "" {
		if limitNum, err := strconv.Atoi(limit); err == nil && limitNum > 0 {
			pagination.Limit = limitNum
		}
	}

	// Parse total items
	if total := headers.Get("X-Total-Items"); total != "" {
		if totalNum, err := strconv.Atoi(total); err == nil {
			pagination.TotalItems = totalNum
		}
	}

	// Parse total pages
	if pages := headers.Get("X-Total-Pages"); pages != "" {
		if pagesNum, err := strconv.Atoi(pages); err == nil {
			pagination.TotalPages = pagesNum
		}
	}

	// Calculate has next/prev
	pagination.HasNext = pagination.Page < pagination.TotalPages
	pagination.HasPrev = pagination.Page > 1

	return pagination
}

// parseTime parses a time string in RFC3339 format
func parseTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, timeStr)
}
