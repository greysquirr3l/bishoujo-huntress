// Package organization contains command handlers for organization operations
package organization

import (
	"context"
	"fmt"
	"time"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/organization"
	"github.com/greysquirr3l/bishoujo-huntress/internal/ports/repository"
)

// UpdateOrganizationCommand represents a request to update an organization
type UpdateOrganizationCommand struct {
	ID          string                 // Organization ID (required)
	Name        string                 // Organization name
	Description string                 // Optional description
	Status      string                 // Organization status
	Industry    string                 // Industry classification
	Tags        []string               // Organization tags
	Settings    map[string]interface{} // Organization settings
	Address     struct {
		Street1 string // Address line 1
		Street2 string // Address line 2 (optional)
		City    string // City
		State   string // State/province
		ZipCode string // ZIP/postal code
		Country string // Country
	}
	ContactInfo struct {
		Name        string // Contact name
		Email       string // Contact email
		PhoneNumber string // Contact phone number
		Title       string // Contact title/position
	}
}

// UpdateOrganizationHandler handles updates to existing organizations
type UpdateOrganizationHandler struct {
	orgRepo repository.OrganizationRepository
}

// NewUpdateOrganizationHandler creates a new organization update handler
func NewUpdateOrganizationHandler(orgRepo repository.OrganizationRepository) *UpdateOrganizationHandler {
	return &UpdateOrganizationHandler{
		orgRepo: orgRepo,
	}
}

// updateAddressFields updates the address fields of the organization if provided in the command.
func updateAddressFields(addr *organization.Address, cmdAddr struct {
	Street1 string
	Street2 string
	City    string
	State   string
	ZipCode string
	Country string
}) {
	if cmdAddr.Street1 != "" {
		addr.Street1 = cmdAddr.Street1
	}
	if cmdAddr.Street2 != "" {
		addr.Street2 = cmdAddr.Street2
	}
	if cmdAddr.City != "" {
		addr.City = cmdAddr.City
	}
	if cmdAddr.State != "" {
		addr.State = cmdAddr.State
	}
	if cmdAddr.ZipCode != "" {
		addr.ZipCode = cmdAddr.ZipCode
	}
	if cmdAddr.Country != "" {
		addr.Country = cmdAddr.Country
	}
}

// updateContactInfoFields updates the contact info fields of the organization if provided in the command.
func updateContactInfoFields(info *organization.ContactInfo, cmdInfo struct {
	Name        string
	Email       string
	PhoneNumber string
	Title       string
}) {
	if cmdInfo.Name != "" {
		info.Name = cmdInfo.Name
	}
	if cmdInfo.Email != "" {
		info.Email = cmdInfo.Email
	}
	if cmdInfo.PhoneNumber != "" {
		info.PhoneNumber = cmdInfo.PhoneNumber
	}
	if cmdInfo.Title != "" {
		info.Title = cmdInfo.Title
	}
}

// mergeSettings merges the provided settings into the organization's settings.
func mergeSettings(existing *map[string]interface{}, updates map[string]interface{}) {
	if *existing == nil {
		*existing = make(map[string]interface{})
	}
	for k, v := range updates {
		(*existing)[k] = v
	}
}

// Handle processes the update organization command
func (h *UpdateOrganizationHandler) Handle(ctx context.Context, cmd UpdateOrganizationCommand) (*organization.Organization, error) {
	// Verify that organization ID is provided
	if cmd.ID == "" {
		return nil, fmt.Errorf("organization ID is required for update")
	}

	// First retrieve the existing organization
	existingOrg, err := h.orgRepo.Get(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve organization: %w", err)
	}

	// Update the organization with new values from command
	// Only update fields that are provided in the command
	if cmd.Name != "" {
		existingOrg.Name = cmd.Name
	}

	if cmd.Description != "" {
		existingOrg.Description = cmd.Description
	}

	if cmd.Status != "" {
		existingOrg.Status = cmd.Status
	}

	if cmd.Industry != "" {
		existingOrg.Industry = cmd.Industry
	}

	if len(cmd.Tags) > 0 {
		existingOrg.Tags = cmd.Tags
	}

	if cmd.Settings != nil {
		mergeSettings(&existingOrg.Settings, cmd.Settings)
	}

	if cmd.Address.Street1 != "" || cmd.Address.Street2 != "" || cmd.Address.City != "" || cmd.Address.State != "" || cmd.Address.ZipCode != "" || cmd.Address.Country != "" {
		updateAddressFields(&existingOrg.Address, cmd.Address)
	}

	if cmd.ContactInfo.Email != "" || cmd.ContactInfo.PhoneNumber != "" || cmd.ContactInfo.Title != "" {
		updateContactInfoFields(&existingOrg.ContactInfo, cmd.ContactInfo)
	}

	// Update the timestamp
	existingOrg.UpdatedAt = time.Now()

	// Validate the updated organization
	if err := existingOrg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid organization update: %w", err)
	}

	// Call repository to update organization
	updatedOrg, err := h.orgRepo.Update(ctx, existingOrg)
	if err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}

	return updatedOrg, nil
}
