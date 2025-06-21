# Data Schemas

## Organization

### Input Schema (ListParams)

```go
type ListParams struct {
    Page     int    `json:"page,omitempty"`     // Page number (1-based)
    Limit    int    `json:"limit,omitempty"`    // Items per page (1-100)
    Search   string `json:"search,omitempty"`   // Search term
    SortBy   string `json:"sort_by,omitempty"`  // Sort field
    SortDir  string `json:"sort_dir,omitempty"` // Sort direction (asc/desc)
}
```

### Output Schema (Organization)

```go
type Organization struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Settings    *Settings `json:"settings,omitempty"`
    Status      string    `json:"status"`              // active, inactive, suspended
    Type        string    `json:"type"`                // customer, partner, internal
    ParentID    *int      `json:"parent_id,omitempty"` // For hierarchical organizations
}

type Settings struct {
    Timezone             string `json:"timezone"`
    DefaultLanguage      string `json:"default_language"`
    AutoArchiveIncidents bool   `json:"auto_archive_incidents"`
    NotificationSettings *NotificationSettings `json:"notification_settings,omitempty"`
}

type NotificationSettings struct {
    EmailEnabled bool     `json:"email_enabled"`
    WebhookURL   string   `json:"webhook_url,omitempty"`
    Channels     []string `json:"channels,omitempty"` // email, webhook, slack
}
```

## Account

### Input Schema (ListParams)

```go
type ListParams struct {
    Page         int    `json:"page,omitempty"`
    Limit        int    `json:"limit,omitempty"`
    Search       string `json:"search,omitempty"`
    Status       string `json:"status,omitempty"`       // active, inactive, suspended
    Role         string `json:"role,omitempty"`         // admin, user, readonly
    Organization int    `json:"organization,omitempty"` // Filter by organization ID
    SortBy       string `json:"sort_by,omitempty"`
    SortDir      string `json:"sort_dir,omitempty"`
}

type CreateParams struct {
    Email           string   `json:"email" validate:"required,email"`
    FirstName       string   `json:"first_name" validate:"required"`
    LastName        string   `json:"last_name" validate:"required"`
    Role            string   `json:"role" validate:"required,oneof=admin user readonly"`
    OrganizationIDs []int    `json:"organization_ids" validate:"required,min=1"`
    Permissions     []string `json:"permissions,omitempty"`
    SendInvite      bool     `json:"send_invite"`
}

type UpdateParams struct {
    FirstName       *string  `json:"first_name,omitempty"`
    LastName        *string  `json:"last_name,omitempty"`
    Role            *string  `json:"role,omitempty" validate:"omitempty,oneof=admin user readonly"`
    OrganizationIDs []int    `json:"organization_ids,omitempty"`
    Permissions     []string `json:"permissions,omitempty"`
    Status          *string  `json:"status,omitempty" validate:"omitempty,oneof=active inactive suspended"`
}
```

### Output Schema (Account)

```go
type Account struct {
    ID              int                `json:"id"`
    Email           string             `json:"email"`
    FirstName       string             `json:"first_name"`
    LastName        string             `json:"last_name"`
    Role            string             `json:"role"`
    Status          string             `json:"status"`
    Organizations   []OrganizationRef  `json:"organizations"`
    Permissions     []string           `json:"permissions"`
    LastLoginAt     *time.Time         `json:"last_login_at,omitempty"`
    CreatedAt       time.Time          `json:"created_at"`
    UpdatedAt       time.Time          `json:"updated_at"`
    ProfileSettings *ProfileSettings   `json:"profile_settings,omitempty"`
}

type OrganizationRef struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Role string `json:"role"` // Role within this organization
}

type ProfileSettings struct {
    Timezone         string `json:"timezone"`
    Language         string `json:"language"`
    EmailNotifications bool `json:"email_notifications"`
    TwoFactorEnabled   bool `json:"two_factor_enabled"`
}
```

## Agent

