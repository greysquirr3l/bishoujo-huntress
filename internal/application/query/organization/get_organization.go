// Package organization contains query handlers for organizations
package organization

import (
	"context"
	"fmt"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/organization"
	"github.com/greysquirr3l/bishoujo-huntress/internal/ports/repository"
)

// GetOrganizationQuery represents a query to retrieve a single organization by ID
type GetOrganizationQuery struct {
	ID string // Organization ID to retrieve
}

// GetOrganizationHandler handles the get organization query
type GetOrganizationHandler struct {
	orgRepo repository.OrganizationRepository
}

// NewGetOrganizationHandler creates a new get organization handler
func NewGetOrganizationHandler(orgRepo repository.OrganizationRepository) *GetOrganizationHandler {
	return &GetOrganizationHandler{
		orgRepo: orgRepo,
	}
}

// Handle executes the get organization query
func (h *GetOrganizationHandler) Handle(ctx context.Context, query GetOrganizationQuery) (*organization.Organization, error) {
	if query.ID == "" {
		return nil, fmt.Errorf("organization ID is required")
	}

	// Call repository to get the organization
	org, err := h.orgRepo.Get(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	return org, nil
}
