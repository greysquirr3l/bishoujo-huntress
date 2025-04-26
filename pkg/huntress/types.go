// Package huntress provides a client for the Huntress API
package huntress

import (
	"fmt"
	"time"
)

// This file serves as the centralized type registry for the Huntress API client.
// It imports and re-exports types from their canonical sources to avoid redeclarations
// while maintaining backward compatibility.

// Core shared types are defined here, while more specialized types are imported
// from their respective packages.

// ----- Common Types -----

// ListParams contains common pagination parameters for list operations
// This is the canonical definition that should be used across all services
type ListParams struct {
	Page     int    `url:"page,omitempty"`
	PerPage  int    `url:"per_page,omitempty"`
	SortBy   string `url:"sort_by,omitempty"`
	SortDesc bool   `url:"sort_desc,omitempty"`
}

// ListOptions is an alias for ListParams for backward compatibility
type ListOptions = ListParams

// Pagination represents pagination information returned by the API
type Pagination struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalPages  int `json:"total_pages"`
	TotalItems  int `json:"total_items"`
}

// ----- ENUM RE-EXPORTS -----

// OrganizationStatus is the status of an organization (active, inactive, pending)
type OrganizationStatus string

const (
	// OrganizationStatusActive indicates the organization is active.
	OrganizationStatusActive OrganizationStatus = "active"
	// OrganizationStatusInactive indicates the organization is inactive.
	OrganizationStatusInactive OrganizationStatus = "inactive"
	// OrganizationStatusPending indicates the organization is pending.
	OrganizationStatusPending OrganizationStatus = "pending"
)

// UserRole is the role of a user in an organization
// (admin, manager, viewer)
type UserRole string

const (
	// UserRoleAdmin indicates the user is an admin.
	UserRoleAdmin UserRole = "admin"
	// UserRoleManager indicates the user is a manager.
	UserRoleManager UserRole = "manager"
	// UserRoleViewer indicates the user is a viewer.
	UserRoleViewer UserRole = "viewer"
)

// AgentStatus is the status of an agent (online, offline, pending, unknown)
type AgentStatus string

const (
	// AgentStatusOnline indicates the agent is online.
	AgentStatusOnline AgentStatus = "online"
	// AgentStatusOffline indicates the agent is offline.
	AgentStatusOffline AgentStatus = "offline"
	// AgentStatusPending indicates the agent is pending.
	AgentStatusPending AgentStatus = "pending"
	// AgentStatusUnknown indicates the agent status is unknown.
	AgentStatusUnknown AgentStatus = "unknown"
)

// AgentPlatform is the platform of an agent (windows, mac, linux)
type AgentPlatform string

const (
	// AgentPlatformWindows indicates the agent is running on Windows.
	AgentPlatformWindows AgentPlatform = "windows"
	// AgentPlatformMac indicates the agent is running on macOS.
	AgentPlatformMac AgentPlatform = "mac"
	// AgentPlatformLinux indicates the agent is running on Linux.
	AgentPlatformLinux AgentPlatform = "linux"
)

// IncidentStatus is the status of an incident (new, in_progress, resolved, closed)
type IncidentStatus string

const (
	// IncidentStatusNew indicates the incident is new.
	IncidentStatusNew IncidentStatus = "new"
	// IncidentStatusInProgress indicates the incident is in progress.
	IncidentStatusInProgress IncidentStatus = "in_progress"
	// IncidentStatusResolved indicates the incident is resolved.
	IncidentStatusResolved IncidentStatus = "resolved"
	// IncidentStatusClosed indicates the incident is closed.
	IncidentStatusClosed IncidentStatus = "closed"
)

// IncidentSeverity is the severity of an incident (critical, high, medium, low)
type IncidentSeverity string

