// Package huntress provides a client for the Huntress API
package huntress

import "time"

// This file contains model types that represent responses from the API.
// Common types shared with parameter definitions have been moved to types.go.
// This file only contains response-specific models.

// ----- Common Types -----

// ID is a common type for entity IDs
type ID string

// ----- Account Types -----

// Account represents a Huntress account
type Account struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Settings    Settings  `json:"settings,omitempty"`
	ContactInfo Contact   `json:"contact_info,omitempty"`
}

// AccountStats represents account statistics
type AccountStats struct {
	OrganizationCount int `json:"organization_count"`
	AgentCount        int `json:"agent_count"`
	IncidentCount     int `json:"incident_count"`
	UserCount         int `json:"user_count"`
}

// Settings represents account or organization settings
type Settings struct {
	NotificationPreferences map[string]interface{} `json:"notification_preferences,omitempty"`
	SecuritySettings        map[string]interface{} `json:"security_settings,omitempty"`
	IntegrationSettings     map[string]interface{} `json:"integration_settings,omitempty"`
}

// ----- User Types -----

// User represents a user in the Huntress system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Role      string    `json:"role,omitempty"`
	Roles     []string  `json:"roles,omitempty"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ----- Organization Types -----

// Organization represents a customer organization
type Organization struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	Tags        []string  `json:"tags,omitempty"`
	Industry    string    `json:"industry,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Settings    Settings  `json:"settings,omitempty"`
	Address     Address   `json:"address,omitempty"`
	ContactInfo Contact   `json:"contact_info,omitempty"`
}

// Address represents a physical address
type Address struct {
	Street1 string `json:"street1,omitempty"`
	Street2 string `json:"street2,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	ZipCode string `json:"zip_code,omitempty"`
	Country string `json:"country,omitempty"`
}

// Contact represents contact information
type Contact struct {
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Title       string `json:"title,omitempty"`
}

// ----- Agent Types -----

// Agent represents a Huntress agent installed on an endpoint
type Agent struct {
	ID             string          `json:"id"`
	Version        string          `json:"version"`
	Hostname       string          `json:"hostname"`
	IPV4Address    string          `json:"ipv4_address"`
	Platform       string          `json:"platform"`
	OS             string          `json:"os"`
	OSVersion      string          `json:"os_version"`
	Status         string          `json:"status"`
	OrganizationID string          `json:"organization_id"`
	LastSeenAt     time.Time       `json:"last_seen_at"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Tags           []string        `json:"tags,omitempty"`
	Stats          AgentStatistics `json:"statistics,omitempty"`
	Settings       AgentSettings   `json:"settings,omitempty"`
	InstalledBy    string          `json:"installed_by,omitempty"`
	MachineUUID    string          `json:"machine_uuid,omitempty"`
	MachineID      string          `json:"machine_id,omitempty"`
}

// ----- Incident Types -----

// Incident represents a security incident
type Incident struct {
	ID             string                 `json:"id"`
	Type           string                 `json:"type"`
	Title          string                 `json:"title"`
	Description    string                 `json:"description"`
	Severity       string                 `json:"severity"`
	Status         string                 `json:"status"`
	OrganizationID string                 `json:"organization_id"`
	AgentID        string                 `json:"agent_id"`
	DetectedAt     time.Time              `json:"detected_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
	ResolvedAt     time.Time              `json:"resolved_at,omitempty"`
	AssignedTo     string                 `json:"assigned_to,omitempty"`
	Tags           []string               `json:"tags,omitempty"`
	Details        map[string]interface{} `json:"details,omitempty"`
}

// ----- Report Types -----

// Report represents a generated report
type Report struct {
	ID             string    `json:"id"`
	Type           string    `json:"type"`
	Status         string    `json:"status"`
	Format         string    `json:"format"`
	OrganizationID string    `json:"organization_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	CompletedAt    time.Time `json:"completed_at,omitempty"`
	URL            string    `json:"url,omitempty"`
}

// DetailedReport represents a detailed report with full content
type DetailedReport struct {
	Report
	Content map[string]interface{} `json:"content"`
}

// SummaryReport represents a summary report
type SummaryReport struct {
	Report
	Summary map[string]interface{} `json:"summary"`
}

// ReportSchedule represents a scheduled report
type ReportSchedule struct {
	ID             string    `json:"id"`
	Type           string    `json:"type"`
	Format         string    `json:"format"`
	Frequency      string    `json:"frequency"`
	Recipients     []string  `json:"recipients"`
	OrganizationID string    `json:"organization_id,omitempty"`
	NextRunAt      time.Time `json:"next_run_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ----- Billing Types -----

// BillingSummary represents billing summary information
type BillingSummary struct {
	AccountID      string  `json:"account_id"`
	CurrentBalance float64 `json:"current_balance"`
	Currency       string  `json:"currency"`
	BillingPeriod  string  `json:"billing_period"`
	DueDate        string  `json:"due_date"`
}

// Invoice represents a billing invoice
type Invoice struct {
	ID            string            `json:"id"`
	InvoiceNumber string            `json:"invoice_number"`
	AccountID     string            `json:"account_id"`
	Amount        float64           `json:"amount"`
	Currency      string            `json:"currency"`
	Status        string            `json:"status"`
	IssuedAt      time.Time         `json:"issued_at"`
	DueAt         time.Time         `json:"due_at"`
	PaidAt        time.Time         `json:"paid_at,omitempty"`
	BillingPeriod string            `json:"billing_period"`
	LineItems     []InvoiceLineItem `json:"line_items,omitempty"`
}

// InvoiceLineItem represents a line item in an invoice
type InvoiceLineItem struct {
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Amount      float64 `json:"amount"`
}
