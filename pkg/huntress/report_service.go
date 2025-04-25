// Package huntress provides a client for the Huntress API
package huntress

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// reportService implements the ReportService interface
type reportService struct {
	client *Client
}

// Generate generates a report
func (s *reportService) Generate(ctx context.Context, input *ReportGenerateInput) (*Report, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/reports", input)
	if err != nil {
		return nil, err
	}

	report := new(Report)
	_, err = s.client.Do(ctx, req, report)
	if err != nil {
		return nil, err
	}

	return report, nil
}

// Get retrieves report details by ID
func (s *reportService) Get(ctx context.Context, id string) (*Report, error) {
	path := fmt.Sprintf("/reports/%s", id)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	report := new(Report)
	_, err = s.client.Do(ctx, req, report)
	if err != nil {
		return nil, err
	}

	return report, nil
}

// List lists reports with optional filtering
func (s *reportService) List(ctx context.Context, opts *ReportListOptions) ([]*Report, *Pagination, error) {
	path := "/reports"
	if opts != nil {
		query, err := addQueryParams(path, opts)
		if err != nil {
			return nil, nil, err
		}
		path = query
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var reports []*Report
	resp, err := s.client.Do(ctx, req, &reports)
	if err != nil {
		return nil, nil, err
	}

	pagination := extractPagination(resp)
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
		return nil, err
	}

	// For downloads, we need to handle the raw response
	if s.client.httpClient == nil {
		return nil, fmt.Errorf("HTTP client is not initialized")
	}

	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("api error: status code %d", resp.StatusCode)
	}

	// Read the full response body
	return io.ReadAll(resp.Body)
}

// GetSummary retrieves a summary report
func (s *reportService) GetSummary(ctx context.Context, params *ReportParams) (*SummaryReport, error) {
	path := "/reports/summary"
	if params != nil {
		query, err := addQueryParams(path, params)
		if err != nil {
			return nil, err
		}
		path = query
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	report := new(SummaryReport)
	_, err = s.client.Do(ctx, req, report)
	if err != nil {
		return nil, err
	}

	return report, nil
}

// GetDetails retrieves a detailed report
func (s *reportService) GetDetails(ctx context.Context, params *ReportParams) (*DetailedReport, error) {
	path := "/reports/detailed"
	if params != nil {
		query, err := addQueryParams(path, params)
		if err != nil {
			return nil, err
		}
		path = query
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	report := new(DetailedReport)
	_, err = s.client.Do(ctx, req, report)
	if err != nil {
		return nil, err
	}

	return report, nil
}

// Schedule schedules a report for delivery
func (s *reportService) Schedule(ctx context.Context, params *ReportScheduleParams) (*ReportSchedule, error) {
	path := "/reports/schedule"
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return nil, err
	}

	schedule := new(ReportSchedule)
	_, err = s.client.Do(ctx, req, schedule)
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("api error: status code %d", resp.StatusCode)
	}

	// Read the full response body
	return io.ReadAll(resp.Body)
}
