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

// List lists incidents with optional filtering
func (s *incidentService) List(ctx context.Context, opts *IncidentListOptions) ([]*Incident, *Pagination, error) {
	path := "/incidents"
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
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

// Update updates an incident
func (s *incidentService) Update(ctx context.Context, id string, input *IncidentUpdateInput) (*Incident, error) {
	path := fmt.Sprintf("/incidents/%s", id)
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, input)
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

// UpdateStatus updates the status of an incident
func (s *incidentService) UpdateStatus(ctx context.Context, id string, status string) error {
	path := fmt.Sprintf("/incidents/%s/status", id)
	statusData := map[string]string{"status": status}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, statusData)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// AddNote adds a note to an incident
func (s *incidentService) AddNote(ctx context.Context, id string, note string) error {
	path := fmt.Sprintf("/incidents/%s/notes", id)
	noteData := map[string]string{"content": note}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, noteData)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// ListNotes lists notes for an incident
func (s *incidentService) ListNotes(ctx context.Context, id string, opts *ListOptions) ([]*IncidentNote, *Pagination, error) {
	path := fmt.Sprintf("/incidents/%s/notes", id)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var notes []*IncidentNote
	resp, err := s.client.Do(ctx, req, &notes)
	if err != nil {
		return nil, nil, err
	}

	pagination := extractPagination(resp)
	return notes, pagination, nil
}
