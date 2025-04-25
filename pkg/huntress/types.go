// Package huntress provides a client for the Huntress API
package huntress

import "time"

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

// ----- Organization Types -----

// OrganizationCreateParams contains parameters for creating an organization
type OrganizationCreateParams struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Industry    string                 `json:"industry,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	Address     *AddressParams         `json:"address,omitempty"`
	ContactInfo *ContactInfoParams     `json:"contact_info,omitempty"`
}

// OrganizationCreateInput is an alias for OrganizationCreateParams for backward compatibility
type OrganizationCreateInput = OrganizationCreateParams

// OrganizationUpdateParams contains parameters for updating an organization
type OrganizationUpdateParams struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Industry    string                 `json:"industry,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	Address     *AddressParams         `json:"address,omitempty"`
	ContactInfo *ContactInfoParams     `json:"contact_info,omitempty"`
}

// OrganizationListOptions contains options for listing organizations
type OrganizationListOptions struct {
	ListParams
	Status string   `url:"status,omitempty"`
	Search string   `url:"search,omitempty"`
	Tags   []string `url:"tags,omitempty"`
}

// ListOrganizationsParams is an alias for OrganizationListOptions
type ListOrganizationsParams = OrganizationListOptions

// ----- Agent Types -----

// AgentListOptions contains options for listing agents
type AgentListOptions struct {
	ListParams
	OrganizationID int        `url:"organization_id,omitempty"`
	Status         string     `url:"status,omitempty"`
	Platform       string     `url:"platform,omitempty"`
	Hostname       string     `url:"hostname,omitempty"`
	IPAddress      string     `url:"ip_address,omitempty"`
	Version        string     `url:"version,omitempty"`
	Search         string     `url:"search,omitempty"`
	Tags           []string   `url:"tags,omitempty"`
	LastSeenSince  *time.Time `url:"last_seen_since,omitempty"`
}

// AgentStatus represents the current status of an agent
type AgentStatus string

const (
	// AgentStatusOnline indicates an online agent
	AgentStatusOnline AgentStatus = "online"
	// AgentStatusOffline indicates an offline agent
	AgentStatusOffline AgentStatus = "offline"
	// AgentStatusPending indicates a pending agent
	AgentStatusPending AgentStatus = "pending"
	// AgentStatusUnknown indicates an agent with unknown status
	AgentStatusUnknown AgentStatus = "unknown"
)

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
	Organization   int        `url:"organization_id,omitempty"`
	AgentID        string     `url:"agent_id,omitempty"`
	Status         string     `url:"status,omitempty"`
	Severity       string     `url:"severity,omitempty"`
	Type           string     `url:"type,omitempty"`
	Search         string     `url:"search,omitempty"`
	AssignedTo     string     `url:"assigned_to,omitempty"`
	Tags           []string   `url:"tags,omitempty"`
	DetectedAfter  *time.Time `url:"detected_after,omitempty"`
	DetectedBefore *time.Time `url:"detected_before,omitempty"`
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
