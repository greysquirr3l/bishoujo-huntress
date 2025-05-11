package httpclient

import (
	"context"
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	rl := NewRateLimiter(2, 100*time.Millisecond)
	if !rl.Allow() {
		t.Error("expected first request to be allowed")
	}
	if !rl.Allow() {
		t.Error("expected second request to be allowed")
	}
	if rl.Allow() {
		t.Error("expected third request to be rate limited")
	}
	// Wait for window to expire
	time.Sleep(110 * time.Millisecond)
	if !rl.Allow() {
		t.Error("expected request to be allowed after window")
	}
}

func TestRateLimiter_Wait_AllowsEventually(t *testing.T) {
	rl := NewRateLimiter(1, 50*time.Millisecond)
	if !rl.Allow() {
		t.Fatal("expected first request to be allowed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	start := time.Now()
	err := rl.Wait(ctx)
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if elapsed < 45*time.Millisecond {
		t.Errorf("expected wait at least 45ms, got %v", elapsed)
	}
}

func TestRateLimiter_Wait_ContextCancel(t *testing.T) {
	rl := NewRateLimiter(0, time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	err := rl.Wait(ctx)
	if err == nil {
		t.Error("expected error on context cancel")
	}
}
