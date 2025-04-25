// Package huntress provides a client for the Huntress API
package huntress

import (
	"context"
)

// AccountService handles Huntress account operations
type AccountService interface {
	// Get retrieves the current account information
	Get(ctx context.Context) (*Account, error)

	// Update updates account settings
	Update(ctx context.Context, account *AccountUpdateParams) (*Account, error)

	// GetStats retrieves account statistics
	GetStats(ctx context.Context) (*AccountStats, error)

	// ListUsers retrieves all users in the account
	ListUsers(ctx context.Context, params *ListParams) ([]*User, *Pagination, error)

	// AddUser adds a new user to the account
	AddUser(ctx context.Context, user *UserCreateParams) (*User, error)

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, userID string, user *UserUpdateParams) (*User, error)

	// RemoveUser removes a user from the account
	RemoveUser(ctx context.Context, userID string) error
}

// OrganizationService handles Huntress organization operations
type OrganizationService interface {
	// Get retrieves a specific organization
	Get(ctx context.Context, id string) (*Organization, error)

	// List returns all organizations with optional filtering
	List(ctx context.Context, params *ListOrganizationsParams) ([]*Organization, *Pagination, error)

	// Create creates a new organization
	Create(ctx context.Context, org *OrganizationCreateParams) (*Organization, error)

	// Update updates an existing organization
	Update(ctx context.Context, id string, org *OrganizationUpdateParams) (*Organization, error)

	// Delete removes an organization
	Delete(ctx context.Context, id string) error

	// ListUsers retrieves all users in an organization
	ListUsers(ctx context.Context, orgID string, params *ListParams) ([]*User, *Pagination, error)

	// AddUser adds a user to an organization
	AddUser(ctx context.Context, orgID string, user *UserCreateParams) (*User, error)

	// RemoveUser removes a user from an organization
	RemoveUser(ctx context.Context, orgID string, userID string) error
}

// AgentService handles Huntress agent operations
type AgentService interface {
	// Get retrieves a specific agent
	Get(ctx context.Context, id string) (*Agent, error)

	// List returns all agents with optional filtering
	List(ctx context.Context, params *AgentListOptions) ([]*Agent, *Pagination, error)

	// GetStats retrieves statistics for a specific agent
	GetStats(ctx context.Context, id string) (*AgentStatistics, error)

	// Update updates an existing agent
	Update(ctx context.Context, id string, agent map[string]interface{}) (*Agent, error)

	// Delete removes an agent
	Delete(ctx context.Context, id string) error
}

// IncidentService handles Huntress incident operations
type IncidentService interface {
	// Get retrieves a specific incident
	Get(ctx context.Context, id string) (*Incident, error)

	// List returns all incidents with optional filtering
	List(ctx context.Context, params *IncidentListOptions) ([]*Incident, *Pagination, error)

	// UpdateStatus updates the status of an incident
	UpdateStatus(ctx context.Context, id string, status string) (*Incident, error)

	// Assign assigns an incident to a user
	Assign(ctx context.Context, id string, userID string) (*Incident, error)
}

// ReportService handles Huntress report operations
type ReportService interface {
	// Generate generates a report
	Generate(ctx context.Context, input *ReportGenerateInput) (*Report, error)

	// Get retrieves a specific report
	Get(ctx context.Context, id string) (*Report, error)

	// List returns all reports with optional filtering
	List(ctx context.Context, opts *ReportListOptions) ([]*Report, *Pagination, error)

	// Download downloads a report
	Download(ctx context.Context, id string, format string) ([]byte, error)

	// GetSummary retrieves a summary report
	GetSummary(ctx context.Context, params *ReportParams) (*SummaryReport, error)

	// GetDetails retrieves a detailed report
	GetDetails(ctx context.Context, params *ReportParams) (*DetailedReport, error)

	// Export exports a report in the specified format
	Export(ctx context.Context, params *ReportExportParams) ([]byte, error)

	// Schedule schedules a report for delivery
	Schedule(ctx context.Context, params *ReportScheduleParams) (*ReportSchedule, error)
}

// BillingService handles Huntress billing operations
type BillingService interface {
	// GetSummary retrieves a billing summary
	GetSummary(ctx context.Context) (*BillingSummary, error)

	// ListInvoices lists all invoices
	ListInvoices(ctx context.Context, params *ListParams) ([]*Invoice, *Pagination, error)

	// GetInvoice retrieves a specific invoice
	GetInvoice(ctx context.Context, id string) (*Invoice, error)

	// GetUsage retrieves usage statistics
	GetUsage(ctx context.Context, params *UsageParams) (*UsageReport, error)
}
