// Package huntress provides a client for the Huntress API
package huntress

import (
	"context"
	"fmt"
	"net/http"
)

// incidentService implements the IncidentService interface
type incidentService struct {
	client *Client
}

// Get retrieves incident details by ID
func (s *incidentService) Get(ctx context.Context, id string) (*Incident, error) {
	path := fmt.Sprintf("/incidents/%s", id)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for Get: %w", err)
	}

	incident := new(Incident)
	resp, err := s.client.Do(ctx, req, incident)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for Get: %w", err)
	}
	if resp != nil {
		errClose := resp.Body.Close()
		if errClose != nil {
			return nil, fmt.Errorf("incident get: error closing response body: %w", errClose)
		}
	}

	return incident, nil
}

// List returns all incidents with optional filtering
func (s *incidentService) List(ctx context.Context, params *IncidentListOptions) ([]*Incident, *Pagination, error) {
	if params != nil {
		if err := params.Validate(); err != nil {
			return nil, nil, fmt.Errorf("invalid incident list params: %w", err)
		}
	}
	var incidents []*Incident
	pagination, err := listResource(ctx, s.client, "/incidents", params, &incidents)
	if err != nil {
		return nil, nil, err
	}
	return incidents, pagination, nil
}

// UpdateStatus updates the status of an incident
func (s *incidentService) UpdateStatus(ctx context.Context, id string, status string) (*Incident, error) {
	path := fmt.Sprintf("/incidents/%s/status", id)
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, map[string]string{"status": status})
	if err != nil {
		return nil, fmt.Errorf("failed to create request for UpdateStatus: %w", err)
	}

	incident := new(Incident)
	resp, err := s.client.Do(ctx, req, incident)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for UpdateStatus: %w", err)
	}
	if resp != nil {
		errClose := resp.Body.Close()
		if errClose != nil {
			return nil, fmt.Errorf("incident update status: error closing response body: %w", errClose)
		}
	}

	return incident, nil
}

// Assign assigns an incident to a user
func (s *incidentService) Assign(ctx context.Context, id string, userID string) (*Incident, error) {
	path := fmt.Sprintf("/incidents/%s/assign", id)
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, map[string]string{
		"user_id": userID,
	})
	if err != nil {
		return nil, err
	}

	incident := new(Incident)
	resp, err := s.client.Do(ctx, req, incident)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}

	return incident, nil
}
