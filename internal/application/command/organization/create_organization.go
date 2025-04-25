// Package organization contains command handlers for organization operations
package organization

import (
	"context"
	"fmt"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/organization"
	"github.com/greysquirr3l/bishoujo-huntress/internal/ports/repository"
)

// CreateOrganizationCommand represents a request to create a new organization
type CreateOrganizationCommand struct {
	AccountID   int                    // Parent account identifier
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

// CreateOrganizationHandler handles creation of new organizations
type CreateOrganizationHandler struct {
	orgRepo repository.OrganizationRepository
}

// NewCreateOrganizationHandler creates a new organization creation handler
func NewCreateOrganizationHandler(orgRepo repository.OrganizationRepository) *CreateOrganizationHandler {
	return &CreateOrganizationHandler{
		orgRepo: orgRepo,
	}
}

// Handle processes the create organization command
func (h *CreateOrganizationHandler) Handle(ctx context.Context, cmd CreateOrganizationCommand) (*organization.Organization, error) {
	// Create domain entity from command
	org := &organization.Organization{
		AccountID:   cmd.AccountID,
		Name:        cmd.Name,
		Description: cmd.Description,
		Status:      cmd.Status,
		Industry:    cmd.Industry,
		Tags:        cmd.Tags,
		Settings:    cmd.Settings,
		Address: organization.Address{
			Street1: cmd.Address.Street1,
			Street2: cmd.Address.Street2,
			City:    cmd.Address.City,
			State:   cmd.Address.State,
			ZipCode: cmd.Address.ZipCode,
			Country: cmd.Address.Country,
		},
		ContactInfo: organization.ContactInfo{
			Name:        cmd.ContactInfo.Name,
			Email:       cmd.ContactInfo.Email,
			PhoneNumber: cmd.ContactInfo.PhoneNumber,
			Title:       cmd.ContactInfo.Title,
		},
	}

	// Apply default status if not provided
	if org.Status == "" {
		org.Status = organization.StatusActive
	}

	// Validate the domain entity
	if err := org.Validate(); err != nil {
		return nil, fmt.Errorf("invalid organization: %w", err)
	}

	// Call repository to create organization
	createdOrg, err := h.orgRepo.Create(ctx, org)
	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	return createdOrg, nil
}
