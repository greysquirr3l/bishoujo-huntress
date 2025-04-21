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

// List lists agents with optional filtering
func (s *agentService) List(ctx context.Context, opts *AgentListOptions) ([]*Agent, *Pagination, error) {
	path := "/agents"
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
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

// Update updates an agent
func (s *agentService) Update(ctx context.Context, id string, input *AgentUpdateInput) (*Agent, error) {
	path := fmt.Sprintf("/agents/%s", id)
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, input)
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

// Delete deletes an agent
func (s *agentService) Delete(ctx context.Context, id string) error {
	path := fmt.Sprintf("/agents/%s", id)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// UpdateStatus updates the status of an agent
func (s *agentService) UpdateStatus(ctx context.Context, id string, status string) error {
	path := fmt.Sprintf("/agents/%s/status", id)
	statusData := map[string]string{"status": status}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, statusData)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// GetStatistics retrieves agent statistics
func (s *agentService) GetStatistics(ctx context.Context, id string) (*AgentStatistics, error) {
	path := fmt.Sprintf("/agents/%s/statistics", id)
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
