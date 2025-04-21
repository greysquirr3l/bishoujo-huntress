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

// GetInvoice retrieves an invoice by ID
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

// ListInvoices lists invoices with optional filtering
func (s *billingService) ListInvoices(ctx context.Context, opts *InvoiceListOptions) ([]*Invoice, *Pagination, error) {
	path := "/billing/invoices"
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
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

// GetUsage retrieves usage information
func (s *billingService) GetUsage(ctx context.Context, period *BillingPeriod) (*UsageReport, error) {
	path := "/billing/usage"
	path, err := addOptions(path, period)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	usageReport := new(UsageReport)
	_, err = s.client.Do(ctx, req, usageReport)
	if err != nil {
		return nil, err
	}

	return usageReport, nil
}
