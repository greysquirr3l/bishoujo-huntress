// Package huntress provides a client for the Huntress API
package huntress

import (
	"context"
	"fmt"
	"net/http"
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
		return nil, err
	}

	agent := new(Agent)
	_, err = s.client.Do(ctx, req, agent)
	if err != nil {
		return nil, err
	}

	return agent, nil
}

// List returns all agents with optional filtering
func (s *agentService) List(ctx context.Context, params *AgentListOptions) ([]*Agent, *Pagination, error) {
	path := "/agents"
	if params != nil {
		query, err := addQueryParams(path, params)
		if err != nil {
			return nil, nil, err
		}
		path = query
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var agents []*Agent
	resp, err := s.client.Do(ctx, req, &agents)
	if err != nil {
		return nil, nil, err
	}

	pagination := extractPagination(resp)
	return agents, pagination, nil
}

// GetStats retrieves statistics for a specific agent
func (s *agentService) GetStats(ctx context.Context, id string) (*AgentStatistics, error) {
	path := fmt.Sprintf("/agents/%s/stats", id)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	stats := new(AgentStatistics)
	_, err = s.client.Do(ctx, req, stats)
	if err != nil {
		return nil, err
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
	_, err = s.client.Do(ctx, req, updatedAgent)
	if err != nil {
		return nil, err
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

	_, err = s.client.Do(ctx, req, nil)
	return err
}
