// Package huntress provides a client for the Huntress API
package huntress

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// agentService implements the AgentService interface
type agentService struct {
	client *Client
}

// Get retrieves agent details by ID
func (s *agentService) Get(ctx context.Context, id string) (*Agent, error) {
	path := fmt.Sprintf("/agents/%s", id)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for Get: %w", err)
	}

	agent := new(Agent)
	resp, err := s.client.Do(ctx, req, agent)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for Get: %w", err)
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}

	return agent, nil
}

// List returns all agents with optional filtering
func (s *agentService) List(ctx context.Context, params *AgentListOptions) ([]*Agent, *Pagination, error) {
	if params != nil {
		if err := params.Validate(); err != nil {
			return nil, nil, fmt.Errorf("invalid agent list params: %w", err)
		}
	}
	// Advanced filtering: convert params to map[string]interface{} using correct types
	filters := map[string]interface{}{}
	if params != nil {
		if params.Page > 0 {
			filters["page"] = params.Page
		}
		if params.PerPage > 0 {
			filters["per_page"] = params.PerPage
		}
		if params.OrganizationID > 0 {
			filters["organization_id"] = params.OrganizationID
		}
		if params.Status != "" {
			filters["status"] = params.Status
		}
		if params.Platform != "" {
			filters["platform"] = params.Platform
		}
		if params.Hostname != "" {
			filters["hostname"] = params.Hostname
		}
		if params.IPAddress != "" {
			filters["ip_address"] = params.IPAddress
		}
		if params.Version != "" {
			filters["version"] = params.Version
		}
		if params.LastSeenSince != nil {
			filters["last_seen_since"] = params.LastSeenSince.Format(time.RFC3339)
		}
		if len(params.Tags) > 0 {
			filters["tags"] = params.Tags
		}
		if params.Search != "" {
			filters["search"] = params.Search
		}
	}
	var agents []*Agent
	pagination, err := listResource(ctx, s.client, "/agents", filters, &agents)
	if err != nil {
		return nil, nil, err
	}
	return agents, pagination, nil
}

// GetStats retrieves statistics for a specific agent
func (s *agentService) GetStats(ctx context.Context, id string) (*AgentStatistics, error) {
	path := fmt.Sprintf("/agents/%s/stats", id)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetStats: %w", err)
	}

	stats := new(AgentStatistics)
	resp, err := s.client.Do(ctx, req, stats)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for GetStats: %w", err)
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}

	return stats, nil
}

// Update updates an existing agent
func (s *agentService) Update(ctx context.Context, id string, agent map[string]interface{}) (*Agent, error) {
	path := fmt.Sprintf("/agents/%s", id)
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, agent)
	if err != nil {
		return nil, err
	}

	updatedAgent := new(Agent)
	resp, err := s.client.Do(ctx, req, updatedAgent)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}

	return updatedAgent, nil
}

// Delete removes an agent
func (s *agentService) Delete(ctx context.Context, id string) error {
	path := fmt.Sprintf("/agents/%s", id)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return err
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}

	return nil
}
