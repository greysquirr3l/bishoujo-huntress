// Package api provides the Bulk Operations API adapter for Huntress.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	endpoint := r.BaseURL + "/bulk/agents/" + action
	reqBody := map[string]interface{}{
		"agent_ids": agentIDs,
		"data":      payload,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(r.APIKey, r.APISecret)
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("bulk agent action failed: %d", resp.StatusCode)
	}
	var out map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

// BulkOrgAction performs a bulk action on organizations (e.g., update, delete).
// (Removed duplicate BulkOrgAction implementation)
func (r *BulkRepository) BulkOrgAction(ctx context.Context, action string, orgIDs []string, payload interface{}) (map[string]interface{}, error) {
	endpoint := r.BaseURL + "/bulk/organizations/" + action
	reqBody := map[string]interface{}{
		"organization_ids": orgIDs,
		"data":             payload,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(r.APIKey, r.APISecret)
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("bulk org action failed: %d", resp.StatusCode)
	}
	var out map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}
