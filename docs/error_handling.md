# Error Handling

## Error Types

### APIError

```go
type APIError struct {
    StatusCode int                    `json:"status_code"`
    Message    string                 `json:"message"`
    Code       string                 `json:"code,omitempty"`
    Details    map[string]interface{} `json:"details,omitempty"`
    RequestID  string                 `json:"request_id,omitempty"`
    Timestamp  time.Time              `json:"timestamp"`
}

func (e *APIError) Error() string {
    if e.Code != "" {
        return fmt.Sprintf("API error %d (%s): %s", e.StatusCode, e.Code, e.Message)
    }
    return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

func (e *APIError) Is(target error) bool {
    var apiErr *APIError
    if errors.As(target, &apiErr) {
        return e.StatusCode == apiErr.StatusCode && e.Code == apiErr.Code
    }
    return false
}
```

### ValidationError

```go
type ValidationError struct {
    Field   string      `json:"field"`
    Value   interface{} `json:"value"`
    Tag     string      `json:"tag"`
    Message string      `json:"message"`
}

type ValidationErrors struct {
    Errors []ValidationError `json:"errors"`
}

func (e *ValidationErrors) Error() string {
    if len(e.Errors) == 1 {
        return fmt.Sprintf("validation failed for field '%s': %s", e.Errors[0].Field, e.Errors[0].Message)
    }
    return fmt.Sprintf("validation failed for %d fields", len(e.Errors))
}
```

### NetworkError

```go
type NetworkError struct {
    Operation string `json:"operation"`
    URL       string `json:"url"`
    Cause     error  `json:"-"`
}

func (e *NetworkError) Error() string {
    return fmt.Sprintf("network error during %s to %s: %v", e.Operation, e.URL, e.Cause)
}

func (e *NetworkError) Unwrap() error {
    return e.Cause
}
```

### AuthenticationError

```go
type AuthenticationError struct {
    Message string `json:"message"`
    Code    string `json:"code"`
}

func (e *AuthenticationError) Error() string {
    return fmt.Sprintf("authentication failed: %s", e.Message)
}
```

### RateLimitError

```go
type RateLimitError struct {
    RetryAfter time.Duration `json:"retry_after"`
    Limit      int           `json:"limit"`
    Remaining  int           `json:"remaining"`
    ResetAt    time.Time     `json:"reset_at"`
}

func (e *RateLimitError) Error() string {
    return fmt.Sprintf("rate limit exceeded, retry after %v", e.RetryAfter)
}
```

## HTTP Status Code Mappings

### 4xx Client Errors

- **400 Bad Request**: Invalid input parameters, malformed JSON

  ```go
  &APIError{StatusCode: 400, Code: "INVALID_REQUEST", Message: "Invalid input parameters"}
  ```

- **401 Unauthorized**: Invalid API credentials, expired token

  ```go
  &AuthenticationError{Message: "Invalid API credentials", Code: "INVALID_CREDENTIALS"}
  ```

- **403 Forbidden**: Insufficient permissions for requested resource

  ```go
  &APIError{StatusCode: 403, Code: "INSUFFICIENT_PERMISSIONS", Message: "Access denied"}
  ```

- **404 Not Found**: Resource does not exist

  ```go
  &APIError{StatusCode: 404, Code: "RESOURCE_NOT_FOUND", Message: "Resource not found"}
  ```

- **409 Conflict**: Resource conflict (e.g., duplicate email)

  ```go
  &APIError{StatusCode: 409, Code: "RESOURCE_CONFLICT", Message: "Email already exists"}
  ```

- **422 Unprocessable Entity**: Validation errors

  ```go
  &ValidationErrors{Errors: []ValidationError{...}}
  ```

- **429 Too Many Requests**: Rate limit exceeded

  ```go
  &RateLimitError{RetryAfter: 60*time.Second, Limit: 1000, Remaining: 0}
  ```

### 5xx Server Errors

- **500 Internal Server Error**: Unexpected server error

  ```go
  &APIError{StatusCode: 500, Code: "INTERNAL_ERROR", Message: "Internal server error"}
  ```

