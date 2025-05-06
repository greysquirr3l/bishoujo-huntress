// Package common provides shared domain error types and helpers.
package common

import (
	"fmt"
)

// DomainError represents a domain-specific error
type DomainError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the wrapped error
func (e *DomainError) Unwrap() error {
	return e.Err
}

// NewDomainError creates a new domain error
func NewDomainError(code string, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

// Error implements the error interface
func (e ValidationErrors) Error() string {
	if len(e) == 1 {
		return fmt.Sprintf("validation error: %s %s", e[0].Field, e[0].Message)
	}
	return fmt.Sprintf("validation errors: %d errors", len(e))
}

// AddError adds a validation error
func (e *ValidationErrors) AddError(field, message string) {
	*e = append(*e, ValidationError{Field: field, Message: message})
}

// HasErrors returns true if there are validation errors
func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

// Common domain error codes
const (
	ErrCodeValidation   = "VALIDATION_ERROR"
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeConflict     = "CONFLICT"
)

var (
	ErrInvalidID         = NewDomainError(ErrCodeValidation, "invalid ID", nil)
	ErrInvalidTimestamp  = NewDomainError(ErrCodeValidation, "invalid timestamp", nil)
	ErrEmptyActor        = NewDomainError(ErrCodeValidation, "actor is required", nil)
	ErrEmptyAction       = NewDomainError(ErrCodeValidation, "action is required", nil)
	ErrEmptyResourceType = NewDomainError(ErrCodeValidation, "resource type is required", nil)
	ErrEmptyResourceID   = NewDomainError(ErrCodeValidation, "resource ID is required", nil)
)
