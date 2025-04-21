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
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
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
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for errors
	if err := s.client.checkResponse(resp); err != nil {
		return nil, err
	}

	// Read the full response body
	return io.ReadAll(resp.Body)
}
