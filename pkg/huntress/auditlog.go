package huntress

import (
	"time"
)

// AuditLog represents a Huntress audit log entry (public model).
type AuditLog struct {
	ID           string                 `json:"id"`
	Timestamp    time.Time              `json:"timestamp"`
	Actor        string                 `json:"actor"`
	Action       string                 `json:"action"`
	ResourceType string                 `json:"resourceType"`
	ResourceID   string                 `json:"resourceId"`
	Description  string                 `json:"description"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// AuditLogListParams defines parameters for listing audit logs.
type AuditLogListParams struct {
	StartTime    *time.Time
	EndTime      *time.Time
	Actor        *string
	Action       *string
	ResourceType *string
	ResourceID   *string
	Page         int
	Limit        int
}
