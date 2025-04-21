package huntress

import (
	"context"
	"time"
)

// AccountService provides access to the account-related API endpoints
type AccountService interface {
	// Get retrieves account details by ID
	Get(ctx context.Context, id int) (*Account, error)

	// GetCurrent retrieves the current authenticated account
	GetCurrent(ctx context.Context) (*Account, error)

	// Update updates account settings
	Update(ctx context.Context, id int, input *AccountUpdateInput) (*Account, error)

	// ListUsers lists users associated with an account
	ListUsers(ctx context.Context, id int, opts *ListOptions) ([]*User, *Pagination, error)

	// GetStatistics retrieves account statistics
	GetStatistics(ctx context.Context, id int) (*AccountStatistics, error)
}

// OrganizationService provides access to the organization-related API endpoints
type OrganizationService interface {
	// Get retrieves organization details by ID
	Get(ctx context.Context, id int) (*Organization, error)

	// Create creates a new organization
	Create(ctx context.Context, input *OrganizationCreateInput) (*Organization, error)

	// Update updates an organization
	Update(ctx context.Context, id int, input *OrganizationUpdateInput) (*Organization, error)

	// Delete deletes an organization
	Delete(ctx context.Context, id int) error

	// List lists organizations with optional filtering
	List(ctx context.Context, opts *OrganizationListOptions) ([]*Organization, *Pagination, error)

	// ListUsers lists users associated with an organization
	ListUsers(ctx context.Context, id int, opts *ListOptions) ([]*User, *Pagination, error)

	// GetStatistics retrieves organization statistics
	GetStatistics(ctx context.Context, id int) (*OrganizationStatistics, error)
}

// AgentService provides access to the agent-related API endpoints
type AgentService interface {
	// Get retrieves agent details by ID
	Get(ctx context.Context, id string) (*Agent, error)

	// List lists agents with optional filtering
	List(ctx context.Context, opts *AgentListOptions) ([]*Agent, *Pagination, error)

	// Update updates an agent
	Update(ctx context.Context, id string, input *AgentUpdateInput) (*Agent, error)

	// Delete deletes an agent
	Delete(ctx context.Context, id string) error

	// UpdateStatus updates the status of an agent
	UpdateStatus(ctx context.Context, id string, status string) error

	// GetStatistics retrieves agent statistics
	GetStatistics(ctx context.Context, id string) (*AgentStatistics, error)
}

// IncidentService provides access to the incident-related API endpoints
type IncidentService interface {
	// Get retrieves incident details by ID
	Get(ctx context.Context, id string) (*Incident, error)

	// List lists incidents with optional filtering
	List(ctx context.Context, opts *IncidentListOptions) ([]*Incident, *Pagination, error)

	// Update updates an incident
	Update(ctx context.Context, id string, input *IncidentUpdateInput) (*Incident, error)

	// UpdateStatus updates the status of an incident
	UpdateStatus(ctx context.Context, id string, status string) error

	// AddNote adds a note to an incident
	AddNote(ctx context.Context, id string, note string) error

	// ListNotes lists notes for an incident
	ListNotes(ctx context.Context, id string, opts *ListOptions) ([]*IncidentNote, *Pagination, error)
}

// ReportService provides access to the report-related API endpoints
type ReportService interface {
	// Generate generates a report
	Generate(ctx context.Context, input *ReportGenerateInput) (*Report, error)

	// Get retrieves report details by ID
	Get(ctx context.Context, id string) (*Report, error)

	// List lists reports with optional filtering
	List(ctx context.Context, opts *ReportListOptions) ([]*Report, *Pagination, error)

	// Download downloads a report
	Download(ctx context.Context, id string, format string) ([]byte, error)
}

// BillingService provides access to the billing-related API endpoints
type BillingService interface {
	// GetInvoice retrieves an invoice by ID
	GetInvoice(ctx context.Context, id string) (*Invoice, error)

	// ListInvoices lists invoices with optional filtering
	ListInvoices(ctx context.Context, opts *InvoiceListOptions) ([]*Invoice, *Pagination, error)

	// GetUsage retrieves usage information
	GetUsage(ctx context.Context, period *BillingPeriod) (*UsageReport, error)
}

// ListOptions contains common options for list operations
type ListOptions struct {
	Page     int    `url:"page,omitempty"`
	PerPage  int    `url:"per_page,omitempty"`
	SortBy   string `url:"sort_by,omitempty"`
	SortDesc bool   `url:"sort_desc,omitempty"`
}

// OrganizationListOptions contains options for listing organizations
type OrganizationListOptions struct {
	ListOptions
	Status  string   `url:"status,omitempty"`
	Search  string   `url:"search,omitempty"`
	Tags    []string `url:"tags,omitempty"`
	Account int      `url:"account_id,omitempty"`
}

// AgentListOptions contains options for listing agents
type AgentListOptions struct {
	ListOptions
	Status         string    `url:"status,omitempty"`
	Platform       string    `url:"platform,omitempty"`
	Search         string    `url:"search,omitempty"`
	Organization   int       `url:"organization_id,omitempty"`
	LastSeenSince  time.Time `url:"last_seen_since,omitempty"`
	LastSeenBefore time.Time `url:"last_seen_before,omitempty"`
}

// IncidentListOptions contains options for listing incidents
type IncidentListOptions struct {
	ListOptions
	Status         string    `url:"status,omitempty"`
	Severity       string    `url:"severity,omitempty"`
	Type           string    `url:"type,omitempty"`
	Organization   int       `url:"organization_id,omitempty"`
	Agent          string    `url:"agent_id,omitempty"`
	Search         string    `url:"search,omitempty"`
	DetectedAfter  time.Time `url:"detected_after,omitempty"`
	DetectedBefore time.Time `url:"detected_before,omitempty"`
	ResolvedAfter  time.Time `url:"resolved_after,omitempty"`
	ResolvedBefore time.Time `url:"resolved_before,omitempty"`
}

// ReportListOptions contains options for listing reports
type ReportListOptions struct {
	ListOptions
	Type            string    `url:"type,omitempty"`
	Organization    int       `url:"organization_id,omitempty"`
	GeneratedAfter  time.Time `url:"generated_after,omitempty"`
	GeneratedBefore time.Time `url:"generated_before,omitempty"`
}

// InvoiceListOptions contains options for listing invoices
type InvoiceListOptions struct {
	ListOptions
	Status       string    `url:"status,omitempty"`
	Organization int       `url:"organization_id,omitempty"`
	IssuedAfter  time.Time `url:"issued_after,omitempty"`
	IssuedBefore time.Time `url:"issued_before,omitempty"`
	PaidAfter    time.Time `url:"paid_after,omitempty"`
	PaidBefore   time.Time `url:"paid_before,omitempty"`
}

// BillingPeriod represents a billing period
type BillingPeriod struct {
	Start time.Time `url:"start"`
	End   time.Time `url:"end"`
}

// Pagination represents pagination information for list responses
type Pagination struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalPages  int `json:"total_pages"`
	TotalItems  int `json:"total_items"`
}
