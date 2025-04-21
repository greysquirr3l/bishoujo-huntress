package repository

import (
	"context"

	"github.com/greysquirr3l18/bishoujo-huntress/internal/domain/organization"
)

// OrganizationRepository defines the interface for organization data operations
type OrganizationRepository interface {
	// Get retrieves a specific organization by its ID
	Get(ctx context.Context, id string) (*organization.Organization, error)

	// List retrieves all organizations with optional filtering
	List(ctx context.Context, params *organization.ListParams) ([]*organization.Organization, *Pagination, error)

	// Create creates a new organization
	Create(ctx context.Context, org *organization.Organization) (*organization.Organization, error)

	// Update updates an existing organization
	Update(ctx context.Context, org *organization.Organization) (*organization.Organization, error)

	// Delete removes an organization by its ID
	Delete(ctx context.Context, id string) error

	// GetUsers retrieves all users associated with an organization
	GetUsers(ctx context.Context, organizationID string) ([]*organization.User, *Pagination, error)

	// AddUser adds a user to an organization
	AddUser(ctx context.Context, organizationID string, user *organization.User) error

	// RemoveUser removes a user from an organization
	RemoveUser(ctx context.Context, organizationID string, userID string) error
}
