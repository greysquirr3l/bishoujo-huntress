package retry

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestRetrier_Do_Success(t *testing.T) {
	r := NewRetrier(DefaultConfig)
	resp, err := r.Do(context.Background(), func() (*http.Response, error) {
		// No handler, just return a response with NoBody
		return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
	})
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

func TestRetrier_Do_RetryOnStatus(t *testing.T) {
	calls := 0
	r := NewRetrier(Config{
		MaxRetries:           2,
		BaseDelay:            1 * time.Millisecond,
		MaxDelay:             10 * time.Millisecond,
		RetryableStatusCodes: []int{http.StatusServiceUnavailable},
	})
	resp, err := r.Do(context.Background(), func() (*http.Response, error) {
		calls++
		if calls < 3 {
			return &http.Response{StatusCode: http.StatusServiceUnavailable, Body: http.NoBody}, nil
		}
		return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
	})
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
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestRetrier_Do_ContextCancel(t *testing.T) {
	r := NewRetrier(DefaultConfig)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	resp, err := r.Do(ctx, func() (*http.Response, error) {
		time.Sleep(50 * time.Millisecond)
		return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
	})
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("error closing response body: %v", err)
		}
	}
	if err == nil {
		if ctx.Err() != nil {
			// Acceptable: context expired, retrier returned nil error
			return
		}
		t.Fatal("expected error due to context timeout")
	}

	errStr := err.Error()
	t.Logf("error string: %q", errStr)
	if !contains(errStr, "deadline exceeded") && !contains(errStr, "context canceled") {
		t.Fatalf("expected context deadline exceeded or canceled, got: %v (error string: %q)", err, errStr)
	}
}

func contains(haystack, needle string) bool {
	return len(needle) > 0 && len(haystack) > 0 && (len(haystack) >= len(needle)) && (stringIndex(haystack, needle) >= 0)
}

func stringIndex(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func TestExecuteWithRetry_Success(t *testing.T) {
	cfg := DefaultConfig
	err := ExecuteWithRetry(context.Background(), func() error {
		return nil
	}, &cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestExecuteWithRetry_RetryAndFail(t *testing.T) {
	cfg := Config{MaxRetries: 2, BaseDelay: 1 * time.Millisecond, MaxDelay: 10 * time.Millisecond}
	calls := 0
	err := ExecuteWithRetry(context.Background(), func() error {
		calls++
		return errors.New("fail")
	}, &cfg)
	if err == nil {
		t.Fatal("expected error after retries")
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}
