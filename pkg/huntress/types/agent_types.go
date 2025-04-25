package types

import "time"

// Agent represents a Huntress agent installed on an endpoint
type Agent struct {
	ID               string     `json:"id"`
	OrganizationID   int        `json:"organization_id"`
	Version          string     `json:"version"`
	Hostname         string     `json:"hostname"`
	IPV4Address      string     `json:"ipv4_address,omitempty"`
	MACAddress       string     `json:"mac_address,omitempty"`
	Platform         string     `json:"platform"`
	OS               string     `json:"os"`
	OSVersion        string     `json:"os_version,omitempty"`
	Status           string     `json:"status"`
	LastSeen         time.Time  `json:"last_seen"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	ExternalID       string     `json:"external_id,omitempty"`
	UserInfo         UserInfo   `json:"user_info,omitempty"`
	SystemInfo       SystemInfo `json:"system_info,omitempty"`
	EncryptionStatus string     `json:"encryption_status,omitempty"`
	Tags             []string   `json:"tags,omitempty"`
}

// UserInfo contains user information related to an agent
type UserInfo struct {
	Username  string    `json:"username,omitempty"`
	Domain    string    `json:"domain,omitempty"`
	IsAdmin   bool      `json:"is_admin,omitempty"`
	LastLogon time.Time `json:"last_logon,omitempty"`
}

// SystemInfo contains system information related to an agent
type SystemInfo struct {
	Manufacturer  string `json:"manufacturer,omitempty"`
	Model         string `json:"model,omitempty"`
	TotalRAM      int64  `json:"total_ram,omitempty"`
	DiskSize      int64  `json:"disk_size,omitempty"`
	ProcessorInfo string `json:"processor_info,omitempty"`
	BIOSVersion   string `json:"bios_version,omitempty"`
}

// AgentUpdateInput contains the fields that can be updated on an agent
type AgentUpdateInput struct {
	Tags []string `json:"tags,omitempty"`
}

// AgentStatistics contains statistics for an agent
type AgentStatistics struct {
	IncidentCount     int   `json:"incident_count"`
	LastActivityTime  int64 `json:"last_activity_time,omitempty"`
	ProcessCount      int   `json:"process_count,omitempty"`
	ServiceCount      int   `json:"service_count,omitempty"`
	InstalledSoftware int   `json:"installed_software,omitempty"`
}

// AgentListOptions contains options for listing agents
type AgentListOptions struct {
	// Pagination
	Page    int `url:"page,omitempty"`
	PerPage int `url:"per_page,omitempty"`

	// Filtering
	Status         string    `url:"status,omitempty"`
	Platform       string    `url:"platform,omitempty"`
	Search         string    `url:"search,omitempty"`
	OrganizationID int       `url:"organization_id,omitempty"`
	LastSeenSince  time.Time `url:"last_seen_since,omitempty"`
	LastSeenBefore time.Time `url:"last_seen_before,omitempty"`

	// Sorting
	SortBy   string `url:"sort_by,omitempty"`
	SortDesc bool   `url:"sort_desc,omitempty"`
}

// AgentList represents a paginated list of agents
type AgentList struct {
	Agents     []Agent     `json:"agents"`
	Pagination *Pagination `json:"pagination"`
}
