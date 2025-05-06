// Package httpclient provides a rate limiter for Huntress API adapters.
package httpclient

import (
	"context"
	"sync"
	"time"
)

// RateLimiter implements a sliding window rate limiter (60 requests/minute).
type RateLimiter struct {
	mu         sync.Mutex
	timestamps []time.Time
	limit      int
	window     time.Duration
}

// NewRateLimiter creates a new RateLimiter.
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:  limit,
		window: window,
	}
}

// Allow returns true if a request is allowed, false if rate limited.
func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-r.window)
	// Remove old timestamps
	var kept []time.Time
	for _, t := range r.timestamps {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	r.timestamps = kept
	if len(r.timestamps) < r.limit {
		r.timestamps = append(r.timestamps, now)
		return true
	}
	return false
}

// Wait blocks until a request is allowed or the context is done.
func (r *RateLimiter) Wait(ctx context.Context) error {
	for {
		if r.Allow() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}
}