- **502 Bad Gateway**: Upstream service error

  ```go
  &APIError{StatusCode: 502, Code: "UPSTREAM_ERROR", Message: "Upstream service unavailable"}
  ```

- **503 Service Unavailable**: Service temporarily unavailable

  ```go
  &APIError{StatusCode: 503, Code: "SERVICE_UNAVAILABLE", Message: "Service temporarily unavailable"}
  ```

- **504 Gateway Timeout**: Request timeout

  ```go
  &APIError{StatusCode: 504, Code: "TIMEOUT", Message: "Request timeout"}
  ```

## Error Handling Patterns

### Basic Error Handling

```go
package main

import (
    "context"
    "errors"
    "log"
    "time"

    "github.com/greysquirr3l/bishoujo-huntress/pkg/huntress"
    "github.com/greysquirr3l/bishoujo-huntress/pkg/huntress/organization"
)

func main() {
    client := huntress.New(
        huntress.WithCredentials("API_KEY", "API_SECRET"),
        huntress.WithTimeout(30*time.Second),
    )

    ctx := context.Background()
    params := &organization.ListParams{Page: 1, Limit: 50}

    orgs, pagination, err := client.Organization.List(ctx, params)
    if err != nil {
        handleError(err)
        return
    }

    log.Printf("Found %d organizations", len(orgs))
}

func handleError(err error) {
    // Check for specific error types
    var apiErr *huntress.APIError
    var authErr *huntress.AuthenticationError
    var rateLimitErr *huntress.RateLimitError
    var validationErr *huntress.ValidationErrors
    var networkErr *huntress.NetworkError

    switch {
    case errors.As(err, &authErr):
        log.Printf("Authentication failed: %s", authErr.Message)
        // Handle authentication error (e.g., refresh credentials)

    case errors.As(err, &rateLimitErr):
        log.Printf("Rate limit exceeded, retry after %v", rateLimitErr.RetryAfter)
        // Handle rate limiting (e.g., implement exponential backoff)

    case errors.As(err, &validationErr):
        log.Printf("Validation failed:")
        for _, valErr := range validationErr.Errors {
            log.Printf("  Field %s: %s", valErr.Field, valErr.Message)
        }

    case errors.As(err, &apiErr):
        switch apiErr.StatusCode {
        case 400:
            log.Printf("Bad request: %s", apiErr.Message)
        case 403:
            log.Printf("Access denied: %s", apiErr.Message)
        case 404:
            log.Printf("Resource not found: %s", apiErr.Message)
        case 409:
            log.Printf("Conflict: %s", apiErr.Message)
        default:
            log.Printf("API error %d: %s", apiErr.StatusCode, apiErr.Message)
        }

    case errors.As(err, &networkErr):
        log.Printf("Network error: %s", networkErr.Error())
        // Handle network issues (e.g., retry with backoff)

    default:
        log.Printf("Unexpected error: %v", err)
    }
}
```

### Advanced Error Handling with Retry Logic

