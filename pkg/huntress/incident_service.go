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
		return nil, err
	}

	incident := new(Incident)
	_, err = s.client.Do(ctx, req, incident)
	if err != nil {
		return nil, err
	}

	return incident, nil
}

// List returns all incidents with optional filtering
func (s *incidentService) List(ctx context.Context, params *IncidentListOptions) ([]*Incident, *Pagination, error) {
	path := "/incidents"
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

	var incidents []*Incident
	resp, err := s.client.Do(ctx, req, &incidents)
	if err != nil {
		return nil, nil, err
	}

	pagination := extractPagination(resp)
	return incidents, pagination, nil
}

// UpdateStatus updates the status of an incident
func (s *incidentService) UpdateStatus(ctx context.Context, id string, status string) (*Incident, error) {
	path := fmt.Sprintf("/incidents/%s/status", id)
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, map[string]string{
		"status": status,
	})
	if err != nil {
		return nil, err
	}

	incident := new(Incident)
	_, err = s.client.Do(ctx, req, incident)
	if err != nil {
		return nil, err
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
	_, err = s.client.Do(ctx, req, incident)
	if err != nil {
		return nil, err
	}

	return incident, nil
}