### Input Schema (ListParams)

```go
type ListParams struct {
    Page         int       `json:"page,omitempty"`
    Limit        int       `json:"limit,omitempty"`
    Search       string    `json:"search,omitempty"`
    Status       string    `json:"status,omitempty"`       // online, offline, error
    Platform     string    `json:"platform,omitempty"`     // windows, linux, macos
    Organization int       `json:"organization,omitempty"`
    LastSeenBefore *time.Time `json:"last_seen_before,omitempty"`
    LastSeenAfter  *time.Time `json:"last_seen_after,omitempty"`
    SortBy       string    `json:"sort_by,omitempty"`
    SortDir      string    `json:"sort_dir,omitempty"`
}

type UpdateParams struct {
    Name        *string           `json:"name,omitempty"`
    Description *string           `json:"description,omitempty"`
    Tags        []string          `json:"tags,omitempty"`
    Settings    *AgentSettings    `json:"settings,omitempty"`
}
```

### Output Schema (Agent)

```go
type Agent struct {
    ID             int            `json:"id"`
    Name           string         `json:"name"`
    Hostname       string         `json:"hostname"`
    IPAddress      string         `json:"ip_address"`
    MACAddress     string         `json:"mac_address"`
    Platform       string         `json:"platform"`
    Architecture   string         `json:"architecture"`
    Version        string         `json:"version"`
    Status         string         `json:"status"`
    OrganizationID int            `json:"organization_id"`
    LastSeenAt     *time.Time     `json:"last_seen_at,omitempty"`
    InstalledAt    time.Time      `json:"installed_at"`
    UpdatedAt      time.Time      `json:"updated_at"`
    Description    string         `json:"description,omitempty"`
    Tags           []string       `json:"tags,omitempty"`
    Settings       *AgentSettings `json:"settings,omitempty"`
    SystemInfo     *SystemInfo    `json:"system_info,omitempty"`
}

type AgentSettings struct {
    ScanInterval     int  `json:"scan_interval"`      // seconds
    AutoUpdate       bool `json:"auto_update"`
    RealtimeProtection bool `json:"realtime_protection"`
    LogLevel         string `json:"log_level"`        // debug, info, warn, error
}

type SystemInfo struct {
    OS              string `json:"os"`
    OSVersion       string `json:"os_version"`
    TotalMemory     int64  `json:"total_memory"`      // bytes
    AvailableMemory int64  `json:"available_memory"`  // bytes
    CPUCores        int    `json:"cpu_cores"`
    DiskSpace       int64  `json:"disk_space"`        // bytes
    Uptime          int64  `json:"uptime"`            // seconds
}
```

## Incident

### Input Schema (ListParams)

```go
type ListParams struct {
    Page         int        `json:"page,omitempty"`
    Limit        int        `json:"limit,omitempty"`
    Search       string     `json:"search,omitempty"`
    Status       string     `json:"status,omitempty"`     // open, investigating, resolved, closed
    Severity     string     `json:"severity,omitempty"`   // low, medium, high, critical
    Category     string     `json:"category,omitempty"`   // malware, suspicious, policy_violation
    Organization int        `json:"organization,omitempty"`
    AssignedTo   int        `json:"assigned_to,omitempty"`
    CreatedAfter *time.Time `json:"created_after,omitempty"`
    CreatedBefore *time.Time `json:"created_before,omitempty"`
    SortBy       string     `json:"sort_by,omitempty"`
    SortDir      string     `json:"sort_dir,omitempty"`
}

type UpdateParams struct {
    Status      *string `json:"status,omitempty" validate:"omitempty,oneof=open investigating resolved closed"`
    AssignedTo  *int    `json:"assigned_to,omitempty"`
    Notes       *string `json:"notes,omitempty"`
    Resolution  *string `json:"resolution,omitempty"`
    Tags        []string `json:"tags,omitempty"`
}

type CreateNoteParams struct {
    Content string `json:"content" validate:"required"`
    Type    string `json:"type" validate:"required,oneof=note action resolution"`
}
```

