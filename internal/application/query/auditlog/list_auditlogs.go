// Package auditlog provides query handlers for audit log queries.
package auditlog

import (
	"context"
	"fmt"

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

// NewListAuditLogsHandler creates a new ListAuditLogsHandler.
func NewListAuditLogsHandler(repo repository.AuditLogRepository) *ListAuditLogsHandler {
	return &ListAuditLogsHandler{Repo: repo}
}

// Handle lists audit logs using the provided query parameters.
func (h *ListAuditLogsHandler) Handle(ctx context.Context, query ListAuditLogsQuery) ([]*auditlog.AuditLog, *common.Pagination, error) {
	logs, pagination, err := h.Repo.List(ctx, query.Params)
	if err != nil {
		return nil, nil, fmt.Errorf("list audit logs handler: %w", err)
	}
	return logs, pagination, nil
}
