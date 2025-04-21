// Package repository provides repository implementations for various domain entities.
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

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

	var incidentDTO incidentDTO
	if err := json.NewDecoder(resp.Body).Decode(&incidentDTO); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return incidentDTO.toDomain(), nil
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
	for i, incidentDTO := range incidentDTOs {
		incidents[i] = incidentDTO.toDomain()
	}

	return incidents, pagination, nil
}

// Create creates a new incident
func (r *IncidentRepositoryImpl) Create(ctx context.Context, inc *incident.Incident) error {
	// Convert domain entity to DTO
	incidentDTO := fromDomainIncident(inc)

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
	var createdIncidentDTO incidentDTO
	if err := json.NewDecoder(resp.Body).Decode(&createdIncidentDTO); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Update the ID of the original incident object
	inc.ID = createdIncidentDTO.ID

	return nil
}

// Update updates an existing incident
func (r *IncidentRepositoryImpl) Update(ctx context.Context, inc *incident.Incident) error {
	// Convert domain entity to DTO
	incidentDTO := fromDomainIncident(inc)

	// Create request body
	body, err := json.Marshal(incidentDTO)
	if err != nil {
		return fmt.Errorf("failed to marshal incident: %w", err)
	}

	path := fmt.Sprintf("/incidents/%s", inc.ID)
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

	// Create request body
	iocDTO := indicatorOfCompromiseDTO{
		Type:        ioc.Type,
		Value:       ioc.Value,
		Description: ioc.Description,
		Source:      ioc.Source,
		Timestamp:   ioc.Timestamp.Format(time.RFC3339),
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

	// Create request body
	artifactDTO := artifactDTO{
		Name:        artifact.Name,
		Type:        artifact.Type,
		Size:        artifact.Size,
		Hash:        artifact.Hash,
		Path:        artifact.Path,
		CreatedAt:   artifact.CreatedAt.Format(time.RFC3339),
		Description: artifact.Description,
	}

	body, err := json.Marshal(artifactDTO)
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

// incidentDTO is the data transfer object for an incident
type incidentDTO struct {
	ID             string                     `json:"id"`
	OrganizationID int                        `json:"organization_id"`
	AgentID        string                     `json:"agent_id,omitempty"`
	Title          string                     `json:"title"`
	Description    string                     `json:"description"`
	Status         string                     `json:"status"`
	Severity       string                     `json:"severity"`
	Type           string                     `json:"type"`
	DetectedAt     string                     `json:"detected_at"`
	ResolvedAt     string                     `json:"resolved_at,omitempty"`
	CreatedAt      string                     `json:"created_at"`
	UpdatedAt      string                     `json:"updated_at"`
	AssignedTo     string                     `json:"assigned_to,omitempty"`
	Tags           []string                   `json:"tags,omitempty"`
	IOCs           []indicatorOfCompromiseDTO `json:"iocs,omitempty"`
	Artifacts      []artifactDTO              `json:"artifacts,omitempty"`
	Notes          []noteDTO                  `json:"notes,omitempty"`
}

type indicatorOfCompromiseDTO struct {
	Type        string `json:"type"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Source      string `json:"source"`
	Timestamp   string `json:"timestamp"`
}

type artifactDTO struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Size        int64  `json:"size"`
	Hash        string `json:"hash"`
	Path        string `json:"path"`
	CreatedAt   string `json:"created_at"`
	Description string `json:"description"`
}

type noteDTO struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// toDomain converts a DTO to a domain entity
func (dto *incidentDTO) toDomain() *incident.Incident {
	inc := &incident.Incident{
		ID:             dto.ID,
		OrganizationID: dto.OrganizationID,
		AgentID:        dto.AgentID,
		Title:          dto.Title,
		Description:    dto.Description,
		Status:         incident.IncidentStatus(dto.Status),
		Severity:       incident.IncidentSeverity(dto.Severity),
		Type:           incident.IncidentType(dto.Type),
		AssignedTo:     dto.AssignedTo,
		Tags:           dto.Tags,
	}

	// Parse timestamps
	if dto.DetectedAt != "" {
		inc.DetectedAt, _ = parseTime(dto.DetectedAt)
	}
	if dto.ResolvedAt != "" {
		resolvedAt, _ := parseTime(dto.ResolvedAt)
		inc.ResolvedAt = &resolvedAt
	}
	if dto.CreatedAt != "" {
		inc.CreatedAt, _ = parseTime(dto.CreatedAt)
	}
	if dto.UpdatedAt != "" {
		inc.UpdatedAt, _ = parseTime(dto.UpdatedAt)
	}

	// Convert IOCs
	inc.IOCs = make([]incident.IndicatorOfCompromise, len(dto.IOCs))
	for i, iocDTO := range dto.IOCs {
		timestamp, _ := parseTime(iocDTO.Timestamp)
		inc.IOCs[i] = incident.IndicatorOfCompromise{
			Type:        iocDTO.Type,
			Value:       iocDTO.Value,
			Description: iocDTO.Description,
			Source:      iocDTO.Source,
			Timestamp:   timestamp,
		}
	}

	// Convert Artifacts
	inc.Artifacts = make([]incident.Artifact, len(dto.Artifacts))
	for i, artifactDTO := range dto.Artifacts {
		createdAt, _ := parseTime(artifactDTO.CreatedAt)
		inc.Artifacts[i] = incident.Artifact{
			Name:        artifactDTO.Name,
			Type:        artifactDTO.Type,
			Size:        artifactDTO.Size,
			Hash:        artifactDTO.Hash,
			Path:        artifactDTO.Path,
			CreatedAt:   createdAt,
			Description: artifactDTO.Description,
		}
	}

	// Convert Notes
	inc.Notes = make([]incident.Note, len(dto.Notes))
	for i, noteDTO := range dto.Notes {
		createdAt, _ := parseTime(noteDTO.CreatedAt)
		updatedAt, _ := parseTime(noteDTO.UpdatedAt)
		inc.Notes[i] = incident.Note{
			ID:        noteDTO.ID,
			Content:   noteDTO.Content,
			CreatedBy: noteDTO.CreatedBy,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
	}

	return inc
}

// fromDomainIncident converts a domain entity to a DTO
func fromDomainIncident(inc *incident.Incident) *incidentDTO {
	dto := &incidentDTO{
		ID:             inc.ID,
		OrganizationID: inc.OrganizationID,
		AgentID:        inc.AgentID,
		Title:          inc.Title,
		Description:    inc.Description,
		Status:         string(inc.Status),
		Severity:       string(inc.Severity),
		Type:           string(inc.Type),
		DetectedAt:     inc.DetectedAt.Format(time.RFC3339),
		AssignedTo:     inc.AssignedTo,
		Tags:           inc.Tags,
	}

	if inc.ResolvedAt != nil {
		dto.ResolvedAt = inc.ResolvedAt.Format(time.RFC3339)
	}

	// Convert IOCs
	dto.IOCs = make([]indicatorOfCompromiseDTO, len(inc.IOCs))
	for i, ioc := range inc.IOCs {
		dto.IOCs[i] = indicatorOfCompromiseDTO{
			Type:        ioc.Type,
			Value:       ioc.Value,
			Description: ioc.Description,
			Source:      ioc.Source,
			Timestamp:   ioc.Timestamp.Format(time.RFC3339),
		}
	}

	// Convert Artifacts
	dto.Artifacts = make([]artifactDTO, len(inc.Artifacts))
	for i, artifact := range inc.Artifacts {
		dto.Artifacts[i] = artifactDTO{
			Name:        artifact.Name,
			Type:        artifact.Type,
			Size:        artifact.Size,
			Hash:        artifact.Hash,
			Path:        artifact.Path,
			CreatedAt:   artifact.CreatedAt.Format(time.RFC3339),
			Description: artifact.Description,
		}
	}

	// Convert Notes
	dto.Notes = make([]noteDTO, len(inc.Notes))
	for i, note := range inc.Notes {
		dto.Notes[i] = noteDTO{
			ID:        note.ID,
			Content:   note.Content,
			CreatedBy: note.CreatedBy,
			CreatedAt: note.CreatedAt.Format(time.RFC3339),
			UpdatedAt: note.UpdatedAt.Format(time.RFC3339),
		}
	}

	return dto
}
