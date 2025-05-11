package api

import (
	"context"
	"net/http"
	"testing"
)

type dummyRoundTripper struct{}

func (d *dummyRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
}

func TestHTTPClient_Do(t *testing.T) {
	client := &http.Client{Transport: &dummyRoundTripper{}}
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "https://api.example.com/test", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}