```go
package main

import (
    "context"
    "errors"
    "math"
    "time"

    "github.com/greysquirr3l/bishoujo-huntress/pkg/huntress"
)

type RetryConfig struct {
    MaxRetries      int
    InitialDelay    time.Duration
    MaxDelay        time.Duration
    BackoffMultiplier float64
}

func WithRetry(ctx context.Context, config RetryConfig, operation func() error) error {
    var lastErr error

    for attempt := 0; attempt <= config.MaxRetries; attempt++ {
        if attempt > 0 {
            delay := calculateDelay(attempt, config)
            timer := time.NewTimer(delay)
            defer timer.Stop()

            select {
            case <-ctx.Done():
                return ctx.Err()
            case <-timer.C:
                // Continue with retry
            }
        }

        err := operation()
        if err == nil {
            return nil // Success
        }

        lastErr = err

        // Check if error is retryable
        if !isRetryableError(err) {
            return err
        }

        // Handle rate limiting with server-specified delay
        var rateLimitErr *huntress.RateLimitError
        if errors.As(err, &rateLimitErr) {
            timer := time.NewTimer(rateLimitErr.RetryAfter)
            defer timer.Stop()

            select {
            case <-ctx.Done():
                return ctx.Err()
            case <-timer.C:
                // Continue with retry after rate limit delay
            }
        }
    }

    return lastErr
}

func calculateDelay(attempt int, config RetryConfig) time.Duration {
    delay := float64(config.InitialDelay) * math.Pow(config.BackoffMultiplier, float64(attempt-1))
    if delay > float64(config.MaxDelay) {
        delay = float64(config.MaxDelay)
    }
    return time.Duration(delay)
}

func isRetryableError(err error) bool {
    var apiErr *huntress.APIError
    var networkErr *huntress.NetworkError
    var rateLimitErr *huntress.RateLimitError

    switch {
    case errors.As(err, &rateLimitErr):
        return true // Always retry rate limit errors
    case errors.As(err, &networkErr):
        return true // Retry network errors
    case errors.As(err, &apiErr):
        // Retry 5xx server errors, but not 4xx client errors
        return apiErr.StatusCode >= 500
    default:
        return false
    }
}

// Example usage
func fetchOrganizationsWithRetry(client *huntress.Client) error {
    config := RetryConfig{
        MaxRetries:        3,
        InitialDelay:      1 * time.Second,
        MaxDelay:          30 * time.Second,
        BackoffMultiplier: 2.0,
    }

    ctx := context.Background()

    return WithRetry(ctx, config, func() error {
        _, _, err := client.Organization.List(ctx, &organization.ListParams{
            Page:  1,
            Limit: 50,
        })
        return err
    })
}
```

### Error Logging and Monitoring

```go
package main

import (
    "context"
    "encoding/json"
    "log/slog"
    "os"

    "github.com/greysquirr3l/bishoujo-huntress/pkg/huntress"
)

type ErrorLogger struct {
    logger *slog.Logger
}

func NewErrorLogger() *ErrorLogger {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))

    return &ErrorLogger{logger: logger}
}

func (e *ErrorLogger) LogError(ctx context.Context, operation string, err error) {
    attrs := []slog.Attr{
        slog.String("operation", operation),
        slog.String("error", err.Error()),
    }

    var apiErr *huntress.APIError
    var authErr *huntress.AuthenticationError
    var rateLimitErr *huntress.RateLimitError
    var validationErr *huntress.ValidationErrors
    var networkErr *huntress.NetworkError

    switch {
    case errors.As(err, &apiErr):
        attrs = append(attrs,
            slog.String("error_type", "api_error"),
            slog.Int("status_code", apiErr.StatusCode),
            slog.String("error_code", apiErr.Code),
            slog.String("request_id", apiErr.RequestID),
        )

        if apiErr.Details != nil {
            if details, err := json.Marshal(apiErr.Details); err == nil {
                attrs = append(attrs, slog.String("details", string(details)))
            }
        }

    case errors.As(err, &authErr):
        attrs = append(attrs,
            slog.String("error_type", "authentication_error"),
            slog.String("error_code", authErr.Code),
        )

    case errors.As(err, &rateLimitErr):
        attrs = append(attrs,
            slog.String("error_type", "rate_limit_error"),
            slog.Duration("retry_after", rateLimitErr.RetryAfter),
            slog.Int("limit", rateLimitErr.Limit),
            slog.Int("remaining", rateLimitErr.Remaining),
        )

    case errors.As(err, &validationErr):
        attrs = append(attrs,
            slog.String("error_type", "validation_error"),
            slog.Int("field_count", len(validationErr.Errors)),
        )

        for i, valErr := range validationErr.Errors {
            attrs = append(attrs,
                slog.String(fmt.Sprintf("field_%d", i), valErr.Field),
                slog.String(fmt.Sprintf("message_%d", i), valErr.Message),
            )
        }

    case errors.As(err, &networkErr):
        attrs = append(attrs,
            slog.String("error_type", "network_error"),
            slog.String("operation", networkErr.Operation),
            slog.String("url", networkErr.URL),
        )

    default:
        attrs = append(attrs, slog.String("error_type", "unknown"))
    }

    e.logger.LogAttrs(ctx, slog.LevelError, "API operation failed", attrs...)
}

// Example usage
func main() {
    client := huntress.New(
        huntress.WithCredentials("API_KEY", "API_SECRET"),
    )

    logger := NewErrorLogger()
    ctx := context.Background()

    _, _, err := client.Organization.List(ctx, &organization.ListParams{
        Page:  1,
        Limit: 50,
    })

    if err != nil {
        logger.LogError(ctx, "list_organizations", err)
    }
}
```

