// Package auditlog provides query handlers for audit log queries.
package auditlog

import (
	"context"
	"fmt"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/auditlog"
	"github.com/greysquirr3l/bishoujo-huntress/internal/ports/repository"
)

// GetAuditLogQuery defines parameters for getting a single audit log entry.
type GetAuditLogQuery struct {
	ID string
}

// GetAuditLogHandler handles retrieving a single audit log entry.
type GetAuditLogHandler struct {
	Repo repository.AuditLogRepository
}

// NewGetAuditLogHandler creates a new GetAuditLogHandler.
func NewGetAuditLogHandler(repo repository.AuditLogRepository) *GetAuditLogHandler {
	return &GetAuditLogHandler{Repo: repo}
}

// Handle retrieves a single audit log entry by ID.
func (h *GetAuditLogHandler) Handle(ctx context.Context, query GetAuditLogQuery) (*auditlog.AuditLog, error) {
	entry, err := h.Repo.Get(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("get audit log handler: %w", err)
	}
	return entry, nil
}
