// Package organization contains command handlers for organization operations
package organization

import (
	"context"
	"fmt"

	"github.com/greysquirr3l/bishoujo-huntress/internal/ports/repository"
)

// DeleteOrganizationCommand represents a request to delete an organization
type DeleteOrganizationCommand struct {
	ID string // Organization ID to delete
}

// DeleteOrganizationHandler handles deletion of organizations
type DeleteOrganizationHandler struct {
	orgRepo repository.OrganizationRepository
}

// NewDeleteOrganizationHandler creates a new organization deletion handler
func NewDeleteOrganizationHandler(orgRepo repository.OrganizationRepository) *DeleteOrganizationHandler {
	return &DeleteOrganizationHandler{
		orgRepo: orgRepo,
	}
}

// Handle processes the delete organization command
func (h *DeleteOrganizationHandler) Handle(ctx context.Context, cmd DeleteOrganizationCommand) error {
	// Verify that organization ID is provided
	if cmd.ID == "" {
		return fmt.Errorf("organization ID is required for deletion")
	}

	// Call repository to delete the organization
	if err := h.orgRepo.Delete(ctx, cmd.ID); err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	return nil
}