## Error Recovery Strategies

### Circuit Breaker Pattern

```go
type CircuitBreaker struct {
    maxFailures int
    resetTimeout time.Duration
    failures    int
    lastFailure time.Time
    state       CircuitState
    mutex       sync.RWMutex
}

type CircuitState int

const (
    Closed CircuitState = iota
    Open
    HalfOpen
)

func (cb *CircuitBreaker) Call(operation func() error) error {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()

    if cb.state == Open {
        if time.Since(cb.lastFailure) > cb.resetTimeout {
            cb.state = HalfOpen
            cb.failures = 0
        } else {
            return errors.New("circuit breaker is open")
        }
    }

    err := operation()
    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()

        if cb.failures >= cb.maxFailures {
            cb.state = Open
        }

        return err
    }

    // Success - reset circuit breaker
    cb.failures = 0
    cb.state = Closed
    return nil
}
```

### Graceful Degradation

```go
type FallbackClient struct {
    primary   *huntress.Client
    secondary *huntress.Client // Backup client or cached data
    cache     *Cache
}

func (f *FallbackClient) GetOrganization(ctx context.Context, id int) (*organization.Organization, error) {
    // Try primary client first
    org, err := f.primary.Organization.Get(ctx, id)
    if err == nil {
        // Cache successful response
        f.cache.Set(fmt.Sprintf("org_%d", id), org, 5*time.Minute)
        return org, nil
    }

    // Check if error is retryable
    if isRetryableError(err) {
        // Try secondary client
        if f.secondary != nil {
            if org, err := f.secondary.Organization.Get(ctx, id); err == nil {
                return org, nil
            }
        }
    }

    // Fall back to cache
    if cached := f.cache.Get(fmt.Sprintf("org_%d", id)); cached != nil {
        if org, ok := cached.(*organization.Organization); ok {
            return org, nil
        }
    }

    return nil, err
}
```

## Testing Error Conditions

```go
package main

import (
    "context"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/greysquirr3l/bishoujo-huntress/pkg/huntress"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestAPIErrorHandling(t *testing.T) {
    tests := []struct {
        name           string
        statusCode     int
        responseBody   string
        expectedError  string
        expectedType   interface{}
    }{
        {
            name:         "401 Unauthorized",
            statusCode:   401,
            responseBody: `{"message": "Invalid credentials", "code": "INVALID_CREDENTIALS"}`,
            expectedError: "authentication failed: Invalid credentials",
            expectedType: &huntress.AuthenticationError{},
        },
        {
            name:         "429 Rate Limit",
            statusCode:   429,
            responseBody: `{"message": "Rate limit exceeded", "retry_after": 60}`,
            expectedError: "rate limit exceeded, retry after 1m0s",
            expectedType: &huntress.RateLimitError{},
        },
        {
            name:         "422 Validation Error",
            statusCode:   422,
            responseBody: `{"errors": [{"field": "email", "message": "Invalid email format"}]}`,
            expectedError: "validation failed for field 'email': Invalid email format",
            expectedType: &huntress.ValidationErrors{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(tt.statusCode)
                w.Write([]byte(tt.responseBody))
            }))
            defer server.Close()

            client := huntress.New(
                huntress.WithBaseURL(server.URL),
                huntress.WithCredentials("test", "test"),
            )

            _, _, err := client.Organization.List(context.Background(), nil)

            require.Error(t, err)
            assert.Equal(t, tt.expectedError, err.Error())
            assert.IsType(t, tt.expectedType, err)
        })
    }
}
```

This comprehensive error handling documentation provides developers with clear guidance on how to handle all types of errors that may occur when using the Bishoujo-Huntress API client, following Go best practices and the project's DDD architecture.
