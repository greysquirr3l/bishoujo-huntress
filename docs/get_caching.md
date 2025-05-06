# Huntress Go Client: GET Response Caching

This document describes the in-memory caching layer for GET requests in the Bishoujo-Huntress Go client library.

## Overview

The Huntress Go client provides an optional, thread-safe in-memory cache for HTTP GET responses. This cache is designed to improve performance and reduce redundant API calls for frequently accessed resources.

## Features

- **In-memory, thread-safe**: Uses a mutex for safe concurrent access.
- **Configurable TTL**: Each cache entry expires after a configurable time-to-live (TTL).
- **Simple API**: Exposes `Get`, `Set`, and `CacheKey` methods.
- **GET-only**: Only responses to HTTP GET requests are cached.

## Usage

### Enabling the Cache

To use the cache, create a `Cache` instance and integrate it with your client logic:

```go
import "github.com/greysquirr3l/bishoujo-huntress/pkg/huntress"

// Create a cache with a 5-minute TTL
cache := huntress.NewCache(5 * time.Minute)

// Generate a cache key for a GET request
key := huntress.CacheKey(req)

// Try to get a cached response
if resp := cache.Get(key); resp != nil {
    // Use cached response
} else {
    // Perform the GET request, then cache the result
    cache.Set(key, responseBytes)
}
```

### API Reference

- `NewCache(ttl time.Duration) *Cache`
  Creates a new cache with the specified TTL for entries.

- `(*Cache) Get(key string) []byte`
  Returns the cached response for the given key, or `nil` if not found or expired.

- `(*Cache) Set(key string, response []byte)`
  Stores the response for the given key, with expiry set to now + TTL.

- `CacheKey(req *http.Request) string`
  Generates a cache key based on the HTTP method and URL.

## Notes & Limitations

- Only GET requests should be cached. Do not cache POST, PUT, PATCH, or DELETE responses.
- Expired entries are not automatically purged, but are ignored on access.
- The cache is in-memory and not persistent. It is cleared when the process exits.
- The cache is best suited for short-lived, frequently accessed data.

## Example

```go
cache := huntress.NewCache(2 * time.Minute)
key := huntress.CacheKey(req)
if data := cache.Get(key); data != nil {
    // Use cached data
} else {
    // Fetch from API, then cache
    cache.Set(key, apiResponse)
}
```

## When to Use

- To reduce API rate limit usage for frequently repeated GET requests.
- To improve performance for dashboards or UIs that poll the same data.
- For short-lived, non-sensitive data where eventual consistency is acceptable.

## When Not to Use

- For sensitive or highly dynamic data.
- When strong consistency is required.
- For long-term or persistent caching needs (use a distributed cache or database instead).

---

For more details, see [`pkg/huntress/cache.go`](../pkg/huntress/cache.go).
