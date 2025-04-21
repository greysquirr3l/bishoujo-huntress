package repository

import (
	"context"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/organization"
)

// OrganizationFilters defines filters for organization queries
type OrganizationFilters struct {
	AccountID *int
	Search    string
	Status    *organization.OrganizationStatus
	Industry  string
	Tags      []string
	Page      int
	Limit     int
	OrderBy   []OrderBy
	TimeRange *TimeRange
}

// OrganizationRepository defines the repository interface for organizations
type OrganizationRepository interface {
	// Get retrieves an organization by ID
	Get(ctx context.Context, id int) (*organization.Organization, error)

	// List retrieves multiple organizations based on filters
	List(ctx context.Context, filters OrganizationFilters) ([]*organization.Organization, Pagination, error)

	// Create creates a new organization
	Create(ctx context.Context, org *organization.Organization) error

	// Update updates an existing organization
	Update(ctx context.Context, org *organization.Organization) error

	// Delete deletes an organization by ID
	Delete(ctx context.Context, id int) error

	// GetByName retrieves an organization by name within an account
	GetByName(ctx context.Context, accountID int, name string) (*organization.Organization, error)

	// GetStatistics retrieves organization statistics
	GetStatistics(ctx context.Context, id int) (*organization.Statistics, error)

	// UpdateTags updates the tags for an organization
	UpdateTags(ctx context.Context, id int, tags []string) error
}
