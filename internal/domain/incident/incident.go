// Package incident contains the domain model for Huntress incidents.
package incident

import (
	"time"

	"github.com/google/uuid"
)

// Status represents the status of an incident.
type Status string

// Severity represents the severity level of an incident.
type Severity string

// Type represents the type/category of an incident.
type Type string

// Incident status constants
const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusResolved   Status = "resolved"
	StatusClosed     Status = "closed"
)

// Incident severity constants
const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
)

// Incident type constants
const (
	TypeMalware      Type = "malware"
	TypeRansomware   Type = "ransomware"
	TypePhishing     Type = "phishing"
	TypeUnauthorized Type = "unauthorized_access"
	TypeOther        Type = "other"
)

// IndicatorOfCompromise (IOC) represents evidence of a potential security breach
type IndicatorOfCompromise struct {
	ID          uuid.UUID `json:"id"`
	IncidentID  uuid.UUID `json:"incident_id"`
	Type        string    `json:"type"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	Source      string    `json:"source"`    // Added field based on repository usage
	Timestamp   time.Time `json:"timestamp"` // Added field based on repository usage
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Artifact represents a file or evidence associated with an incident
type Artifact struct {
	ID          uuid.UUID `json:"id"`
	IncidentID  uuid.UUID `json:"incident_id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Size        int64     `json:"size"`
	Hash        string    `json:"hash"`        // Added field based on repository usage
	Path        string    `json:"path"`        // Added field based on repository usage
	Description string    `json:"description"` // Added field based on repository usage
	ContentHash string    `json:"content_hash"`
	StoragePath string    `json:"storage_path"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Note represents a comment or observation added to an incident
type Note struct {
	ID         uuid.UUID `json:"id"`
	IncidentID uuid.UUID `json:"incident_id"`
	Content    string    `json:"content"`
	CreatedBy  string    `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Incident represents a security incident within the system
type Incident struct {
	ID             uuid.UUID               `json:"id"`
	OrganizationID uuid.UUID               `json:"organization_id"`
	AgentID        string                  `json:"agent_id"` // Added field based on repository usage
	Title          string                  `json:"title"`
	Description    string                  `json:"description"`
	Status         Status                  `json:"status"`
	Severity       Severity                `json:"severity"`
	Type           Type                    `json:"type"`
	AssignedTo     string                  `json:"assigned_to"`
	Reporter       string                  `json:"reporter"`
	Notes          []Note                  `json:"notes"`
	IOCs           []IndicatorOfCompromise `json:"iocs"`
	Artifacts      []Artifact              `json:"artifacts"`
	Tags           []string                `json:"tags"`
	DetectedAt     time.Time               `json:"detected_at"` // Added field based on repository usage
	CreatedAt      time.Time               `json:"created_at"`
	UpdatedAt      time.Time               `json:"updated_at"`
	ResolvedAt     *time.Time              `json:"resolved_at,omitempty"`
}
