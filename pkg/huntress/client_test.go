package huntress

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

const (
	fooKey = "foo"
	barVal = "bar"
)

// Use roundTripFunc from testhelpers_test.go (do not redeclare here)

func TestClient_New(t *testing.T) {
	client := New(
		WithCredentials("test-key", "test-secret"),
	)
	if client == nil {
		t.Fatal("expected client to be non-nil")
	}
}

func TestAccountService_Get_NotImplemented(t *testing.T) {
	client := New(WithCredentials("test", "test"))
	_, err := client.Account.Get(context.Background())
	if err == nil {
		t.Error("expected error for Get, got nil")
	}
}

// No aliasing needed; use package-level identifiers directly.

func TestClient_NewRequest_sets_headers(t *testing.T) {
	client := New(WithCredentials("foo", "bar"), WithUserAgent("myagent/1.2.3"))
	req, err := client.NewRequest(context.Background(), "GET", "/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Header.Get("Authorization") == "" {
		t.Error("expected Authorization header to be set")
	}
	if req.Header.Get("User-Agent") != "myagent/1.2.3" {
		t.Errorf("expected User-Agent 'myagent/1.2.3', got '%s'", req.Header.Get("User-Agent"))
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", req.Header.Get("Content-Type"))
	}
	if req.Header.Get("Accept") != "application/json" {
		t.Errorf("expected Accept 'application/json', got '%s'", req.Header.Get("Accept"))
	}
}

func TestClient_NewRequest_with_body(t *testing.T) {
	client := New(WithCredentials("foo", "bar"))
	type payload struct{ Foo string }
	req, err := client.NewRequest(context.Background(), "POST", "/test", payload{Foo: "bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("unexpected error reading body: %v", err)
	}
	if string(b) == "" || string(b) == "null" {
		t.Error("expected non-empty body")
	}
}

func TestClient_Do_success(t *testing.T) {
	respBody := map[string]string{"foo": "bar"}
	body, _ := json.Marshal(respBody)
	client := New(
		WithCredentials("foo", "bar"),
		WithHTTPClient(&http.Client{Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     make(http.Header),
			}
		})}),
	)
	req, _ := client.NewRequest(context.Background(), "GET", "/test", nil)
	var out map[string]string
	resp, err := client.Do(context.Background(), req, &out)
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[fooKey] != barVal {
		t.Errorf("expected foo=bar, got %v", out)
	}
}

func TestClient_Do_error_status(t *testing.T) {
	client := New(
		WithCredentials("foo", "bar"),
		WithHTTPClient(&http.Client{Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewReader([]byte("fail"))),
				Header:     make(http.Header),
			}
		})}),
	)
	req, _ := client.NewRequest(context.Background(), "GET", "/fail", nil)
	var out map[string]string
	resp, err := client.Do(context.Background(), req, &out)
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}
	if err == nil || err.Error() == "" {
		t.Error("expected error for HTTP 500")
	}
}

func TestClient_Do_cache(t *testing.T) {
	respBody := map[string]string{"foo": "bar"}
	body, _ := json.Marshal(respBody)
	called := 0
	client := New(
		WithCredentials("foo", "bar"),
		WithCacheTTL(1*time.Minute),
		WithHTTPClient(&http.Client{Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			called++
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(body)),
				Header:     make(http.Header),
			}
		})}),
	)
	req, _ := client.NewRequest(context.Background(), "GET", "/cache", nil)
	var out map[string]string
	resp, err := client.Do(context.Background(), req, &out)
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}
	if err != nil && err.Error() != "response served from cache" {
		t.Fatalf("unexpected error: %v", err)
	}
	// Second call should hit cache
	req2, _ := client.NewRequest(context.Background(), "GET", "/cache", nil)
	var out2 map[string]string
	resp2, err2 := client.Do(context.Background(), req2, &out2)
	if resp2 != nil && resp2.Body != nil {
		if err := resp2.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}
	if err2 == nil || err2.Error() != "response served from cache" {
		t.Errorf("expected cache error, got %v", err2)
	}
	if out2[fooKey] != barVal {
		t.Errorf("expected foo=bar from cache, got %v", out2)
	}
	if called != 1 {
		t.Errorf("expected http call only once, got %d", called)
	}
}
