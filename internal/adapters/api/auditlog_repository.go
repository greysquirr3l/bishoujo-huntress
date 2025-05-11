// Package api provides the Audit Log API adapter for Huntress.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/auditlog"
	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/common"
)

// (Removed duplicate Get method)

// Search allows searching audit logs with advanced filters.
func (r *AuditLogRepository) Search(ctx context.Context, filters map[string]string) ([]map[string]interface{}, error) {
	return doGetWithQueryAndDecode(ctx, r.Client, r.BaseURL, "/audit-logs/search", r.APIKey, r.APISecret, filters)
}

// AuditLogRepository provides access to Huntress audit logs.
// Implements repository.AuditLogRepository.
type AuditLogRepository struct {
	Client    *http.Client
	BaseURL   string
	APIKey    string
	APISecret string
}

// Get retrieves a specific audit log entry by ID.
func (r *AuditLogRepository) Get(ctx context.Context, id string) (*auditlog.AuditLog, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", r.BaseURL+"/audit-logs/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("auditlog get: %w", err)
	}
	req.SetBasicAuth(r.APIKey, r.APISecret)
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("auditlog get: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing response body: %v\n", err)
		}
	}()
	if resp.StatusCode == http.StatusNotFound {
		return nil, common.NewDomainError(common.ErrCodeNotFound, "audit log not found", nil)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("auditlog get: unexpected status: %d", resp.StatusCode)
	}
	var out auditlog.AuditLog
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("auditlog get: decode: %w", err)
	}
	return &out, nil
}

// List returns all audit log entries (optionally with filters).
func (r *AuditLogRepository) List(ctx context.Context, params *auditlog.ListParams) ([]*auditlog.AuditLog, *common.Pagination, error) {
	// Build query params
	q := make(map[string]string)
	if params != nil {
		if params.StartTime != nil {
			q["start_time"] = params.StartTime.Format(time.RFC3339)
		}
		if params.EndTime != nil {
			q["end_time"] = params.EndTime.Format(time.RFC3339)
		}
		if params.Actor != nil {
			q["actor"] = *params.Actor
		}
		if params.Action != nil {
			q["action"] = *params.Action
		}
		if params.ResourceType != nil {
			q["resource_type"] = *params.ResourceType
		}
		if params.ResourceID != nil {
			q["resource_id"] = *params.ResourceID
		}
		if params.Page > 0 {
			q["page"] = fmt.Sprintf("%d", params.Page)
		}
		if params.Limit > 0 {
			q["per_page"] = fmt.Sprintf("%d", params.Limit)
		}
	}
	req, err := http.NewRequestWithContext(ctx, "GET", r.BaseURL+"/audit-logs", nil)
	if err != nil {
		return nil, nil, fmt.Errorf("auditlog list: %w", err)
	}
	query := req.URL.Query()
	for k, v := range q {
		query.Set(k, v)
	}
	req.URL.RawQuery = query.Encode()
	req.SetBasicAuth(r.APIKey, r.APISecret)
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("auditlog list: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Fprintf(os.Stderr, "error closing response body: %v\n", cerr)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("auditlog list: unexpected status: %d", resp.StatusCode)
	}
	var out struct {
		Data []*auditlog.AuditLog `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, nil, fmt.Errorf("auditlog list: decode: %w", err)
	}
	pagination := &common.Pagination{
		Page:       params.Page,
		PerPage:    params.Limit,
		TotalItems: len(out.Data), // TODO: parse from headers if available
		TotalPages: 1,             // TODO: parse from headers if available
	}
	return out.Data, pagination, nil
}
