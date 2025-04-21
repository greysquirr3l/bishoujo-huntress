package incident

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrInvalidIncidentID indicates that the incident ID is invalid
	ErrInvalidIncidentID = errors.New("incident ID is invalid")
	// ErrInvalidTitle indicates that the incident title is invalid
	ErrInvalidTitle = errors.New("incident title is required")
	// ErrInvalidSeverity indicates that the incident severity is invalid
	ErrInvalidSeverity = errors.New("incident severity must be one of: low, medium, high, critical")
	// ErrInvalidStatus indicates that the incident status is invalid
	ErrInvalidStatus = errors.New("incident status must be one of: open, investigating, contained, resolved")
)

// Severity represents the severity level of an incident
type Severity string

const (
	// SeverityLow represents a low severity incident
	SeverityLow Severity = "low"
	// SeverityMedium represents a medium severity incident
	SeverityMedium Severity = "medium"
	// SeverityHigh represents a high severity incident
	SeverityHigh Severity = "high"
	// SeverityCritical represents a critical severity incident
	SeverityCritical Severity = "critical"
)

// Status represents the current status of an incident
type Status string

const (
	// StatusOpen indicates the incident is open and unresolved
	StatusOpen Status = "open"
	// StatusInvestigating indicates the incident is being investigated
	StatusInvestigating Status = "investigating"
	// StatusContained indicates the incident has been contained
	StatusContained Status = "contained"
	// StatusResolved indicates the incident has been resolved
	StatusResolved Status = "resolved"
)

// Incident represents a security incident detected by Huntress
type Incident struct {
	ID          uuid.UUID
	Title       string
	Description string
	Severity    Severity
	Status      Status
	DetectedAt  time.Time
	ResolvedAt  *time.Time
	AgentID     uuid.UUID
	OrgID       uuid.UUID
	Indicators  []Indicator
	Notes       []Note
}

// Indicator represents an indicator of compromise associated with an incident
type Indicator struct {
	Type  string
	Value string
}

// Note represents a note attached to an incident
type Note struct {
	Content   string
	CreatedAt time.Time
	CreatedBy string
}

// NewIncident creates a new incident with the given parameters
func NewIncident(
	id uuid.UUID,
	title string,
	description string,
	severity Severity,
	agentID uuid.UUID,
	orgID uuid.UUID,
	detectedAt time.Time,
) (*Incident, error) {
	incident := &Incident{
		ID:          id,
		Title:       title,
		Description: description,
		Severity:    severity,
		Status:      StatusOpen,
		DetectedAt:  detectedAt,
		AgentID:     agentID,
		OrgID:       orgID,
		Indicators:  []Indicator{},
		Notes:       []Note{},
	}

	if err := incident.Validate(); err != nil {
		return nil, err
	}

	return incident, nil
}

// Validate ensures that the incident has valid data
func (i *Incident) Validate() error {
	if i.ID == uuid.Nil {
		return ErrInvalidIncidentID
	}

	if i.Title == "" {
		return ErrInvalidTitle
	}

	if err := i.validateSeverity(); err != nil {
		return err
	}

	if err := i.validateStatus(); err != nil {
		return err
	}

	return nil
}

// AddNote adds a new note to the incident
func (i *Incident) AddNote(content string, createdBy string) {
	note := Note{
		Content:   content,
		CreatedAt: time.Now().UTC(),
		CreatedBy: createdBy,
	}
	i.Notes = append(i.Notes, note)
}

// AddIndicator adds a new indicator of compromise to the incident
func (i *Incident) AddIndicator(indicatorType string, value string) {
	indicator := Indicator{
		Type:  indicatorType,
		Value: value,
	}
	i.Indicators = append(i.Indicators, indicator)
}

// UpdateStatus updates the status of the incident
func (i *Incident) UpdateStatus(status Status) error {
	if err := i.validateStatusTransition(status); err != nil {
		return err
	}

	i.Status = status

	// If the incident is being resolved, set the resolved time
	if status == StatusResolved {
		now := time.Now().UTC()
		i.ResolvedAt = &now
	}

	return nil
}

// UpdateSeverity updates the severity of the incident
func (i *Incident) UpdateSeverity(severity Severity) error {
	if !isValidSeverity(severity) {
		return ErrInvalidSeverity
	}
	i.Severity = severity
	return nil
}

func (i *Incident) validateSeverity() error {
	if !isValidSeverity(i.Severity) {
		return ErrInvalidSeverity
	}
	return nil
}

func (i *Incident) validateStatus() error {
	if !isValidStatus(i.Status) {
		return ErrInvalidStatus
	}
	return nil
}

func (i *Incident) validateStatusTransition(newStatus Status) error {
	if !isValidStatus(newStatus) {
		return ErrInvalidStatus
	}

	// Add validation logic for valid status transitions if needed
	// For example, preventing resolved incidents from being reopened
	if i.Status == StatusResolved && newStatus != StatusResolved {
		return errors.New("cannot change status of resolved incident")
	}

	return nil
}

func isValidSeverity(severity Severity) bool {
	switch severity {
	case SeverityLow, SeverityMedium, SeverityHigh, SeverityCritical:
		return true
	default:
		return false
	}
}

func isValidStatus(status Status) bool {
	switch status {
	case StatusOpen, StatusInvestigating, StatusContained, StatusResolved:
		return true
	default:
		return false
	}
}
