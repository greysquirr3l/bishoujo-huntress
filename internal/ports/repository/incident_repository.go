package repository

import (
	"context"
	"time"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/incident"
)

// IncidentFilters defines filters for incident queries
type IncidentFilters struct {
	OrganizationID *int
	AgentID        string
	Status         *incident.IncidentStatus
	Severity       *incident.IncidentSeverity
	Type           *incident.IncidentType
	Search         string
	AssignedTo     string
	DetectedAfter  *time.Time
	DetectedBefore *time.Time
	Tags           []string
	Page           int
	Limit          int
	OrderBy        []OrderBy
}

// IncidentRepository defines the repository interface for incidents
type IncidentRepository interface {
	// Get retrieves an incident by ID
	Get(ctx context.Context, id string) (*incident.Incident, error)

	// List retrieves multiple incidents based on filters
	List(ctx context.Context, filters IncidentFilters) ([]*incident.Incident, Pagination, error)

	// Create creates a new incident
	Create(ctx context.Context, incident *incident.Incident) error

	// Update updates an existing incident
	Update(ctx context.Context, incident *incident.Incident) error

	// Delete deletes an incident by ID (if allowed)
	Delete(ctx context.Context, id string) error

	// UpdateStatus updates the status of an incident
	UpdateStatus(ctx context.Context, id string, status incident.IncidentStatus) error

	// UpdateAssignee updates the assignee of an incident
	UpdateAssignee(ctx context.Context, id string, assignee string) error

	// ListByOrganization retrieves all incidents for an organization
	ListByOrganization(ctx context.Context, organizationID int, filters IncidentFilters) ([]*incident.Incident, Pagination, error)

	// ListByAgent retrieves all incidents for a specific agent
	ListByAgent(ctx context.Context, agentID string, filters IncidentFilters) ([]*incident.Incident, Pagination, error)

	// AddNote adds a note to an incident
	AddNote(ctx context.Context, incidentID string, note incident.Note) error

	// AddIOC adds an indicator of compromise to an incident
	AddIOC(ctx context.Context, incidentID string, ioc incident.IndicatorOfCompromise) error

	// AddArtifact adds an artifact to an incident
	AddArtifact(ctx context.Context, incidentID string, artifact incident.Artifact) error

	// UpdateTags updates the tags for an incident
	UpdateTags(ctx context.Context, id string, tags []string) error
}
