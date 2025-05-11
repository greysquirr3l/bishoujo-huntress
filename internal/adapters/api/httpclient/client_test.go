package httpclient

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockRetrier struct {
	maxRetries  int
	shouldRetry bool
	backoff     time.Duration
}

func (m *mockRetrier) MaxRetries() int                      { return m.maxRetries }
func (m *mockRetrier) IsRetryableStatusCode(_ int) bool     { return m.shouldRetry }
func (m *mockRetrier) CalculateBackoff(_ int) time.Duration { return m.backoff }

func TestClient_Do_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		if _, err := w.Write([]byte(`{"ok":true}`)); err != nil {
			t.Fatalf("error writing response: %v", err)
		}
	}))
	defer srv.Close()
	client := New(2*time.Second, nil, nil)
	req, _ := http.NewRequestWithContext(context.Background(), "GET", srv.URL, nil)
	resp, err := client.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
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

func TestClient_Do_RateLimit(t *testing.T) {
	r := NewRateLimiter(0, time.Second) // always rate limited
	client := New(2*time.Second, r, nil)
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "http://example.com", nil)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	resp, err := client.Do(ctx, req)
	if resp != nil && resp.Body != nil {
		if cerr := resp.Body.Close(); cerr != nil {
			t.Errorf("error closing response body: %v", cerr)
		}
	}
	if err == nil {
		t.Errorf("expected error due to rate limiting, got nil")
	}
	if !errors.Is(err, &RateLimitError{}) && !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected RateLimitError or context.DeadlineExceeded, got %v", err)
	}
}

func TestClient_Do_Retry(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(500)
		} else {
			w.WriteHeader(200)
		}
	}))
	defer srv.Close()
	retrier := &mockRetrier{maxRetries: 2, shouldRetry: true, backoff: 10 * time.Millisecond}
	client := New(2*time.Second, nil, retrier)
	req, _ := http.NewRequestWithContext(context.Background(), "GET", srv.URL, nil)
	resp, err := client.Do(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestClient_Do_ContextCancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(200)
	}))
	defer srv.Close()
	retrier := &mockRetrier{maxRetries: 2, shouldRetry: true, backoff: 50 * time.Millisecond}
	client := New(2*time.Second, nil, retrier)
	req, _ := http.NewRequestWithContext(context.Background(), "GET", srv.URL, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	resp, err := client.Do(ctx, req)
	if resp != nil && resp.Body != nil {
		if cerr := resp.Body.Close(); cerr != nil {
			t.Errorf("error closing response body: %v", cerr)
		}
	}
	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context deadline exceeded, got %v", err)
	}
}

func TestClient_DoJSON_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		if _, err := w.Write([]byte(`{"foo":"bar"}`)); err != nil {
			t.Fatalf("error writing response: %v", err)
		}
	}))
	defer srv.Close()
	client := New(2*time.Second, nil, nil)
	req, _ := http.NewRequestWithContext(context.Background(), "GET", srv.URL, nil)
	var v map[string]string
	resp, err := client.DoJSON(context.Background(), req, &v)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}
	if v["foo"] != "bar" {
		t.Errorf("expected foo=bar, got %v", v)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestClient_DoJSON_ErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(400)
		if _, err := w.Write([]byte(`{"error":"bad request"}`)); err != nil {
			t.Fatalf("error writing response: %v", err)
		}
	}))
	defer srv.Close()
	client := New(2*time.Second, nil, nil)
	req, _ := http.NewRequestWithContext(context.Background(), "GET", srv.URL, nil)
	var v map[string]string
	resp, err := client.DoJSON(context.Background(), req, &v)
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}
	if err == nil || resp.StatusCode != 400 {
		t.Errorf("expected error and 400, got %v, %d", err, resp.StatusCode)
	}
}

func Test_decodeJSON_Error(t *testing.T) {
	bad := bytes.NewBufferString("notjson")
	var v map[string]interface{}
	err := decodeJSON(bad, &v)
	if err == nil {
		t.Error("expected error decoding bad json")
	}
}

func TestRateLimitError_Error(t *testing.T) {
	err := &RateLimitError{Message: "foo"}
	if err.Error() != "foo" {
		t.Errorf("expected 'foo', got '%s'", err.Error())
	}
}
