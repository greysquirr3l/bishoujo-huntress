// Package huntress provides a client for the Huntress API
package huntress

import (
	"errors"
	"fmt"
	"net/http"
)

// Error is the base error interface for all errors returned by this package
type Error interface {
	error
	// Code returns the error code
	Code() string
	// StatusCode returns the HTTP status code associated with this error, if any
	StatusCode() int
}

// InternalAPIError defines the structure for API errors received from the Huntress API
type InternalAPIError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
	StatusCode int    `json:"-"` // Not part of the JSON response, set from HTTP status code
}

// Granular API error types for specific Huntress API error codes
var (
	ErrWebhookValidationFailed = &APIError{internal: &InternalAPIError{StatusCode: 422, Code: "WEBHOOK_VALIDATION_FAILED", Message: "Webhook payload validation failed"}}
	ErrWebhookParseFailed      = &APIError{internal: &InternalAPIError{StatusCode: 400, Code: "WEBHOOK_PARSE_FAILED", Message: "Webhook payload could not be parsed"}}
	ErrAgentNotFound           = &APIError{internal: &InternalAPIError{StatusCode: 404, Code: "AGENT_NOT_FOUND", Message: "Agent not found"}}
	ErrOrgNotFound             = &APIError{internal: &InternalAPIError{StatusCode: 404, Code: "ORG_NOT_FOUND", Message: "Organization not found"}}
	ErrWebhookNotFound         = &APIError{internal: &InternalAPIError{StatusCode: 404, Code: "WEBHOOK_NOT_FOUND", Message: "Webhook not found"}}
	ErrIntegrationNotFound     = &APIError{internal: &InternalAPIError{StatusCode: 404, Code: "INTEGRATION_NOT_FOUND", Message: "Integration not found"}}
	ErrInvalidEventType        = &APIError{internal: &InternalAPIError{StatusCode: 400, Code: "INVALID_EVENT_TYPE", Message: "Invalid event type for webhook"}}
	// Add more as needed for other API error codes
)

func (e *InternalAPIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// InternalRateLimitError defines the structure for rate limit errors
type InternalRateLimitError struct {
	RetryAfter int `json:"retry_after"`
}

func (e *InternalRateLimitError) Error() string {
	return fmt.Sprintf("rate limit exceeded, retry after %d seconds", e.RetryAfter)
}

// InternalRequestError defines the structure for request errors
type InternalRequestError struct {
	Err error
}

func (e *InternalRequestError) Error() string {
	if e.Err == nil {
		return "request error"
	}
	return fmt.Sprintf("request error: %s", e.Err.Error())
}

// APIError represents an error returned by the Huntress API
type APIError struct {
	internal *InternalAPIError
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.internal == nil {
		return "API error"
	}
	return e.internal.Error()
}

// Code returns the API error code
func (e *APIError) Code() string {
	if e.internal == nil {
		return ""
	}
	return e.internal.Code
}

// StatusCode returns the HTTP status code
func (e *APIError) StatusCode() int {
	if e.internal == nil {
		return 0
	}
	return e.internal.StatusCode
}

// Message returns the error message
func (e *APIError) Message() string {
	if e.internal == nil {
		return ""
	}
	return e.internal.Message
}

// Details returns additional error details
func (e *APIError) Details() string {
	if e.internal == nil {
		return ""
	}
	return e.internal.Details
}

// IsNotFound returns true if the error is a 404 Not Found
func (e *APIError) IsNotFound() bool {
	return e.internal != nil && e.internal.StatusCode == http.StatusNotFound
}

// IsUnauthorized returns true if the error is a 401 Unauthorized
func (e *APIError) IsUnauthorized() bool {
	return e.internal != nil && e.internal.StatusCode == http.StatusUnauthorized
}

// IsForbidden returns true if the error is a 403 Forbidden
func (e *APIError) IsForbidden() bool {
	return e.internal != nil && e.internal.StatusCode == http.StatusForbidden
}

// RateLimitError indicates that the API rate limit has been exceeded
type RateLimitError struct {
	internal *InternalRateLimitError
}

// Error implements the error interface
func (e *RateLimitError) Error() string {
	if e.internal == nil {
		return "rate limit exceeded"
	}
	return e.internal.Error()
}

// Code returns the error code
func (e *RateLimitError) Code() string {
	return "RATE_LIMIT_EXCEEDED"
}

// StatusCode returns the HTTP status code
func (e *RateLimitError) StatusCode() int {
	return http.StatusTooManyRequests
}

// RetryAfter returns the number of seconds to wait before retrying
func (e *RateLimitError) RetryAfter() int {
	if e.internal == nil {
		return 60 // Default to 60 seconds
	}
	return e.internal.RetryAfter
}

// RequestError represents an error that occurred while making a request
type RequestError struct {
	internal *InternalRequestError
}

// Error implements the error interface
func (e *RequestError) Error() string {
	if e.internal == nil {
		return "request error"
	}
	return e.internal.Error()
}

// Code returns the error code
func (e *RequestError) Code() string {
	return "REQUEST_ERROR"
}

// StatusCode returns the HTTP status code
func (e *RequestError) StatusCode() int {
	return 0 // No HTTP status for request errors
}

// Unwrap returns the underlying error
func (e *RequestError) Unwrap() error {
	if e.internal == nil || e.internal.Err == nil {
		return nil
	}
	return e.internal.Err
}

// IsAPIError returns true if the error is an APIError
func IsAPIError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr)
}

// IsRateLimitError returns true if the error is a RateLimitError
func IsRateLimitError(err error) bool {
	var rateLimitErr *RateLimitError
	return errors.As(err, &rateLimitErr)
}

// IsRequestError returns true if the error is a RequestError
func IsRequestError(err error) bool {
	var requestErr *RequestError
	return errors.As(err, &requestErr)
}

// IsNotFoundError returns true if the error is a 404 Not Found error
func IsNotFoundError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.IsNotFound()
}

// IsAuthError returns true if the error is an authentication error (401 or 403)
func IsAuthError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && (apiErr.IsUnauthorized() || apiErr.IsForbidden())
}

// AsAPIError attempts to convert an error to an APIError
func AsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr, true
	}

	var internalAPIErr *InternalAPIError
	if errors.As(err, &internalAPIErr) {
		return &APIError{internal: internalAPIErr}, true
	}

	return nil, false
}

// AsRateLimitError attempts to convert an error to a RateLimitError
func AsRateLimitError(err error) (*RateLimitError, bool) {
	var rateLimitErr *RateLimitError
	if errors.As(err, &rateLimitErr) {
		return rateLimitErr, true
	}

	var internalRateLimitErr *InternalRateLimitError
	if errors.As(err, &internalRateLimitErr) {
		return &RateLimitError{internal: internalRateLimitErr}, true
	}

	return nil, false
}

// NewError creates a new general error with the given message
func NewError(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}
