// Package auditlog provides query handlers for audit log queries.
package auditlog

import (
	"context"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/auditlog"
	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/common"
	"github.com/greysquirr3l/bishoujo-huntress/internal/ports/repository"
)

// ListAuditLogsQuery defines parameters for listing audit logs.
type ListAuditLogsQuery struct {
	Params *auditlog.ListParams
}

// ListAuditLogsHandler handles listing audit logs.
type ListAuditLogsHandler struct {
	Repo repository.AuditLogRepository
}

func NewListAuditLogsHandler(repo repository.AuditLogRepository) *ListAuditLogsHandler {
	return &ListAuditLogsHandler{Repo: repo}
}

func (h *ListAuditLogsHandler) Handle(ctx context.Context, query ListAuditLogsQuery) ([]*auditlog.AuditLog, *common.Pagination, error) {
	return h.Repo.List(ctx, query.Params)
}
