package huntress

import (
	"errors"
	"fmt"
	"net/http"
)

// Common errors returned by the Huntress API client
var (
	// ErrBadRequest represents a 400 Bad Request error
	ErrBadRequest = errors.New("bad request")
	// ErrUnauthorized represents a 401 Unauthorized error
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden represents a 403 Forbidden error
	ErrForbidden = errors.New("forbidden")
	// ErrNotFound represents a 404 Not Found error
	ErrNotFound = errors.New("not found")
	// ErrRateLimit represents a 429 Too Many Requests error
	ErrRateLimit = errors.New("rate limit exceeded")
	// ErrInternal represents a 500 Internal Server Error
	ErrInternal = errors.New("internal server error")
	// ErrTimeout represents a client-side timeout error
	ErrTimeout = errors.New("request timed out")
	// ErrInvalidInput represents invalid client-side input
	ErrInvalidInput = errors.New("invalid input")
)

// ErrorResponse represents an error response from the Huntress API
type ErrorResponse struct {
	StatusCode int    `json:"status"`
	Code       string `json:"code,omitempty"`
	Message    string `json:"message"`
	RequestID  string `json:"request_id,omitempty"`
	Details    any    `json:"details,omitempty"`
}

// Error implements the error interface
func (e *ErrorResponse) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("API error %d (%s): %s [request-id: %s]",
			e.StatusCode, e.Code, e.Message, e.RequestID)
	}
	return fmt.Sprintf("API error %d (%s): %s", e.StatusCode, e.Code, e.Message)
}

// Unwrap returns the underlying error for the given status code
func (e *ErrorResponse) Unwrap() error {
	return StatusCodeToError(e.StatusCode, e.Message)
}

// StatusCodeToError converts an HTTP status code to a corresponding error
func StatusCodeToError(statusCode int, message string) error {
	switch statusCode {
	case http.StatusBadRequest:
		return fmt.Errorf("%w: %s", ErrBadRequest, message)
	case http.StatusUnauthorized:
		return fmt.Errorf("%w: %s", ErrUnauthorized, message)
	case http.StatusForbidden:
		return fmt.Errorf("%w: %s", ErrForbidden, message)
	case http.StatusNotFound:
		return fmt.Errorf("%w: %s", ErrNotFound, message)
	case http.StatusTooManyRequests:
		return fmt.Errorf("%w: %s", ErrRateLimit, message)
	case http.StatusInternalServerError:
		return fmt.Errorf("%w: %s", ErrInternal, message)
	default:
		if statusCode >= 400 && statusCode < 500 {
			return fmt.Errorf("client error %d: %s", statusCode, message)
		} else if statusCode >= 500 {
			return fmt.Errorf("server error %d: %s", statusCode, message)
		}
		return fmt.Errorf("unexpected status code %d: %s", statusCode, message)
	}
}

// IsNotFound returns true if the error is a Not Found error
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsUnauthorized returns true if the error is an Unauthorized error
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsForbidden returns true if the error is a Forbidden error
func IsForbidden(err error) bool {
	return errors.Is(err, ErrForbidden)
}

// IsRateLimit returns true if the error is a Rate Limit error
func IsRateLimit(err error) bool {
	return errors.Is(err, ErrRateLimit)
}

// IsBadRequest returns true if the error is a Bad Request error
func IsBadRequest(err error) bool {
	return errors.Is(err, ErrBadRequest)
}

// IsTimeout returns true if the error is a Timeout error
func IsTimeout(err error) bool {
	return errors.Is(err, ErrTimeout)
}