### Output Schema (Incident)

```go
type Incident struct {
    ID             int                `json:"id"`
    Title          string             `json:"title"`
    Description    string             `json:"description"`
    Status         string             `json:"status"`
    Severity       string             `json:"severity"`
    Category       string             `json:"category"`
    OrganizationID int                `json:"organization_id"`
    AgentID        *int               `json:"agent_id,omitempty"`
    AssignedTo     *int               `json:"assigned_to,omitempty"`
    CreatedBy      int                `json:"created_by"`
    CreatedAt      time.Time          `json:"created_at"`
    UpdatedAt      time.Time          `json:"updated_at"`
    ResolvedAt     *time.Time         `json:"resolved_at,omitempty"`
    Tags           []string           `json:"tags,omitempty"`
    Metadata       map[string]interface{} `json:"metadata,omitempty"`
    Notes          []IncidentNote     `json:"notes,omitempty"`
    Artifacts      []IncidentArtifact `json:"artifacts,omitempty"`
}

type IncidentNote struct {
    ID        int       `json:"id"`
    Content   string    `json:"content"`
    Type      string    `json:"type"`    // note, action, resolution
    CreatedBy int       `json:"created_by"`
    CreatedAt time.Time `json:"created_at"`
}

type IncidentArtifact struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Type        string    `json:"type"`        // file, url, hash
    Value       string    `json:"value"`
    Description string    `json:"description,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
}
```

## Report

### Input Schema (ListParams)

```go
type ListParams struct {
    Page         int        `json:"page,omitempty"`
    Limit        int        `json:"limit,omitempty"`
    Search       string     `json:"search,omitempty"`
    Type         string     `json:"type,omitempty"`         // summary, detailed, custom
    Status       string     `json:"status,omitempty"`       // pending, completed, failed
    Organization int        `json:"organization,omitempty"`
    CreatedAfter *time.Time `json:"created_after,omitempty"`
    CreatedBefore *time.Time `json:"created_before,omitempty"`
    SortBy       string     `json:"sort_by,omitempty"`
    SortDir      string     `json:"sort_dir,omitempty"`
}

type GenerateParams struct {
    Name         string                 `json:"name" validate:"required"`
    Type         string                 `json:"type" validate:"required,oneof=summary detailed custom"`
    Format       string                 `json:"format" validate:"required,oneof=pdf csv json"`
    DateRange    DateRange              `json:"date_range" validate:"required"`
    Filters      map[string]interface{} `json:"filters,omitempty"`
    Recipients   []string               `json:"recipients,omitempty" validate:"dive,email"`
    Schedule     *ReportSchedule        `json:"schedule,omitempty"`
}

type DateRange struct {
    StartDate time.Time `json:"start_date" validate:"required"`
    EndDate   time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
}

