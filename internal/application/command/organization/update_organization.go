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
		// Merge settings instead of replacing
		if existingOrg.Settings == nil {
			existingOrg.Settings = make(map[string]interface{})
		}
		for k, v := range cmd.Settings {
			existingOrg.Settings[k] = v
		}
	}

	// Update address if any address field is provided
	if cmd.Address.Street1 != "" || cmd.Address.Street2 != "" || cmd.Address.City != "" ||
		cmd.Address.State != "" || cmd.Address.ZipCode != "" || cmd.Address.Country != "" {

		if cmd.Address.Street1 != "" {
			existingOrg.Address.Street1 = cmd.Address.Street1
		}
		if cmd.Address.Street2 != "" {
			existingOrg.Address.Street2 = cmd.Address.Street2
		}
		if cmd.Address.City != "" {
			existingOrg.Address.City = cmd.Address.City
		}
		if cmd.Address.State != "" {
			existingOrg.Address.State = cmd.Address.State
		}
		if cmd.Address.ZipCode != "" {
			existingOrg.Address.ZipCode = cmd.Address.ZipCode
		}
		if cmd.Address.Country != "" {
			existingOrg.Address.Country = cmd.Address.Country
		}
	}

	// Update contact info if any contact field is provided
	if cmd.ContactInfo.Name != "" || cmd.ContactInfo.Email != "" ||
		cmd.ContactInfo.PhoneNumber != "" || cmd.ContactInfo.Title != "" {

		if cmd.ContactInfo.Name != "" {
			existingOrg.ContactInfo.Name = cmd.ContactInfo.Name
		}
		if cmd.ContactInfo.Email != "" {
			existingOrg.ContactInfo.Email = cmd.ContactInfo.Email
		}
		if cmd.ContactInfo.PhoneNumber != "" {
			existingOrg.ContactInfo.PhoneNumber = cmd.ContactInfo.PhoneNumber
		}
		if cmd.ContactInfo.Title != "" {
			existingOrg.ContactInfo.Title = cmd.ContactInfo.Title
		}
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
