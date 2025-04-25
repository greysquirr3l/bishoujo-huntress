// Package organization contains domain entities related to organizations in the system.
package organization

import (
	"fmt"
	"strings"
	"time"
)

// Organization represents a customer organization within a Huntress account.
type Organization struct {
	ID          string                 // Unique identifier
	AccountID   int                    // Parent account identifier
	Name        string                 // Organization name
	Description string                 // Optional description
	Status      string                 // Organization status (active, inactive, etc.)
	CreatedAt   time.Time              // Creation timestamp
	UpdatedAt   time.Time              // Last update timestamp
	Settings    map[string]interface{} // Organization-specific settings
	Tags        []string               // Organization tags for categorization
	Industry    string                 // Industry classification
	AgentCount  int                    // Number of agents deployed
	Address     Address                // Physical address
	ContactInfo ContactInfo            // Primary contact information
}

// Address represents a physical address.
type Address struct {
	Street1 string // Address line 1
	Street2 string // Address line 2 (optional)
	City    string // City
	State   string // State/province
	ZipCode string // ZIP/postal code
	Country string // Country
}

// ContactInfo represents organization contact information.
type ContactInfo struct {
	Name        string // Contact name
	Email       string // Contact email
	PhoneNumber string // Contact phone number
	Title       string // Contact title/position
}

// User represents a user associated with an organization.
type User struct {
	ID        string // Unique identifier
	Email     string // User email address
	FirstName string // First name
	LastName  string // Last name
	Role      string // User role (admin, viewer, etc.)
	Status    string // User status (active, inactive, etc.)
}

// ListParams defines optional parameters for listing organizations.
type ListParams struct {
	Page      int      // Page number for pagination
	Limit     int      // Number of items per page
	AccountID int      // Filter by account ID
	Status    string   // Filter by status
	Search    string   // Search term
	Industry  string   // Filter by industry
	Tags      []string // Filter by tags
}

// Status constants for organizations
const (
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusPending  = "pending"
)

// Role constants for users
const (
	RoleAdmin   = "admin"
	RoleManager = "manager"
	RoleViewer  = "viewer"
)

// Validate performs validation checks on the Organization entity.
func (o *Organization) Validate() error {
	if o.Name == "" {
		return fmt.Errorf("organization name is required")
	}

	if o.AccountID <= 0 {
		return fmt.Errorf("invalid account ID: must be greater than 0")
	}

	if o.Status != "" && !isValidStatus(o.Status) {
		return fmt.Errorf("invalid status: %s", o.Status)
	}

	// If contact info is provided, validate it
	if o.ContactInfo.Email != "" {
		if !isValidEmail(o.ContactInfo.Email) {
			return fmt.Errorf("invalid contact email format")
		}
	}

	return nil
}

// IsActive returns true if the organization is active.
func (o *Organization) IsActive() bool {
	return o.Status == StatusActive
}

// CanHaveAgents returns true if the organization can have agents associated with it.
func (o *Organization) CanHaveAgents() bool {
	return o.IsActive()
}

// FullName returns the user's full name.
func (u *User) FullName() string {
	return strings.TrimSpace(fmt.Sprintf("%s %s", u.FirstName, u.LastName))
}

// Validate performs validation checks on the User entity.
func (u *User) Validate() error {
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}

	if !isValidEmail(u.Email) {
		return fmt.Errorf("invalid email format")
	}

	if u.Role != "" && !isValidRole(u.Role) {
		return fmt.Errorf("invalid role: %s", u.Role)
	}

	return nil
}

// HasAdminAccess returns true if the user has admin access.
func (u *User) HasAdminAccess() bool {
	return u.Role == RoleAdmin
}

// HasManagerAccess returns true if the user has manager or admin access.
func (u *User) HasManagerAccess() bool {
	return u.Role == RoleAdmin || u.Role == RoleManager
}

// Helper functions

// isValidStatus checks if the provided status is valid.
func isValidStatus(status string) bool {
	validStatuses := []string{StatusActive, StatusInactive, StatusPending}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

// isValidRole checks if the provided role is valid.
func isValidRole(role string) bool {
	validRoles := []string{RoleAdmin, RoleManager, RoleViewer}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

// isValidEmail performs a basic check if the email format is valid.
func isValidEmail(email string) bool {
	// Simple email validation: contains @ and at least one dot after @
	atIndex := strings.Index(email, "@")
	if atIndex < 1 || atIndex == len(email)-1 {
		return false
	}

	domain := email[atIndex+1:]
	return strings.Contains(domain, ".")
}