type ReportSchedule struct {
    Frequency string `json:"frequency" validate:"required,oneof=daily weekly monthly"`
    DayOfWeek *int   `json:"day_of_week,omitempty" validate:"omitempty,min=0,max=6"` // 0=Sunday
    DayOfMonth *int  `json:"day_of_month,omitempty" validate:"omitempty,min=1,max=31"`
    Enabled   bool   `json:"enabled"`
}
```

### Output Schema (Report)

```go
type Report struct {
    ID             int                    `json:"id"`
    Name           string                 `json:"name"`
    Type           string                 `json:"type"`
    Format         string                 `json:"format"`
    Status         string                 `json:"status"`
    OrganizationID int                    `json:"organization_id"`
    CreatedBy      int                    `json:"created_by"`
    CreatedAt      time.Time              `json:"created_at"`
    CompletedAt    *time.Time             `json:"completed_at,omitempty"`
    DateRange      DateRange              `json:"date_range"`
    Filters        map[string]interface{} `json:"filters,omitempty"`
    Recipients     []string               `json:"recipients,omitempty"`
    Schedule       *ReportSchedule        `json:"schedule,omitempty"`
    DownloadURL    string                 `json:"download_url,omitempty"`
    FileSize       *int64                 `json:"file_size,omitempty"` // bytes
    Error          string                 `json:"error,omitempty"`
}
```

## Billing

### Input Schema (ListParams)

```go
type ListParams struct {
    Page         int        `json:"page,omitempty"`
    Limit        int        `json:"limit,omitempty"`
    Organization int        `json:"organization,omitempty"`
    Status       string     `json:"status,omitempty"`     // pending, paid, overdue, cancelled
    PeriodStart  *time.Time `json:"period_start,omitempty"`
    PeriodEnd    *time.Time `json:"period_end,omitempty"`
    SortBy       string     `json:"sort_by,omitempty"`
    SortDir      string     `json:"sort_dir,omitempty"`
}
```

### Output Schema (Billing)

```go
type Invoice struct {
    ID             int           `json:"id"`
    InvoiceNumber  string        `json:"invoice_number"`
    OrganizationID int           `json:"organization_id"`
    Status         string        `json:"status"`
    Amount         float64       `json:"amount"`         // total amount
    Currency       string        `json:"currency"`       // USD, EUR, etc.
    PeriodStart    time.Time     `json:"period_start"`
    PeriodEnd      time.Time     `json:"period_end"`
    IssuedAt       time.Time     `json:"issued_at"`
    DueAt          time.Time     `json:"due_at"`
    PaidAt         *time.Time    `json:"paid_at,omitempty"`
    LineItems      []LineItem    `json:"line_items"`
    TaxAmount      float64       `json:"tax_amount"`
    DiscountAmount float64       `json:"discount_amount"`
    DownloadURL    string        `json:"download_url,omitempty"`
}

type LineItem struct {
    ID          int     `json:"id"`
    Description string  `json:"description"`
    Quantity    int     `json:"quantity"`
    UnitPrice   float64 `json:"unit_price"`
    Amount      float64 `json:"amount"`
    Type        string  `json:"type"`        // subscription, usage, one_time
}

type Usage struct {
    OrganizationID int                    `json:"organization_id"`
    Period         Period                 `json:"period"`
    Metrics        map[string]UsageMetric `json:"metrics"`
    UpdatedAt      time.Time              `json:"updated_at"`
}

type Period struct {
    Start time.Time `json:"start"`
    End   time.Time `json:"end"`
}

type UsageMetric struct {
    Name        string  `json:"name"`
    Value       float64 `json:"value"`
    Unit        string  `json:"unit"`        // agents, scans, gb
    Limit       *float64 `json:"limit,omitempty"`
    Percentage  *float64 `json:"percentage,omitempty"` // if limit exists
}
```

## Webhook

### Input Schema (ListParams)

```go
type ListParams struct {
    Page         int    `json:"page,omitempty"`
    Limit        int    `json:"limit,omitempty"`
    Search       string `json:"search,omitempty"`
    Status       string `json:"status,omitempty"`       // active, inactive
    EventType    string `json:"event_type,omitempty"`   // incident.created, agent.status_changed
    Organization int    `json:"organization,omitempty"`
    SortBy       string `json:"sort_by,omitempty"`
    SortDir      string `json:"sort_dir,omitempty"`
}

type CreateParams struct {
    Name         string   `json:"name" validate:"required"`
    URL          string   `json:"url" validate:"required,url"`
    EventTypes   []string `json:"event_types" validate:"required,min=1"`
    Secret       string   `json:"secret,omitempty"`
    Active       bool     `json:"active"`
    Headers      map[string]string `json:"headers,omitempty"`
    RetryPolicy  *RetryPolicy `json:"retry_policy,omitempty"`
}

