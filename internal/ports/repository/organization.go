package repository

import (
	"context"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/organization"
)

// OrganizationRepository defines the interface for organization data operations
type OrganizationRepository interface {
	// Get retrieves a specific organization by its ID
	Get(ctx context.Context, id string) (*organization.Organization, error)

	// List retrieves all organizations with optional filtering
	List(ctx context.Context, filters map[string]interface{}) ([]*organization.Organization, *Pagination, error)

	// Create creates a new organization
	Create(ctx context.Context, org *organization.Organization) (*organization.Organization, error)

	// Update updates an existing organization
	Update(ctx context.Context, org *organization.Organization) (*organization.Organization, error)

	// Delete removes an organization by its ID
	Delete(ctx context.Context, id string) error
}
