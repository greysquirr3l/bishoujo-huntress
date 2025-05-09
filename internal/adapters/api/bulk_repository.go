// Package api provides the Bulk Operations API adapter for Huntress.
package api

import (
	"context"
	"net/http"
)

// BulkRepository provides access to Huntress bulk operations.
type BulkRepository struct {
	Client    *http.Client
	BaseURL   string
	APIKey    string
	APISecret string
}

// BulkAgentAction performs a bulk action on agents (e.g., update, delete).
func (r *BulkRepository) BulkAgentAction(ctx context.Context, action string, agentIDs []string, payload interface{}) (map[string]interface{}, error) {
	endpoint := "/bulk/agents/" + action
	return doPostBulkActionAndDecode(ctx, r.Client, r.BaseURL, endpoint, r.APIKey, r.APISecret, "agent_ids", agentIDs, payload)
}

// BulkOrgAction performs a bulk action on organizations (e.g., update, delete).
// (Removed duplicate BulkOrgAction implementation)
func (r *BulkRepository) BulkOrgAction(ctx context.Context, action string, orgIDs []string, payload interface{}) (map[string]interface{}, error) {
	endpoint := "/bulk/organizations/" + action
	return doPostBulkActionAndDecode(ctx, r.Client, r.BaseURL, endpoint, r.APIKey, r.APISecret, "organization_ids", orgIDs, payload)
}
