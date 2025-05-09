// Package huntress provides a simple in-memory cache for GET requests.
package huntress

import (
	"net/http"
	"sync"
	"time"
)

type cacheEntry struct {
	response []byte
	expiry   time.Time
}

// Cache is a simple in-memory cache for GET requests.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]*cacheEntry
	ttl     time.Duration
}

// NewCache creates a new Cache with the given TTL.
func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]*cacheEntry),
		ttl:     ttl,
	}
}

// Get returns the cached response for the given key, or nil if not found or expired.
func (c *Cache) Get(key string) []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	if !ok || time.Now().After(entry.expiry) {
		return nil
	}
	return entry.response
}

// Set stores the response for the given key.
func (c *Cache) Set(key string, response []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = &cacheEntry{
		response: response,
		expiry:   time.Now().Add(c.ttl),
	}
}

// CacheKey generates a cache key for a GET request.
func CacheKey(req *http.Request) string {
	return req.Method + ":" + req.URL.String()
}
