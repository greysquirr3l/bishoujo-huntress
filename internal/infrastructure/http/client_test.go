package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/greysquirr3l/bishoujo-huntress/internal/infrastructure/http/retry"
)

func newTestClient(handler http.Handler) *Client {
	ts := httptest.NewServer(handler)
	client, _ := NewClient(ts.URL, "test-key", "test-secret")
	client.HTTPClient = ts.Client()
	return client
}

func TestClient_Do_Success(t *testing.T) {
	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Request-Id", "req-123")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"ok":true}`)); err != nil {
			t.Errorf("error writing response: %v", err)
		}
	}))
	var result map[string]interface{}
	resp, err := client.Do(context.Background(), http.MethodGet, "/test", nil, &result, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if v, ok := result["ok"]; !ok || v != true {
		t.Errorf("expected result ok=true, got %v", result)
	}
}

func TestClient_Do_ErrorResponse(t *testing.T) {
	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Request-Id", "req-err")
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte(`{"message":"bad request"}`)); err != nil {
			t.Errorf("error writing response: %v", err)
		}
	}))
	var result map[string]interface{}
	resp, err := client.Do(context.Background(), http.MethodGet, "/fail", nil, &result, nil)
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", apiErr.StatusCode)
	}
	if apiErr.Message != "bad request" {
		t.Errorf("expected message 'bad request', got %q", apiErr.Message)
	}
	if apiErr.RequestID != "req-err" {
		t.Errorf("expected request id 'req-err', got %q", apiErr.RequestID)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestClient_Do_ContextCancel(t *testing.T) {
	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	var result map[string]interface{}
	resp, err := client.Do(ctx, http.MethodGet, "/timeout", nil, &result, nil)
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}
	if err == nil {
		t.Fatal("expected error due to context timeout")
	}
	if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
		t.Errorf("expected context error, got %v", err)
	}
}

func TestClient_Do_AuthHeader(t *testing.T) {
	var gotAuth string
	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"ok":true}`)); err != nil {
			t.Errorf("error writing response: %v", err)
		}
	}))
	var result map[string]interface{}
	resp, err := client.Do(context.Background(), http.MethodGet, "/auth", nil, &result, nil)
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if gotAuth == "" {
		t.Error("expected Authorization header to be set")
	}
}

func TestClient_Do_RetryLogic(t *testing.T) {
	calls := 0
	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		if calls < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			if _, err := w.Write([]byte(`{"message":"unavailable"}`)); err != nil {
				t.Errorf("error writing response: %v", err)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"ok":true}`)); err != nil {
			t.Errorf("error writing response: %v", err)
		}
	}))
	client.RetryConfig = &retry.Config{
		MaxRetries:           2,
		BaseDelay:            10 * time.Millisecond,
		MaxDelay:             50 * time.Millisecond,
		RetryableStatusCodes: []int{http.StatusServiceUnavailable},
	}
	var result map[string]interface{}
	resp, err := client.Do(context.Background(), http.MethodGet, "/retry", nil, &result, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestClient_Get_Post_Put_Delete(t *testing.T) {
	client := newTestClient(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"ok":true}`)); err != nil {
			t.Errorf("error writing response: %v", err)
		}
	}))
	var result map[string]interface{}
	resp, err := client.Get(context.Background(), "/get", &result, nil)
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	resp, err = client.Post(context.Background(), "/post", map[string]string{"foo": "bar"}, &result, nil)
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}
	if err != nil {
		t.Errorf("Post failed: %v", err)
	}
	resp, err = client.Put(context.Background(), "/put", map[string]string{"foo": "baz"}, &result, nil)
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}
	if err != nil {
		t.Errorf("Put failed: %v", err)
	}
	resp, err = client.Delete(context.Background(), "/delete", &result, nil)
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}
}

func TestGetPagination_Defaults(t *testing.T) {
	resp := &http.Response{Header: http.Header{}}
	pagination, err := GetPagination(resp)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pagination.CurrentPage != 1 || pagination.TotalPages != 1 {
		t.Errorf("expected default pagination, got %+v", pagination)
	}
}
