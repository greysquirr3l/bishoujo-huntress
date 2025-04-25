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
		return nil, err
	}

	summary := new(BillingSummary)
	_, err = s.client.Do(ctx, req, summary)
	if err != nil {
		return nil, err
	}

	return summary, nil
}

// ListInvoices lists all invoices
func (s *billingService) ListInvoices(ctx context.Context, params *ListParams) ([]*Invoice, *Pagination, error) {
	path := "/billing/invoices"
	if params != nil {
		query, err := addQueryParams(path, params)
		if err != nil {
			return nil, nil, err
		}
		path = query
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var invoices []*Invoice
	resp, err := s.client.Do(ctx, req, &invoices)
	if err != nil {
		return nil, nil, err
	}

	pagination := extractPagination(resp)
	return invoices, pagination, nil
}

// GetInvoice retrieves a specific invoice
func (s *billingService) GetInvoice(ctx context.Context, id string) (*Invoice, error) {
	path := fmt.Sprintf("/billing/invoices/%s", id)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	invoice := new(Invoice)
	_, err = s.client.Do(ctx, req, invoice)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

// GetUsage retrieves usage statistics
func (s *billingService) GetUsage(ctx context.Context, params *UsageParams) (*UsageReport, error) {
	path := "/billing/usage"
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

	usage := new(UsageReport)
	_, err = s.client.Do(ctx, req, usage)
	if err != nil {
		return nil, err
	}

	return usage, nil
}
