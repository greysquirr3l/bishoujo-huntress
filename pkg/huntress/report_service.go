// Package huntress provides a client for the Huntress API
package huntress

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// reportService implements the ReportService interface
type reportService struct {
	client *Client
}

// Generate generates a report
func (s *reportService) Generate(ctx context.Context, input *ReportGenerateInput) (*Report, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/reports", input)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for Generate: %w", err)
	}

	report := new(Report)
	resp, err := s.client.Do(ctx, req, report)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for Generate: %w", err)
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	return report, nil
}

// Get retrieves report details by ID
func (s *reportService) Get(ctx context.Context, id string) (*Report, error) {
	path := fmt.Sprintf("/reports/%s", id)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for Get: %w", err)
	}

	report := new(Report)
	resp, err := s.client.Do(ctx, req, report)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for Get: %w", err)
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	return report, nil
}

// List lists reports with optional filtering
func (s *reportService) List(ctx context.Context, opts *ReportListOptions) ([]*Report, *Pagination, error) {
	// Advanced filtering: convert opts to map[string]interface{} using correct types
	filters := map[string]interface{}{}
	if opts != nil {
		if opts.Page > 0 {
			filters["page"] = opts.Page
		}
		if opts.PerPage > 0 {
			filters["per_page"] = opts.PerPage
		}
		if opts.Type != "" {
			filters["type"] = opts.Type
		}
		if opts.Format != "" {
			filters["format"] = opts.Format
		}
		if opts.Status != "" {
			filters["status"] = opts.Status
		}
		if opts.OrganizationID != "" {
			filters["organization_id"] = opts.OrganizationID
		}
		if opts.CreatedAfter != nil {
			filters["created_after"] = opts.CreatedAfter.Format(time.RFC3339)
		}
		if opts.CreatedBefore != nil {
			filters["created_before"] = opts.CreatedBefore.Format(time.RFC3339)
		}
	}
	var reports []*Report
	pagination, err := listResource(ctx, s.client, "/reports", filters, &reports)
	if err != nil {
		return nil, nil, err
	}
	return reports, pagination, nil
}

// Download downloads a report
func (s *reportService) Download(ctx context.Context, id string, format string) ([]byte, error) {
	path := fmt.Sprintf("/reports/%s/download", id)
	if format != "" {
		path = fmt.Sprintf("%s?format=%s", path, format)
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for Download: %w", err)
	}

	// For downloads, we need to handle the raw response
	if s.client.httpClient == nil {
		return nil, fmt.Errorf("HTTP client is not initialized")
	}

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request for Download: %w", err)
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}

	// Check for errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("api error: status code %d", resp.StatusCode)
	}

	// Read the full response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body in Download: %w", err)
	}
	return data, nil
}

// GetSummary retrieves a summary report
func (s *reportService) GetSummary(ctx context.Context, params *ReportParams) (*SummaryReport, error) {
	path := "/reports/summary"
	if params != nil {
		query, err := addQueryParams(path, params)
		if err != nil {
			return nil, fmt.Errorf("failed to add query params in GetSummary: %w", err)
		}
		path = query
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetSummary: %w", err)
	}

	report := new(SummaryReport)
	resp, err := s.client.Do(ctx, req, report)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for GetSummary: %w", err)
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	return report, nil
}

// GetDetails retrieves a detailed report
func (s *reportService) GetDetails(ctx context.Context, params *ReportParams) (*DetailedReport, error) {
	path := "/reports/detailed"
	if params != nil {
		query, err := addQueryParams(path, params)
		if err != nil {
			return nil, fmt.Errorf("failed to add query params in GetDetails: %w", err)
		}
		path = query
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetDetails: %w", err)
	}

	report := new(DetailedReport)
	resp, err := s.client.Do(ctx, req, report)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for GetDetails: %w", err)
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	return report, nil
}

// Schedule schedules a report for delivery
func (s *reportService) Schedule(ctx context.Context, params *ReportScheduleParams) (*ReportSchedule, error) {
	path := "/reports/schedule"
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for Schedule: %w", err)
	}

	schedule := new(ReportSchedule)
	resp, err := s.client.Do(ctx, req, schedule)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for Schedule: %w", err)
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	return schedule, nil
}

// Export exports a report in the specified format
func (s *reportService) Export(ctx context.Context, params *ReportExportParams) ([]byte, error) {
	path := fmt.Sprintf("/reports/%s/export", params.ReportID)
	if params.Format != "" {
		path = fmt.Sprintf("%s?format=%s", path, params.Format)
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for Export: %w", err)
	}

	// Set Accept header based on format
	if params.Format != "" {
		switch params.Format {
		case "pdf":
			req.Header.Set("Accept", "application/pdf")
		case "csv":
			req.Header.Set("Accept", "text/csv")
		case "json":
			req.Header.Set("Accept", "application/json")
		}
	}

	// For exports, we need to handle the raw response
	if s.client.httpClient == nil {
		return nil, fmt.Errorf("HTTP client is not initialized")
	}

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request for Export: %w", err)
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}

	// Check for errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("api error: status code %d", resp.StatusCode)
	}

	// Read the full response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body in Export: %w", err)
	}
	return data, nil
}
