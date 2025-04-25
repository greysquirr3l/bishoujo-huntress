// Package repository provides repository implementations for various domain entities.
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/incident"
	"github.com/greysquirr3l/bishoujo-huntress/internal/ports/repository"
)

// IncidentRepositoryImpl implements the incident repository interface
type IncidentRepositoryImpl struct {
	httpClient  HTTPClient
	baseURL     string
	authHeaders map[string]string
}

// NewIncidentRepository creates a new incident repository
func NewIncidentRepository(httpClient HTTPClient, baseURL string, authHeaders map[string]string) repository.IncidentRepository {
	return &IncidentRepositoryImpl{
		httpClient:  httpClient,
		baseURL:     baseURL,
		authHeaders: authHeaders,
	}
}

// Get retrieves an incident by ID
func (r *IncidentRepositoryImpl) Get(ctx context.Context, id string) (*incident.Incident, error) {
	path := fmt.Sprintf("/incidents/%s", id)
	req, err := createRequest(ctx, http.MethodGet, r.baseURL+path, nil, r.authHeaders)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(resp)
	}

	// Parse the response into our DTO structure
	var dto incidentDTO
	if err := json.NewDecoder(resp.Body).Decode(&dto); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert the DTO to a domain model
	incident, err := ToDomainIncident(&dto)
	if err != nil {
		return nil, fmt.Errorf("failed to convert DTO to domain model: %w", err)
	}

	return incident, nil
}

// List retrieves multiple incidents based on filters
func (r *IncidentRepositoryImpl) List(ctx context.Context, filters repository.IncidentFilters) ([]*incident.Incident, repository.Pagination, error) {
	// Construct query parameters
	query := url.Values{}
	if filters.Page > 0 {
		query.Set("page", strconv.Itoa(filters.Page))
	}
	if filters.Limit > 0 {
		query.Set("limit", strconv.Itoa(filters.Limit))
	}
	if filters.OrganizationID != nil {
		query.Set("organization_id", strconv.Itoa(*filters.OrganizationID))
	}
	if filters.AgentID != "" {
		query.Set("agent_id", filters.AgentID)
	}
	if filters.Status != nil {
		query.Set("status", string(*filters.Status))
	}
	if filters.Severity != nil {
		query.Set("severity", string(*filters.Severity))
	}
	if filters.Type != nil {
		query.Set("type", string(*filters.Type))
	}
	if filters.Search != "" {
		query.Set("search", filters.Search)
	}
	if filters.AssignedTo != "" {
		query.Set("assigned_to", filters.AssignedTo)
	}
	if filters.DetectedAfter != nil {
		query.Set("detected_after", filters.DetectedAfter.Format(time.RFC3339))
	}
	if filters.DetectedBefore != nil {
		query.Set("detected_before", filters.DetectedBefore.Format(time.RFC3339))
	}
	// Add tags if present
	for _, tag := range filters.Tags {
		query.Add("tags", tag)
	}
	// Add ordering
	for _, order := range filters.OrderBy {
		direction := "asc"
		if order.Direction == repository.OrderDesc {
			direction = "desc"
		}
		query.Add("sort", fmt.Sprintf("%s:%s", order.Field, direction))
	}

	path := "/incidents"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	req, err := createRequest(ctx, http.MethodGet, r.baseURL+path, nil, r.authHeaders)
	if err != nil {
		return nil, repository.Pagination{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, repository.Pagination{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, repository.Pagination{}, handleErrorResponse(resp)
	}

	var incidentDTOs []incidentDTO
	if err := json.NewDecoder(resp.Body).Decode(&incidentDTOs); err != nil {
		return nil, repository.Pagination{}, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract pagination information from headers
	pagination := extractPagination(resp.Header)

	// Convert DTOs to domain entities
	incidents := make([]*incident.Incident, len(incidentDTOs))
	for i, dto := range incidentDTOs {
		incident, err := ToDomainIncident(&dto)
		if err != nil {
			return nil, repository.Pagination{}, fmt.Errorf("failed to convert incident DTO to domain model at index %d: %w", i, err)
		}
		incidents[i] = incident
	}

	return incidents, pagination, nil
}

// Create creates a new incident
func (r *IncidentRepositoryImpl) Create(ctx context.Context, inc *incident.Incident) error {
	// Convert domain entity to DTO
	incidentDTO := ToIncidentDTO(inc)

	// Create request body
	body, err := json.Marshal(incidentDTO)
	if err != nil {
		return fmt.Errorf("failed to marshal incident: %w", err)
	}

	req, err := createRequest(ctx, http.MethodPost, r.baseURL+"/incidents", body, r.authHeaders)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return handleErrorResponse(resp)
	}

	// Parse the response to get the created incident's ID
	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Parse the ID from the response and update the original incident object
	if idStr, ok := responseData["id"].(string); ok {
		createdID, err := uuid.Parse(idStr)
		if err != nil {
			return fmt.Errorf("failed to parse created incident ID: %w", err)
		}
		inc.ID = createdID
	} else {
		return fmt.Errorf("failed to extract ID from response")
	}

	return nil
}

// Update updates an existing incident
func (r *IncidentRepositoryImpl) Update(ctx context.Context, inc *incident.Incident) error {
	// Convert domain entity to DTO
	incidentDTO := ToIncidentDTO(inc)

	// Create request body
	body, err := json.Marshal(incidentDTO)
	if err != nil {
		return fmt.Errorf("failed to marshal incident: %w", err)
	}

	path := fmt.Sprintf("/incidents/%s", inc.ID.String())
	req, err := createRequest(ctx, http.MethodPut, r.baseURL+path, body, r.authHeaders)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return handleErrorResponse(resp)
	}

	return nil
}

// Delete deletes an incident by ID
func (r *IncidentRepositoryImpl) Delete(ctx context.Context, id string) error {
	path := fmt.Sprintf("/incidents/%s", id)
	req, err := createRequest(ctx, http.MethodDelete, r.baseURL+path, nil, r.authHeaders)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return handleErrorResponse(resp)
	}

	return nil
}

