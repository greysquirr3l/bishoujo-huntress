// Package errors defines API error types for the Huntress domain.
package errors

import (
	"fmt"
	"net/http"
)

// More granular API error types for specific Huntress API error codes
var (
	ErrWebhookValidationFailed = NewAPIError(422, "WEBHOOK_VALIDATION_FAILED", "Webhook payload validation failed", "")
	ErrWebhookParseFailed      = NewAPIError(400, "WEBHOOK_PARSE_FAILED", "Webhook payload could not be parsed", "")
	ErrAgentNotFound           = NewAPIError(404, "AGENT_NOT_FOUND", "Agent not found", "")
	ErrOrgNotFound             = NewAPIError(404, "ORG_NOT_FOUND", "Organization not found", "")
	ErrWebhookNotFound         = NewAPIError(404, "WEBHOOK_NOT_FOUND", "Webhook not found", "")
	ErrIntegrationNotFound     = NewAPIError(404, "INTEGRATION_NOT_FOUND", "Integration not found", "")
	ErrInvalidEventType        = NewAPIError(400, "INVALID_EVENT_TYPE", "Invalid event type for webhook", "")
	// Add more as needed for other API error codes
)

// APIError represents an error returned by the API
type APIError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code,omitempty"`
	Message    string `json:"message"`
	Details    string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%d] %s: %s - %s", e.StatusCode, e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%d] %s: %s", e.StatusCode, e.Code, e.Message)
}

// NewAPIError creates a new API error
func NewAPIError(statusCode int, code string, message string, details string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Details:    details,
	}
}

// Common API error codes
const (
	ErrCodeBadRequest          = "BAD_REQUEST"
	ErrCodeUnauthorized        = "UNAUTHORIZED"
	ErrCodeForbidden           = "FORBIDDEN"
	ErrCodeNotFound            = "NOT_FOUND"
	ErrCodeMethodNotAllowed    = "METHOD_NOT_ALLOWED"
	ErrCodeConflict            = "CONFLICT"
	ErrCodeTooManyRequests     = "TOO_MANY_REQUESTS"
	ErrCodeInternalServerError = "INTERNAL_SERVER_ERROR"
	ErrCodeServiceUnavailable  = "SERVICE_UNAVAILABLE"
)

// Common API errors
var (
	ErrBadRequest          = NewAPIError(http.StatusBadRequest, ErrCodeBadRequest, "Bad request", "")
	ErrUnauthorized        = NewAPIError(http.StatusUnauthorized, ErrCodeUnauthorized, "Unauthorized", "")
	ErrForbidden           = NewAPIError(http.StatusForbidden, ErrCodeForbidden, "Forbidden", "")
	ErrNotFound            = NewAPIError(http.StatusNotFound, ErrCodeNotFound, "Not found", "")
	ErrMethodNotAllowed    = NewAPIError(http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed, "Method not allowed", "")
	ErrConflict            = NewAPIError(http.StatusConflict, ErrCodeConflict, "Conflict", "")
	ErrTooManyRequests     = NewAPIError(http.StatusTooManyRequests, ErrCodeTooManyRequests, "Rate limit exceeded", "")
	ErrInternalServerError = NewAPIError(http.StatusInternalServerError, ErrCodeInternalServerError, "Internal server error", "")
	ErrServiceUnavailable  = NewAPIError(http.StatusServiceUnavailable, ErrCodeServiceUnavailable, "Service unavailable", "")
)
