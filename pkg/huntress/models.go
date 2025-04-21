package huntress

import "time"

// Pagination represents pagination information
type Pagination struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalPages  int `json:"total_pages"`
	TotalItems  int `json:"total_items"`
}

// ListOptions represents common options for list operations
type ListOptions struct {
	Page     int    `url:"page,omitempty"`
	PerPage  int    `url:"per_page,omitempty"`
	SortBy   string `url:"sort_by,omitempty"`
	SortDesc bool   `url:"sort_desc,omitempty"`
}

// Account represents a Huntress account in the public API
type Account struct {
	ID            int                    `json:"id"`
	Name          string                 `json:"name"`
	Status        string                 `json:"status"`
	ContactName   string                 `json:"contact_name"`
	ContactEmail  string                 `json:"contact_email"`
	ContactPhone  string                 `json:"contact_phone,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	BillingMethod string                 `json:"billing_method,omitempty"`
	Features      []string               `json:"features,omitempty"`
	Settings      map[string]interface{} `json:"settings,omitempty"`
}

// AccountUpdateInput contains the fields that can be updated on an account
type AccountUpdateInput struct {
	Name         *string                `json:"name,omitempty"`
	ContactName  *string                `json:"contact_name,omitempty"`
	ContactEmail *string                `json:"contact_email,omitempty"`
	ContactPhone *string                `json:"contact_phone,omitempty"`
	Settings     map[string]interface{} `json:"settings,omitempty"`
}

// AccountStatistics contains statistics for an account
type AccountStatistics struct {
	OrganizationCount  int `json:"organization_count"`
	ActiveAgentCount   int `json:"active_agent_count"`
	TotalAgentCount    int `json:"total_agent_count"`
	OpenIncidentCount  int `json:"open_incident_count"`
	TotalIncidentCount int `json:"total_incident_count"`
}

// Organization represents a Huntress organization in the public API
type Organization struct {
	ID          int                    `json:"id"`
	AccountID   int                    `json:"account_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Status      string                 `json:"status"`
	Address     Address                `json:"address,omitempty"`
	ContactInfo ContactInfo            `json:"contact_info,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Industry    string                 `json:"industry,omitempty"`
	AgentCount  int                    `json:"agent_count,omitempty"`
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

// ContactInfo represents contact information
type ContactInfo struct {
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Title       string `json:"title,omitempty"`
}

// OrganizationCreateInput contains the fields needed to create an organization
type OrganizationCreateInput struct {
	AccountID   int                    `json:"account_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Status      string                 `json:"status,omitempty"`
	Address     *Address               `json:"address,omitempty"`
	ContactInfo *ContactInfo           `json:"contact_info,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Industry    string                 `json:"industry,omitempty"`
}

// OrganizationUpdateInput contains the fields that can be updated on an organization
type OrganizationUpdateInput struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	Status      *string                `json:"status,omitempty"`
	Address     *Address               `json:"address,omitempty"`
	ContactInfo *ContactInfo           `json:"contact_info,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Industry    *string                `json:"industry,omitempty"`
}

// OrganizationStatistics contains statistics for an organization
type OrganizationStatistics struct {
	ActiveAgentCount   int `json:"active_agent_count"`
	TotalAgentCount    int `json:"total_agent_count"`
	OpenIncidentCount  int `json:"open_incident_count"`
	TotalIncidentCount int `json:"total_incident_count"`
}

// Agent represents a Huntress agent in the public API
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

// Incident represents a security incident in the public API
type Incident struct {
	ID             string     `json:"id"`
	OrganizationID int        `json:"organization_id"`
	AgentID        string     `json:"agent_id,omitempty"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Status         string     `json:"status"`
	Severity       string     `json:"severity"`
	Type           string     `json:"type"`
	DetectedAt     time.Time  `json:"detected_at"`
	ResolvedAt     *time.Time `json:"resolved_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	AssignedTo     string     `json:"assigned_to,omitempty"`
	Tags           []string   `json:"tags,omitempty"`
}

// IncidentUpdateInput contains the fields that can be updated on an incident
type IncidentUpdateInput struct {
	Status     *string  `json:"status,omitempty"`
	AssignedTo *string  `json:"assigned_to,omitempty"`
	Tags       []string `json:"tags,omitempty"`
}

// IncidentNote represents a note on an incident
type IncidentNote struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Report represents a generated report in the public API
type Report struct {
	ID             string    `json:"id"`
	OrganizationID int       `json:"organization_id,omitempty"`
	AccountID      int       `json:"account_id,omitempty"`
	Type           string    `json:"type"`
	Name           string    `json:"name"`
	Description    string    `json:"description,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	GeneratedAt    time.Time `json:"generated_at"`
	Status         string    `json:"status"`
	Format         string    `json:"format,omitempty"` // pdf, csv, etc.
	URL            string    `json:"url,omitempty"`
	Size           int64     `json:"size,omitempty"`
}

// ReportGenerateInput contains the fields needed to generate a report
type ReportGenerateInput struct {
	Type           string    `json:"type"`
	OrganizationID *int      `json:"organization_id,omitempty"`
	AccountID      *int      `json:"account_id,omitempty"`
	StartDate      time.Time `json:"start_date,omitempty"`
	EndDate        time.Time `json:"end_date,omitempty"`
	Format         string    `json:"format,omitempty"` // pdf, csv, etc.
	Name           string    `json:"name,omitempty"`
	Description    string    `json:"description,omitempty"`
}

// Invoice represents a billing invoice in the public API
type Invoice struct {
	ID             string        `json:"id"`
	AccountID      int           `json:"account_id"`
	OrganizationID int           `json:"organization_id,omitempty"`
	Status         string        `json:"status"`
	Amount         float64       `json:"amount"`
	Currency       string        `json:"currency"`
	IssuedDate     time.Time     `json:"issued_date"`
	DueDate        time.Time     `json:"due_date"`
	PaidDate       *time.Time    `json:"paid_date,omitempty"`
	Items          []InvoiceItem `json:"items"`
	PDF            string        `json:"pdf_url,omitempty"`
}

// InvoiceItem represents an item on an invoice
type InvoiceItem struct {
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Amount      float64 `json:"amount"`
}

// UsageReport represents usage information for billing
type UsageReport struct {
	AccountID         int                    `json:"account_id"`
	PeriodStart       time.Time              `json:"period_start"`
	PeriodEnd         time.Time              `json:"period_end"`
	AgentCount        int                    `json:"agent_count"`
	AgentHours        int64                  `json:"agent_hours"`
	OrganizationCount int                    `json:"organization_count"`
	Details           map[string]interface{} `json:"details,omitempty"`
}

// User represents a Huntress user in the public API
type User struct {
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Role        string     `json:"role"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}