// UpdateStatus updates the status of an incident
func (r *IncidentRepositoryImpl) UpdateStatus(ctx context.Context, id string, status incident.IncidentStatus) error {
	path := fmt.Sprintf("/incidents/%s/status", id)

	// Create request body
	statusData := map[string]string{"status": string(status)}
	body, err := json.Marshal(statusData)
	if err != nil {
		return fmt.Errorf("failed to marshal status: %w", err)
	}

	req, err := createRequest(ctx, http.MethodPut, r.baseURL+path, body, r.authHeaders)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return handleErrorResponse(resp)
	}

	return nil
}

// UpdateAssignee updates the assignee of an incident
func (r *IncidentRepositoryImpl) UpdateAssignee(ctx context.Context, id string, assignee string) error {
	path := fmt.Sprintf("/incidents/%s/assign", id)

	// Create request body
	assignData := map[string]string{"assigned_to": assignee}
	body, err := json.Marshal(assignData)
	if err != nil {
		return fmt.Errorf("failed to marshal assignee: %w", err)
	}

	req, err := createRequest(ctx, http.MethodPut, r.baseURL+path, body, r.authHeaders)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return handleErrorResponse(resp)
	}

	return nil
}

// ListByOrganization retrieves all incidents for an organization
func (r *IncidentRepositoryImpl) ListByOrganization(ctx context.Context, organizationID int, filters repository.IncidentFilters) ([]*incident.Incident, repository.Pagination, error) {
	// Clone filters and set organization ID
	filters.OrganizationID = &organizationID
	return r.List(ctx, filters)
}

// ListByAgent retrieves all incidents for a specific agent
func (r *IncidentRepositoryImpl) ListByAgent(ctx context.Context, agentID string, filters repository.IncidentFilters) ([]*incident.Incident, repository.Pagination, error) {
	// Clone filters and set agent ID
	filters.AgentID = agentID
	return r.List(ctx, filters)
}

// AddNote adds a note to an incident
func (r *IncidentRepositoryImpl) AddNote(ctx context.Context, incidentID string, note incident.Note) error {
	path := fmt.Sprintf("/incidents/%s/notes", incidentID)

	// Create request body
	noteData := map[string]string{"content": note.Content}
	body, err := json.Marshal(noteData)
	if err != nil {
		return fmt.Errorf("failed to marshal note: %w", err)
	}

	req, err := createRequest(ctx, http.MethodPost, r.baseURL+path, body, r.authHeaders)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return handleErrorResponse(resp)
	}

	return nil
}

// AddIOC adds an indicator of compromise to an incident
func (r *IncidentRepositoryImpl) AddIOC(ctx context.Context, incidentID string, ioc incident.IndicatorOfCompromise) error {
	path := fmt.Sprintf("/incidents/%s/iocs", incidentID)

	// Create simplified DTO for API request (only sending required fields, not full conversion)
	iocDTO := map[string]interface{}{
		"type":        ioc.Type,
		"value":       ioc.Value,
		"description": ioc.Description,
		"source":      ioc.Source,
	}

	body, err := json.Marshal(iocDTO)
	if err != nil {
		return fmt.Errorf("failed to marshal IOC: %w", err)
	}

	req, err := createRequest(ctx, http.MethodPost, r.baseURL+path, body, r.authHeaders)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return handleErrorResponse(resp)
	}

	return nil
}

// AddArtifact adds an artifact to an incident
func (r *IncidentRepositoryImpl) AddArtifact(ctx context.Context, incidentID string, artifact incident.Artifact) error {
	path := fmt.Sprintf("/incidents/%s/artifacts", incidentID)

	// Create simplified map for API request (only sending required fields)
	artifactData := map[string]interface{}{
		"name":         artifact.Name,
		"type":         artifact.Type,
		"size":         artifact.Size,
		"content_type": strings.TrimSpace(artifact.ContentHash), // Use ContentHash as ContentType
		"url":          artifact.StoragePath,                    // Use StoragePath as URL
	}

	// Create request body
	body, err := json.Marshal(artifactData)
	if err != nil {
		return fmt.Errorf("failed to marshal artifact: %w", err)
	}

	req, err := createRequest(ctx, http.MethodPost, r.baseURL+path, body, r.authHeaders)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return handleErrorResponse(resp)
	}

	return nil
}

// UpdateTags updates the tags for an incident
func (r *IncidentRepositoryImpl) UpdateTags(ctx context.Context, id string, tags []string) error {
	path := fmt.Sprintf("/incidents/%s/tags", id)

	// Create request body
	tagData := map[string][]string{"tags": tags}
	body, err := json.Marshal(tagData)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	req, err := createRequest(ctx, http.MethodPut, r.baseURL+path, body, r.authHeaders)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return handleErrorResponse(resp)
	}

	return nil
}

// Note: We're removing these obsolete conversion functions since they have type mismatches.
// We'll use the properly implemented ToDomainIncident and ToIncidentDTO functions from incident_dtos.go instead.

// Note: We have removed these additional conversion functions that were causing type conflicts.
// We'll rely exclusively on the properly implemented ToDomainIncident and ToIncidentDTO functions
// from incident_dtos.go, which correctly handle the type conversions between domain models and DTOs.
