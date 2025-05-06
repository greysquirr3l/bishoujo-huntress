// Package auditlog contains the domain model for Huntress audit logs.
package auditlog

import (
	"time"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/common"
)

// AuditLog represents a single audit log entry.
type AuditLog struct {
	ID           string         `json:"id"`
	Timestamp    time.Time      `json:"timestamp"`
	Actor        string         `json:"actor"`
	Action       string         `json:"action"`
	ResourceType string         `json:"resourceType"`
	ResourceID   string         `json:"resourceId"`
	Description  string         `json:"description"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

// ListParams defines optional parameters for listing audit logs.
type ListParams struct {
	StartTime    *time.Time
	EndTime      *time.Time
	Actor        *string
	Action       *string
	ResourceType *string
	ResourceID   *string
	Page         int
	Limit        int
}

// Validate checks if the audit log entry is valid.
func (a *AuditLog) Validate() error {
	if a.ID == "" {
		return common.ErrInvalidID
	}
	if a.Timestamp.IsZero() {
		return common.ErrInvalidTimestamp
	}
	if a.Actor == "" {
		return common.ErrEmptyActor
	}
	if a.Action == "" {
		return common.ErrEmptyAction
	}
	if a.ResourceType == "" {
		return common.ErrEmptyResourceType
	}
	if a.ResourceID == "" {
		return common.ErrEmptyResourceID
	}
	return nil
}