type UpdateParams struct {
    Name        *string           `json:"name,omitempty"`
    URL         *string           `json:"url,omitempty" validate:"omitempty,url"`
    EventTypes  []string          `json:"event_types,omitempty"`
    Secret      *string           `json:"secret,omitempty"`
    Active      *bool             `json:"active,omitempty"`
    Headers     map[string]string `json:"headers,omitempty"`
    RetryPolicy *RetryPolicy      `json:"retry_policy,omitempty"`
}

type RetryPolicy struct {
    MaxRetries    int           `json:"max_retries" validate:"min=0,max=10"`
    RetryInterval time.Duration `json:"retry_interval" validate:"min=1s,max=1h"`
    BackoffFactor float64       `json:"backoff_factor" validate:"min=1,max=10"`
}
```

### Output Schema (Webhook)

```go
type Webhook struct {
    ID             int                   `json:"id"`
    Name           string                `json:"name"`
    URL            string                `json:"url"`
    EventTypes     []string              `json:"event_types"`
    Active         bool                  `json:"active"`
    OrganizationID int                   `json:"organization_id"`
    CreatedBy      int                   `json:"created_by"`
    CreatedAt      time.Time             `json:"created_at"`
    UpdatedAt      time.Time             `json:"updated_at"`
    LastTriggered  *time.Time            `json:"last_triggered,omitempty"`
    Headers        map[string]string     `json:"headers,omitempty"`
    RetryPolicy    *RetryPolicy          `json:"retry_policy,omitempty"`
    Statistics     *WebhookStatistics    `json:"statistics,omitempty"`
}

type WebhookStatistics struct {
    TotalDeliveries    int     `json:"total_deliveries"`
    SuccessfulDeliveries int   `json:"successful_deliveries"`
    FailedDeliveries   int     `json:"failed_deliveries"`
    SuccessRate        float64 `json:"success_rate"`       // percentage
    AverageResponseTime time.Duration `json:"average_response_time"`
}

type WebhookDelivery struct {
    ID          int       `json:"id"`
    WebhookID   int       `json:"webhook_id"`
    EventType   string    `json:"event_type"`
    Payload     string    `json:"payload"`     // JSON string
    Status      string    `json:"status"`      // pending, delivered, failed
    StatusCode  *int      `json:"status_code,omitempty"`
    Response    string    `json:"response,omitempty"`
    DeliveredAt *time.Time `json:"delivered_at,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    Attempts    int       `json:"attempts"`
    NextRetry   *time.Time `json:"next_retry,omitempty"`
}
```

## Common Types

### Pagination

```go
type Pagination struct {
    Page       int  `json:"page"`
    Limit      int  `json:"limit"`
    TotalPages int  `json:"total_pages"`
    TotalItems int  `json:"total_items"`
    HasNext    bool `json:"has_next"`
    HasPrev    bool `json:"has_prev"`
}

type PaginatedResponse[T any] struct {
    Data       []T        `json:"data"`
    Pagination Pagination `json:"pagination"`
}
```

### Filters and Search

```go
type DateFilter struct {
    After  *time.Time `json:"after,omitempty"`
    Before *time.Time `json:"before,omitempty"`
}

type NumericFilter struct {
    Min *float64 `json:"min,omitempty"`
    Max *float64 `json:"max,omitempty"`
}

type StringFilter struct {
    Contains []string `json:"contains,omitempty"`
    Excludes []string `json:"excludes,omitempty"`
    Exact    *string  `json:"exact,omitempty"`
}
```

## Validation Tags

All input schemas use struct validation tags:

- `required`: Field is mandatory
- `omitempty`: Field can be omitted if empty
- `email`: Valid email format
- `url`: Valid URL format
- `oneof`: Value must be one of specified options
- `min/max`: Minimum/maximum value or length
- `dive`: Validate each element in slice/map
- `gtfield`: Greater than another field's value

Example usage with the validator package:

```go
import "github.com/go-playground/validator/v10"

validate := validator.New()
err := validate.Struct(params)
if err != nil {
    // Handle validation errors
}
```
