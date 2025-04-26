// Package organization contains query handlers for organizations
package organization

import (
	"context"
	"fmt"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/organization"
	"github.com/greysquirr3l/bishoujo-huntress/internal/ports/repository"
)

// ListOrganizationsQuery represents a query to list organizations with filters
type ListOrganizationsQuery struct {
	Page      int      // Page number for pagination
	Limit     int      // Number of items per page
	AccountID int      // Filter by account ID
	Status    string   // Filter by status
	Search    string   // Search term
	Industry  string   // Filter by industry
	Tags      []string // Filter by tags
}

// ListOrganizationsResult represents the result of a list organizations query
type ListOrganizationsResult struct {
	Organizations []*organization.Organization
	Pagination    *repository.Pagination
}

// ListOrganizationsHandler handles the list organizations query
type ListOrganizationsHandler struct {
	orgRepo repository.OrganizationRepository
}

// NewListOrganizationsHandler creates a new list organizations handler
func NewListOrganizationsHandler(orgRepo repository.OrganizationRepository) *ListOrganizationsHandler {
	return &ListOrganizationsHandler{
		orgRepo: orgRepo,
	}
}

// Handle executes the list organizations query
func (h *ListOrganizationsHandler) Handle(ctx context.Context, query ListOrganizationsQuery) (*ListOrganizationsResult, error) {
	// Map query to domain list params
	params := &organization.ListParams{
		Page:      query.Page,
		Limit:     query.Limit,
		AccountID: query.AccountID,
		Status:    query.Status,
		Search:    query.Search,
		Industry:  query.Industry,
		Tags:      query.Tags,
	}

	// Convert params to map[string]interface{} for the repository
	filtersMap := map[string]interface{}{
		"page":       params.Page,
		"limit":      params.Limit,
		"account_id": params.AccountID,
	}

	// Only add optional filters if they have values
	if params.Status != "" {
		filtersMap["status"] = params.Status
	}
	if params.Search != "" {
		filtersMap["search"] = params.Search
	}
	if params.Industry != "" {
		filtersMap["industry"] = params.Industry
	}
	if len(params.Tags) > 0 {
		filtersMap["tags"] = params.Tags
	}

	// Call repository to get organizations
	orgs, pagination, err := h.orgRepo.List(ctx, filtersMap)
	if err != nil {
		return nil, fmt.Errorf("failed to list organizations: %w", err)
	}

	return &ListOrganizationsResult{
		Organizations: orgs,
		Pagination:    pagination,
	}, nil
}
