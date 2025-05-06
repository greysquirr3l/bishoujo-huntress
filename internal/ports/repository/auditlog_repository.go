// Package repository defines repository interfaces for audit logs.
package repository

import (
	"context"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/auditlog"
	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/common"
)

type AuditLogRepository interface {
	Get(ctx context.Context, id string) (*auditlog.AuditLog, error)
	List(ctx context.Context, params *auditlog.ListParams) ([]*auditlog.AuditLog, *common.Pagination, error)
}
