// Package huntress provides a client for the Huntress API
package huntress

import "time"

// ----- Account Types -----

// AccountUpdateParams contains parameters for updating an account
type AccountUpdateParams struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	ContactInfo *ContactInfoParams     `json:"contact_info,omitempty"`
}

// ----- User Types -----

// UserCreateParams contains parameters for creating a user
type UserCreateParams struct {
	Email     string     `json:"email"`
	FirstName string     `json:"first_name,omitempty"`
	LastName  string     `json:"last_name,omitempty"`
	Role      UserRole   `json:"role,omitempty"`
	Roles     []UserRole `json:"roles,omitempty"`
	Active    bool       `json:"active,omitempty"`
}

// UserUpdateParams contains parameters for updating a user
type UserUpdateParams struct {
	Email     string     `json:"email,omitempty"`
	FirstName string     `json:"first_name,omitempty"`
	LastName  string     `json:"last_name,omitempty"`
	Role      UserRole   `json:"role,omitempty"`
	Roles     []UserRole `json:"roles,omitempty"`
	Active    *bool      `json:"active,omitempty"`
}

// ----- Report Types -----

// ReportGenerateInput contains parameters for generating a report
type ReportGenerateInput struct {
	Type           string                 `json:"type"`
	Format         string                 `json:"format"`
	OrganizationID string                 `json:"organization_id,omitempty"`
	Filter         map[string]interface{} `json:"filter,omitempty"`
	TimeRange      *TimeRange             `json:"time_range,omitempty"`
}

// TimeRange represents a time range for reports
type TimeRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// ReportListOptions contains options for listing reports
type ReportListOptions struct {
	ListOptions
	Type           string     `url:"type,omitempty"`
	Format         string     `url:"format,omitempty"`
	Status         string     `url:"status,omitempty"`
	OrganizationID string     `url:"organization_id,omitempty"`
	CreatedAfter   *time.Time `url:"created_after,omitempty"`
	CreatedBefore  *time.Time `url:"created_before,omitempty"`
}

// ReportParams contains common parameters for retrieving reports
type ReportParams struct {
	OrganizationID string     `url:"organization_id,omitempty"`
	From           *time.Time `url:"from,omitempty"`
	To             *time.Time `url:"to,omitempty"`
	Format         string     `url:"format,omitempty"`
}

// ReportExportParams contains parameters for exporting a report
type ReportExportParams struct {
	ReportID string `url:"-"` // Not a query parameter, used in URL path
	Format   string `url:"format,omitempty"`
}

// ReportScheduleParams contains parameters for scheduling a report
type ReportScheduleParams struct {
	Type           string     `json:"type"`
	Format         string     `json:"format"`
	Frequency      string     `json:"frequency"`
	Recipients     []string   `json:"recipients"`
	OrganizationID string     `json:"organization_id,omitempty"`
	NextRunAt      *time.Time `json:"next_run_at,omitempty"`
}

// ----- Usage Types -----

// UsageParams contains parameters for retrieving usage statistics
type UsageParams struct {
	OrganizationID string     `url:"organization_id,omitempty"`
	From           *time.Time `url:"from,omitempty"`
	To             *time.Time `url:"to,omitempty"`
}

// UsageReport represents a usage report
type UsageReport struct {
	AccountID      string         `json:"account_id"`
	OrganizationID string         `json:"organization_id,omitempty"`
	From           time.Time      `json:"from"`
	To             time.Time      `json:"to"`
	AgentCount     int            `json:"agent_count"`
	ActiveAgents   int            `json:"active_agents"`
	Usage          map[string]int `json:"usage"`
}
