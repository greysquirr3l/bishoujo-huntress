package api

import (
	"context"
	"net/http"
	"testing"
)

// roundTripFunc allows us to mock http.RoundTripper inline for tests.
type roundTripFunc func(_ *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(_ *http.Request) (*http.Response, error) {
	return f(nil), nil
}

func TestDefaultClient_Do_Success(t *testing.T) {
	dummy := &http.Client{Transport: roundTripFunc(func(_ *http.Request) *http.Response {
		return &http.Response{StatusCode: 200, Body: http.NoBody}
	})}
	c := &DefaultClient{HTTPClient: dummy}
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "https://api.example.com/test", nil)
	resp, err := c.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestDefaultClient_Do_Error(t *testing.T) {
	dummy := &http.Client{Transport: roundTripFunc(func(_ *http.Request) *http.Response {
		return nil
	})}
	c := &DefaultClient{HTTPClient: dummy}
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "https://api.example.com/test", nil)
	resp, err := c.Do(context.Background(), req)
	if resp != nil && resp.Body != nil {
		if cerr := resp.Body.Close(); cerr != nil {
			t.Errorf("error closing response body: %v", cerr)
		}
	}
	if err == nil {
		t.Error("expected error for nil response")
	}
}
