package huntress

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestBillingService_GetSummary(t *testing.T) {
	respSummary := &BillingSummary{CurrentBalance: 42.0}
	body, _ := json.Marshal(respSummary)
	client := newTestClient(roundTripFunc(func(r *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header:     make(http.Header),
		}
	}))
	svc := &billingService{client: &Client{httpClient: client.Client.httpClient}}
	sum, err := svc.GetSummary(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum.CurrentBalance != 42.0 {
		t.Errorf("expected current balance 42.0, got %v", sum.CurrentBalance)
	}
}

func TestBillingService_GetSummary_HTTPError(t *testing.T) {
	client := newTestClient(roundTripFunc(func(r *http.Request) *http.Response {
		return &http.Response{StatusCode: 500, Body: io.NopCloser(bytes.NewReader([]byte("fail")))}
	}))
	svc := &billingService{client: &Client{httpClient: client.Client.httpClient}}
	_, err := svc.GetSummary(context.Background())
	if err == nil {
		t.Error("expected error for HTTP 500")
	}
}

// Additional tests for ListInvoices, GetInvoice, GetUsage would follow the same pattern.