const (
	// IncidentSeverityCritical indicates a critical severity incident.
	IncidentSeverityCritical IncidentSeverity = "critical"
	// IncidentSeverityHigh indicates a high severity incident.
	IncidentSeverityHigh IncidentSeverity = "high"
	// IncidentSeverityMedium indicates a medium severity incident.
	IncidentSeverityMedium IncidentSeverity = "medium"
	// IncidentSeverityLow indicates a low severity incident.
	IncidentSeverityLow IncidentSeverity = "low"
)

// IncidentType is the type/category of an incident
// (malware, ransomware, phishing, unauthorized_access, other)
type IncidentType string

const (
	// IncidentTypeMalware indicates a malware incident.
	IncidentTypeMalware IncidentType = "malware"
	// IncidentTypeRansomware indicates a ransomware incident.
	IncidentTypeRansomware IncidentType = "ransomware"
	// IncidentTypePhishing indicates a phishing incident.
	IncidentTypePhishing IncidentType = "phishing"
	// IncidentTypeUnauthorized indicates an unauthorized access incident.
	IncidentTypeUnauthorized IncidentType = "unauthorized_access"
	// IncidentTypeOther indicates an incident of another type.
	IncidentTypeOther IncidentType = "other"
)

// ----- Organization Types -----

// OrganizationCreateParams contains parameters for creating an organization
type OrganizationCreateParams struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Status      OrganizationStatus     `json:"status,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Industry    string                 `json:"industry,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	Address     *AddressParams         `json:"address,omitempty"`
	ContactInfo *ContactInfoParams     `json:"contact_info,omitempty"`
}

// Validate checks if the OrganizationCreateParams are valid
func (p *OrganizationCreateParams) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("organization name is required")
	}
	if p.Status != "" && p.Status != OrganizationStatusActive && p.Status != OrganizationStatusInactive && p.Status != OrganizationStatusPending {
		return fmt.Errorf("invalid organization status: %s", p.Status)
	}
	return nil
}

// OrganizationUpdateParams contains parameters for updating an organization
type OrganizationUpdateParams struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Status      OrganizationStatus     `json:"status,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Industry    string                 `json:"industry,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	Address     *AddressParams         `json:"address,omitempty"`
	ContactInfo *ContactInfoParams     `json:"contact_info,omitempty"`
}

// Validate checks if the OrganizationUpdateParams are valid
func (p *OrganizationUpdateParams) Validate() error {
	if p.Status != "" && p.Status != OrganizationStatusActive && p.Status != OrganizationStatusInactive && p.Status != OrganizationStatusPending {
		return fmt.Errorf("invalid organization status: %s", p.Status)
	}
	return nil
}

// OrganizationListOptions contains options for listing organizations
type OrganizationListOptions struct {
	ListParams
	Status OrganizationStatus `url:"status,omitempty"`
	Search string             `url:"search,omitempty"`
	Tags   []string           `url:"tags,omitempty"`
}

// ListOrganizationsParams is an alias for OrganizationListOptions
type ListOrganizationsParams = OrganizationListOptions

// ----- Agent Types -----

// AgentListOptions contains options for listing agents
type AgentListOptions struct {
	ListParams
	OrganizationID int           `url:"organization_id,omitempty"`
	Status         AgentStatus   `url:"status,omitempty"`
	Platform       AgentPlatform `url:"platform,omitempty"`
	Hostname       string        `url:"hostname,omitempty"`
	IPAddress      string        `url:"ip_address,omitempty"`
	Version        string        `url:"version,omitempty"`
	Search         string        `url:"search,omitempty"`
	Tags           []string      `url:"tags,omitempty"`
	LastSeenSince  *time.Time    `url:"last_seen_since,omitempty"`
}

// AgentStatistics represents statistics for an agent
type AgentStatistics struct {
	TotalDetections  int            `json:"total_detections"`
	LastDetection    time.Time      `json:"last_detection,omitempty"`
	UpTime           float64        `json:"up_time"` // percentage
	LastUpdated      time.Time      `json:"last_updated"`
	DetectionsByType map[string]int `json:"detections_by_type,omitempty"`
}

// ----- Incident Types -----

