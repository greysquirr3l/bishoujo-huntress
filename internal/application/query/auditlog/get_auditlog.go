// Package auditlog provides query handlers for audit log queries.
package auditlog

import (
	"context"

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

func NewGetAuditLogHandler(repo repository.AuditLogRepository) *GetAuditLogHandler {
	return &GetAuditLogHandler{Repo: repo}
}

func (h *GetAuditLogHandler) Handle(ctx context.Context, query GetAuditLogQuery) (*auditlog.AuditLog, error) {
	return h.Repo.Get(ctx, query.ID)
}
