package incident

import (
	"time"

	"github.com/google/uuid"
)

// IncidentStatus represents the status of an incident
type IncidentStatus string

// Incident status constants
const (
	StatusNew        IncidentStatus = "new"
	StatusInProgress IncidentStatus = "in_progress"
	StatusResolved   IncidentStatus = "resolved"
	StatusClosed     IncidentStatus = "closed"
)

// IncidentSeverity represents the severity level of an incident
type IncidentSeverity string

// Incident severity constants
const (
	SeverityCritical IncidentSeverity = "critical"
	SeverityHigh     IncidentSeverity = "high"
	SeverityMedium   IncidentSeverity = "medium"
	SeverityLow      IncidentSeverity = "low"
)

// IncidentType represents the type of an incident
type IncidentType string

// Incident type constants
const (
	TypeMalware      IncidentType = "malware"
	TypeRansomware   IncidentType = "ransomware"
	TypePhishing     IncidentType = "phishing"
	TypeUnauthorized IncidentType = "unauthorized_access"
	TypeOther        IncidentType = "other"
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
	Status         IncidentStatus          `json:"status"`
	Severity       IncidentSeverity        `json:"severity"`
	Type           IncidentType            `json:"type"`
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