// IncidentListOptions contains options for listing incidents
type IncidentListOptions struct {
	ListOptions
	Organization   int              `url:"organization_id,omitempty"`
	AgentID        string           `url:"agent_id,omitempty"`
	Status         IncidentStatus   `url:"status,omitempty"`
	Severity       IncidentSeverity `url:"severity,omitempty"`
	Type           IncidentType     `url:"type,omitempty"`
	Search         string           `url:"search,omitempty"`
	AssignedTo     string           `url:"assigned_to,omitempty"`
	Tags           []string         `url:"tags,omitempty"`
	DetectedAfter  *time.Time       `url:"detected_after,omitempty"`
	DetectedBefore *time.Time       `url:"detected_before,omitempty"`
}

// ListIncidentsParams is an alias for IncidentListOptions
type ListIncidentsParams = IncidentListOptions

// ----- Helper Types -----

// AddressParams represents physical address parameters
type AddressParams struct {
	Street1 string `json:"street1,omitempty"`
	Street2 string `json:"street2,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	ZipCode string `json:"zip_code,omitempty"`
	Country string `json:"country,omitempty"`
}

// ContactInfoParams represents contact information parameters
type ContactInfoParams struct {
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Title       string `json:"title,omitempty"`
}

// Validate checks if the AgentListOptions are valid
func (p *AgentListOptions) Validate() error {
	if p.Status != "" && p.Status != AgentStatusOnline && p.Status != AgentStatusOffline && p.Status != AgentStatusPending && p.Status != AgentStatusUnknown {
		return fmt.Errorf("invalid agent status: %s", p.Status)
	}
	if p.Platform != "" && p.Platform != AgentPlatformWindows && p.Platform != AgentPlatformMac && p.Platform != AgentPlatformLinux {
		return fmt.Errorf("invalid agent platform: %s", p.Platform)
	}
	return nil
}

// Validate checks if the IncidentListOptions are valid
func (p *IncidentListOptions) Validate() error {
	if p.Status != "" && p.Status != IncidentStatusNew && p.Status != IncidentStatusInProgress && p.Status != IncidentStatusResolved && p.Status != IncidentStatusClosed {
		return fmt.Errorf("invalid incident status: %s", p.Status)
	}
	if p.Severity != "" && p.Severity != IncidentSeverityCritical && p.Severity != IncidentSeverityHigh && p.Severity != IncidentSeverityMedium && p.Severity != IncidentSeverityLow {
		return fmt.Errorf("invalid incident severity: %s", p.Severity)
	}
	if p.Type != "" && p.Type != IncidentTypeMalware && p.Type != IncidentTypeRansomware && p.Type != IncidentTypePhishing && p.Type != IncidentTypeUnauthorized && p.Type != IncidentTypeOther {
		return fmt.Errorf("invalid incident type: %s", p.Type)
	}
	return nil
}

// Validate checks if the UserCreateParams are valid
func (p *UserCreateParams) Validate() error {
	if p.Email == "" {
		return fmt.Errorf("user email is required")
	}
	if p.Role != "" && p.Role != UserRoleAdmin && p.Role != UserRoleManager && p.Role != UserRoleViewer {
		return fmt.Errorf("invalid user role: %s", p.Role)
	}
	for _, r := range p.Roles {
		if r != UserRoleAdmin && r != UserRoleManager && r != UserRoleViewer {
			return fmt.Errorf("invalid user role in roles: %s", r)
		}
	}
	return nil
}

// Validate checks if the UserUpdateParams are valid
func (p *UserUpdateParams) Validate() error {
	if p.Role != "" && p.Role != UserRoleAdmin && p.Role != UserRoleManager && p.Role != UserRoleViewer {
		return fmt.Errorf("invalid user role: %s", p.Role)
	}
	for _, r := range p.Roles {
		if r != UserRoleAdmin && r != UserRoleManager && r != UserRoleViewer {
			return fmt.Errorf("invalid user role in roles: %s", r)
		}
	}
	return nil
}
