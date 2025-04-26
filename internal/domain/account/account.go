// Package account contains the domain model for Huntress accounts.
package account

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Common errors related to account domain
var (
	ErrInvalidID          = errors.New("invalid account ID")
	ErrEmptyName          = errors.New("account name cannot be empty")
	ErrInvalidEmail       = errors.New("invalid email address")
	ErrInvalidPhoneNumber = errors.New("invalid phone number format")
	ErrEmptyContactName   = errors.New("contact name cannot be empty")
	ErrEmptyContactEmail  = errors.New("contact email cannot be empty")
)

const (
	// StatusTrialing indicates an account in trial period
	StatusTrialing Status = "trialing"
)

// Account represents a Huntress account entity
type Account struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Status          Status    `json:"status"`
	PrimaryContact  Contact   `json:"primaryContact"`
	BillingContact  Contact   `json:"billingContact"`
	PreferredCDT    string    `json:"preferredCDT"`
	Timezone        string    `json:"timezone"`
	Created         time.Time `json:"created"`
	Modified        time.Time `json:"modified"`
	WebhookURL      string    `json:"webhookUrl,omitempty"`
	WebhookUsername string    `json:"webhookUsername,omitempty"`
	WebhookPassword string    `json:"webhookPassword,omitempty"`
}

// Contact represents contact information for an account
type Contact struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone,omitempty"`
	Address string `json:"address,omitempty"`
}

// Validate ensures the account data is valid
func (a *Account) Validate() error {
	if a.Name == "" {
		return ErrEmptyName
	}

	if err := a.PrimaryContact.Validate(); err != nil {
		return err
	}

	if err := a.BillingContact.Validate(); err != nil {
		return err
	}

	return nil
}

// Validate ensures the contact information is valid
func (c *Contact) Validate() error {
	if c.Name == "" {
		return ErrEmptyContactName
	}

	if c.Email == "" {
		return ErrEmptyContactEmail
	}

	return nil
}

// IsActive returns true if the account is in active status
func (a *Account) IsActive() bool {
	return a.Status == StatusActive
}
