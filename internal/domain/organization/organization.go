package organization

import (
	"errors"
	"time"
)

// Common errors for organization domain
var (
	ErrInvalidID      = errors.New("invalid organization ID")
	ErrEmptyName      = errors.New("organization name cannot be empty")
	ErrInvalidAccount = errors.New("invalid account ID")
)

// OrganizationStatus represents the current status of an organization
type OrganizationStatus string

const (
	// OrganizationStatusActive indicates an active organization
	OrganizationStatusActive OrganizationStatus = "active"
	// OrganizationStatusInactive indicates an inactive organization
	OrganizationStatusInactive OrganizationStatus = "inactive"
	// OrganizationStatusPending indicates a pending organization
	OrganizationStatusPending OrganizationStatus = "pending"
)

// Organization represents a customer organization within a Huntress account
type Organization struct {
	ID          int
	AccountID   int
	Name        string
	Description string
	Status      OrganizationStatus
	Address     Address
	ContactInfo ContactInfo
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Settings    map[string]interface{}
	Tags        []string
	Industry    string
	AgentCount  int
}

// Address represents an organization's physical address
type Address struct {
	Street1 string
	Street2 string
	City    string
	State   string
	ZipCode string
	Country string
}

// ContactInfo represents contact information for an organization
type ContactInfo struct {
	Name        string
	Email       string
	PhoneNumber string
	Title       string
}

// Validate checks if the organization has valid data
func (o *Organization) Validate() error {
	if o.ID <= 0 {
		return ErrInvalidID
	}
	if o.AccountID <= 0 {
		return ErrInvalidAccount
	}
	if o.Name == "" {
		return ErrEmptyName
	}
	return nil
}

// IsActive returns whether the organization is active
func (o *Organization) IsActive() bool {
	return o.Status == OrganizationStatusActive
}

// HasTag checks if the organization has a specific tag
func (o *Organization) HasTag(tag string) bool {
	for _, t := range o.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// AddTag adds a tag to the organization if it doesn't already exist
func (o *Organization) AddTag(tag string) {
	if tag == "" || o.HasTag(tag) {
		return
	}
	o.Tags = append(o.Tags, tag)
}

// RemoveTag removes a tag from the organization
func (o *Organization) RemoveTag(tag string) {
	if tag == "" {
		return
	}

	for i, t := range o.Tags {
		if t == tag {
			// Remove the tag by replacing it with the last element
			// and then truncating the slice
			o.Tags[i] = o.Tags[len(o.Tags)-1]
			o.Tags = o.Tags[:len(o.Tags)-1]
			return
		}
	}
}
