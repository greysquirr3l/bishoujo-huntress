package huntress

import (
	"context"
	"fmt"

	api "github.com/greysquirr3l/bishoujo-huntress/internal/adapters/api"
	internal_auditlog "github.com/greysquirr3l/bishoujo-huntress/internal/domain/auditlog"
	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/common"
	"github.com/greysquirr3l/bishoujo-huntress/internal/ports/repository"
)

// NewAuditLogService creates a new AuditLogService for testing or custom wiring.
func NewAuditLogService(repo repository.AuditLogRepository) AuditLogService {
	return &auditLogService{repo: repo}
}

// AuditLogService provides access to Huntress audit logs (typed, public).
type auditLogService struct {
	repo repository.AuditLogRepository
}

func newInternalAuditLogRepo(client *Client) repository.AuditLogRepository {
	// This function should return the internal API adapter for audit logs.
	// For now, assume it is in internal/adapters/api and named AuditLogRepository.
	return &internalAuditLogRepoAdapter{client: client}
}

// List returns audit logs matching the given parameters.
func (s *auditLogService) List(ctx context.Context, params *AuditLogListParams) ([]*AuditLog, *Pagination, error) {
	// Convert public params to internal
	internalParams := &internal_auditlog.ListParams{
		StartTime:    params.StartTime,
		EndTime:      params.EndTime,
		Actor:        params.Actor,
		Action:       params.Action,
		ResourceType: params.ResourceType,
		ResourceID:   params.ResourceID,
		Page:         params.Page,
		Limit:        params.Limit,
	}
	logs, pag, err := s.repo.List(ctx, internalParams)
	if err != nil {
		return nil, nil, fmt.Errorf("listing audit logs: %w", err)
	}
	// Convert internal to public
	result := make([]*AuditLog, len(logs))
	for i, l := range logs {
		result[i] = &AuditLog{
			ID:           l.ID,
			Timestamp:    l.Timestamp,
			Actor:        l.Actor,
			Action:       l.Action,
			ResourceType: l.ResourceType,
			ResourceID:   l.ResourceID,
			Description:  l.Description,
			Metadata:     l.Metadata,
		}
	}
	return result, &Pagination{
		CurrentPage: pag.Page,
		PerPage:     pag.PerPage,
		TotalPages:  pag.TotalPages,
		TotalItems:  pag.TotalItems,
	}, nil
}

// Get returns a single audit log entry by ID.
func (s *auditLogService) Get(ctx context.Context, id string) (*AuditLog, error) {
	log, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting audit log by id: %w", err)
	}
	return &AuditLog{
		ID:           log.ID,
		Timestamp:    log.Timestamp,
		Actor:        log.Actor,
		Action:       log.Action,
		ResourceType: log.ResourceType,
		ResourceID:   log.ResourceID,
		Description:  log.Description,
		Metadata:     log.Metadata,
	}, nil
}

// internalAuditLogRepoAdapter adapts the internal API repo to the public interface.
type internalAuditLogRepoAdapter struct {
	client *Client
}

// Implement repository.AuditLogRepository methods by delegating to internal/adapters/api.AuditLogRepository
func (a *internalAuditLogRepoAdapter) Get(ctx context.Context, id string) (*internal_auditlog.AuditLog, error) {
	apiRepo := a.getAPIRepo()
	log, err := apiRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("apiRepo.Get: %w", err)
	}
	return log, nil
}

func (a *internalAuditLogRepoAdapter) List(ctx context.Context, params *internal_auditlog.ListParams) ([]*internal_auditlog.AuditLog, *common.Pagination, error) {
	apiRepo := a.getAPIRepo()
	logs, pag, err := apiRepo.List(ctx, params)
	if err != nil {
		return nil, nil, fmt.Errorf("apiRepo.List: %w", err)
	}
	return logs, pag, nil
}

// getAPIRepo returns an instance of the internal API adapter for audit logs.
func (a *internalAuditLogRepoAdapter) getAPIRepo() *api.AuditLogRepository {
	return &api.AuditLogRepository{
		Client:    a.client.httpClient,
		BaseURL:   a.client.baseURL,
		APIKey:    a.client.apiKey,
		APISecret: a.client.apiSecret,
	}
}
