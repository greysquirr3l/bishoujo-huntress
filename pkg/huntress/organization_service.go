package huntress

import (
	"context"
	"fmt"
	"net/http"
)

// organizationService implements the OrganizationService interface
type organizationService struct {
	client *Client
}

// Get retrieves organization details by ID
func (s *organizationService) Get(ctx context.Context, id int) (*Organization, error) {
	path := fmt.Sprintf("/organizations/%d", id)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	org := new(Organization)
	_, err = s.client.Do(ctx, req, org)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// Create creates a new organization
func (s *organizationService) Create(ctx context.Context, input *OrganizationCreateInput) (*Organization, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "/organizations", input)
	if err != nil {
		return nil, err
	}

	org := new(Organization)
	_, err = s.client.Do(ctx, req, org)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// Update updates an organization
func (s *organizationService) Update(ctx context.Context, id int, input *OrganizationUpdateInput) (*Organization, error) {
	path := fmt.Sprintf("/organizations/%d", id)
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, input)
	if err != nil {
		return nil, err
	}

	org := new(Organization)
	_, err = s.client.Do(ctx, req, org)
	if err != nil {
		return nil, err
	}

	return org, nil
}

// Delete deletes an organization
func (s *organizationService) Delete(ctx context.Context, id int) error {
	path := fmt.Sprintf("/organizations/%d", id)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List lists organizations with optional filtering
func (s *organizationService) List(ctx context.Context, opts *OrganizationListOptions) ([]*Organization, *Pagination, error) {
	path := "/organizations"
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var orgs []*Organization
	resp, err := s.client.Do(ctx, req, &orgs)
	if err != nil {
		return nil, nil, err
	}

	pagination := extractPagination(resp)
	return orgs, pagination, nil
}

// ListUsers lists users associated with an organization
func (s *organizationService) ListUsers(ctx context.Context, id int, opts *ListOptions) ([]*User, *Pagination, error) {
	path := fmt.Sprintf("/organizations/%d/users", id)
	path, err := addOptions(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var users []*User
	resp, err := s.client.Do(ctx, req, &users)
	if err != nil {
		return nil, nil, err
	}

	pagination := extractPagination(resp)
	return users, pagination, nil
}

// GetStatistics retrieves organization statistics
func (s *organizationService) GetStatistics(ctx context.Context, id int) (*OrganizationStatistics, error) {
	path := fmt.Sprintf("/organizations/%d/statistics", id)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	stats := new(OrganizationStatistics)
	_, err = s.client.Do(ctx, req, stats)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
