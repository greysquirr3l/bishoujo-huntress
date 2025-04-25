// Package repository provides repository implementations for various domain entities.
package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/incident"
)

// incidentDTO represents the data transfer object for incidents
type incidentDTO struct {
	ID             string                     `json:"id"`
	Title          string                     `json:"title"`
	Description    string                     `json:"description"`
	OrganizationID int                        `json:"organization_id"`
	AgentID        string                     `json:"agent_id,omitempty"`
	Status         string                     `json:"status"`
	Severity       string                     `json:"severity"`
	Type           string                     `json:"type"`
	DetectedAt     time.Time                  `json:"detected_at"`
	ResolvedAt     *time.Time                 `json:"resolved_at,omitempty"`
	CreatedAt      time.Time                  `json:"created_at"`
	UpdatedAt      time.Time                  `json:"updated_at"`
	Tags           []string                   `json:"tags,omitempty"`
	AssignedTo     string                     `json:"assigned_to,omitempty"`
	Notes          []noteDTO                  `json:"notes,omitempty"`
	Artifacts      []artifactDTO              `json:"artifacts,omitempty"`
	IOCs           []indicatorOfCompromiseDTO `json:"iocs,omitempty"`
}

// noteDTO represents the data transfer object for incident notes
type noteDTO struct {
	ID         string    `json:"id"`
	IncidentID string    `json:"incident_id"`
	Content    string    `json:"content"`
	CreatedBy  string    `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
}

// artifactDTO represents the data transfer object for incident artifacts
type artifactDTO struct {
	ID          string    `json:"id"`
	IncidentID  string    `json:"incident_id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	ContentType string    `json:"content_type,omitempty"`
	Size        int64     `json:"size,omitempty"`
	URL         string    `json:"url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// indicatorOfCompromiseDTO represents the data transfer object for IOCs
type indicatorOfCompromiseDTO struct {
	ID          string    `json:"id"`
	IncidentID  string    `json:"incident_id"`
	Type        string    `json:"type"`
	Value       string    `json:"value"`
	Description string    `json:"description,omitempty"`
	Source      string    `json:"source,omitempty"`
	Timestamp   time.Time `json:"timestamp,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToDomainIncident converts an incidentDTO to a domain Incident entity
func ToDomainIncident(dto *incidentDTO) (*incident.Incident, error) {
	// Parse UUID from string ID
	incidentID, err := uuid.Parse(dto.ID)
	if err != nil {
		return nil, err
	}

	// Parse organization ID UUID
	orgID, err := uuid.Parse(fmt.Sprintf("%d", dto.OrganizationID))
	if err != nil {
		// Use a placeholder UUID if conversion fails
		orgID = uuid.New()
	}

	// Create a domain incident with properly converted types
	inc := &incident.Incident{
		ID:             incidentID,
		OrganizationID: orgID,
		Title:          dto.Title,
		Description:    dto.Description,
		Status:         incident.IncidentStatus(dto.Status),
		Severity:       incident.IncidentSeverity(dto.Severity),
		Type:           incident.IncidentType(dto.Type),
		DetectedAt:     dto.DetectedAt,
		CreatedAt:      dto.CreatedAt,
		UpdatedAt:      dto.UpdatedAt,
		Tags:           dto.Tags,
		AssignedTo:     dto.AssignedTo,
		// AgentID is a string, no need for conversion
		AgentID: dto.AgentID,
	}

	// Set optional fields if they exist
	if dto.ResolvedAt != nil {
		inc.ResolvedAt = dto.ResolvedAt
	}

	// Convert notes
	if len(dto.Notes) > 0 {
		inc.Notes = make([]incident.Note, len(dto.Notes))
		for i, noteDto := range dto.Notes {
			noteID, err := uuid.Parse(noteDto.ID)
			if err != nil {
				return nil, err
			}

			inc.Notes[i] = incident.Note{
				ID:        noteID,
				Content:   noteDto.Content,
				CreatedBy: noteDto.CreatedBy,
				CreatedAt: noteDto.CreatedAt,
				UpdatedAt: noteDto.CreatedAt, // Using CreatedAt as a fallback since DTO doesn't have UpdatedAt
			}
		}
	}

	// Convert artifacts
	if len(dto.Artifacts) > 0 {
		inc.Artifacts = make([]incident.Artifact, len(dto.Artifacts))
		for i, artifactDto := range dto.Artifacts {
			artifactID, err := uuid.Parse(artifactDto.ID)
			if err != nil {
				return nil, err
			}

			incidentID, err := uuid.Parse(artifactDto.IncidentID)
			if err != nil {
				return nil, err
			}

			inc.Artifacts[i] = incident.Artifact{
				ID:         artifactID,
				IncidentID: incidentID,
				Name:       artifactDto.Name,
				Type:       artifactDto.Type,
				Size:       artifactDto.Size,
				// Use empty strings for fields not in the DTO
				Hash:        "",
				Path:        "",
				Description: "",
				ContentHash: "",
				StoragePath: "",
				CreatedAt:   artifactDto.CreatedAt,
				UpdatedAt:   artifactDto.CreatedAt, // Using CreatedAt as a fallback
			}
		}
	}

	// Convert IOCs
	if len(dto.IOCs) > 0 {
		inc.IOCs = make([]incident.IndicatorOfCompromise, len(dto.IOCs))
		for i, iocDto := range dto.IOCs {
			iocID, err := uuid.Parse(iocDto.ID)
			if err != nil {
				return nil, err
			}

			incidentID, err := uuid.Parse(iocDto.IncidentID)
			if err != nil {
				return nil, err
			}

			inc.IOCs[i] = incident.IndicatorOfCompromise{
				ID:          iocID,
				IncidentID:  incidentID,
				Type:        iocDto.Type,
				Value:       iocDto.Value,
				Description: iocDto.Description,
				Source:      iocDto.Source,
				Timestamp:   iocDto.Timestamp,
				CreatedAt:   iocDto.CreatedAt,
				UpdatedAt:   iocDto.CreatedAt, // Using CreatedAt as a fallback
			}
		}
	}

	return inc, nil
}

// ToIncidentDTO converts a domain Incident entity to an incidentDTO
func ToIncidentDTO(inc *incident.Incident) *incidentDTO {
	dto := &incidentDTO{
		ID:          inc.ID.String(),
		Title:       inc.Title,
		Description: inc.Description,
		// Convert UUID to int for organization ID - this is a simplification
		OrganizationID: 0, // In a real implementation, you'd map correctly
		Status:         string(inc.Status),
		Severity:       string(inc.Severity),
		Type:           string(inc.Type),
		DetectedAt:     inc.DetectedAt,
		CreatedAt:      inc.CreatedAt,
		UpdatedAt:      inc.UpdatedAt,
		Tags:           inc.Tags,
		AssignedTo:     inc.AssignedTo,
	}

	// Set optional fields if they exist
	if inc.ResolvedAt != nil {
		dto.ResolvedAt = inc.ResolvedAt
	}

	// AgentID is already a string in the domain model, direct assignment
	dto.AgentID = inc.AgentID

	// Convert notes
	if len(inc.Notes) > 0 {
		dto.Notes = make([]noteDTO, len(inc.Notes))
		for i, note := range inc.Notes {
			dto.Notes[i] = noteDTO{
				ID:         note.ID.String(),
				IncidentID: inc.ID.String(),
				Content:    note.Content,
				CreatedBy:  note.CreatedBy,
				CreatedAt:  note.CreatedAt,
			}
		}
	}

	// Convert artifacts
	if len(inc.Artifacts) > 0 {
		dto.Artifacts = make([]artifactDTO, len(inc.Artifacts))
		for i, artifact := range inc.Artifacts {
			dto.Artifacts[i] = artifactDTO{
				ID:         artifact.ID.String(),
				IncidentID: inc.ID.String(),
				Name:       artifact.Name,
				Type:       artifact.Type,
				// We don't have direct mapping for ContentType in domain model
				ContentType: "",
				Size:        artifact.Size,
				// We don't have URL in domain model, could use StoragePath instead if needed
				URL:       artifact.StoragePath,
				CreatedAt: artifact.CreatedAt,
			}
		}
	}

	// Convert IOCs
	if len(inc.IOCs) > 0 {
		dto.IOCs = make([]indicatorOfCompromiseDTO, len(inc.IOCs))
		for i, ioc := range inc.IOCs {
			dto.IOCs[i] = indicatorOfCompromiseDTO{
				ID:          ioc.ID.String(),
				IncidentID:  inc.ID.String(),
				Type:        ioc.Type,
				Value:       ioc.Value,
				Description: ioc.Description,
				Source:      ioc.Source,
				Timestamp:   ioc.Timestamp,
				CreatedAt:   ioc.CreatedAt,
			}
		}
	}

	return dto
}
