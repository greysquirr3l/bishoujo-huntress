// Package huntress provides a client for the Huntress API
package huntress

import (
	"context"
	"fmt"
	"net/http"
)

// billingService implements the BillingService interface
type billingService struct {
	client *Client
}

// GetSummary retrieves a billing summary
func (s *billingService) GetSummary(ctx context.Context) (*BillingSummary, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/billing/summary", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetSummary: %w", err)
	}

	summary := new(BillingSummary)
	resp, err := s.client.Do(ctx, req, summary)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for GetSummary: %w", err)
	}
	if resp != nil {
		if err := resp.Body.Close(); err != nil {
			return nil, fmt.Errorf("billing get summary: error closing response body: %w", err)
		}
	}

	return summary, nil
}

// ListInvoices lists all invoices
func (s *billingService) ListInvoices(ctx context.Context, params *ListParams) ([]*Invoice, *Pagination, error) {
	var invoices []*Invoice
	pagination, err := listResource(ctx, s.client, "/billing/invoices", params, &invoices)
	if err != nil {
		return nil, nil, err
	}
	return invoices, pagination, nil
}

// GetInvoice retrieves a specific invoice
func (s *billingService) GetInvoice(ctx context.Context, id string) (*Invoice, error) {
	path := fmt.Sprintf("/billing/invoices/%s", id)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetInvoice: %w", err)
	}

	invoice := new(Invoice)
	resp, err := s.client.Do(ctx, req, invoice)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for GetInvoice: %w", err)
	}
	if resp != nil {
		if err := resp.Body.Close(); err != nil {
			return nil, fmt.Errorf("billing get invoice: error closing response body: %w", err)
		}
	}

	return invoice, nil
}

// GetUsage retrieves usage statistics
func (s *billingService) GetUsage(ctx context.Context, params *UsageParams) (*UsageReport, error) {
	path := "/billing/usage"
	if params != nil {
		query, err := addQueryParams(path, params)
		if err != nil {
			return nil, fmt.Errorf("failed to add query params in GetUsage: %w", err)
		}
		path = query
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetUsage: %w", err)
	}

	usage := new(UsageReport)
	resp, err := s.client.Do(ctx, req, usage)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for GetUsage: %w", err)
	}
	if resp != nil {
		if err := resp.Body.Close(); err != nil {
			return nil, fmt.Errorf("billing get usage: error closing response body: %w", err)
		}
	}

	return usage, nil
}
