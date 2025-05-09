// Package repository defines repository interfaces for audit logs.
package repository

import (
	"context"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/auditlog"
	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/common"
)

// AuditLogRepository defines repository operations for audit logs.
type AuditLogRepository interface {
	// Get retrieves a single audit log entry by ID.
	Get(ctx context.Context, id string) (*auditlog.AuditLog, error)
	// List retrieves audit log entries with optional filters and pagination.
	List(ctx context.Context, params *auditlog.ListParams) ([]*auditlog.AuditLog, *common.Pagination, error)
}
