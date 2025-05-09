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
		return &http.Response{StatusCode: 200}, nil
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
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
			return &http.Response{StatusCode: http.StatusServiceUnavailable}, nil
		}
		return &http.Response{StatusCode: 200}, nil
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
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
	_, err := r.Do(ctx, func() (*http.Response, error) {
		time.Sleep(10 * time.Millisecond)
		return &http.Response{StatusCode: 200}, nil
	})
	if err == nil {
		t.Fatal("expected error due to context timeout")
	}
	if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
		t.Errorf("expected context error, got %v", err)
	}
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
