// Package huntress provides a client for the Huntress API
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
func (s *organizationService) Get(ctx context.Context, id string) (*Organization, error) {
	path := fmt.Sprintf("/organizations/%s", id)
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

// List returns all organizations with optional filtering
func (s *organizationService) List(ctx context.Context, params *ListOrganizationsParams) ([]*Organization, *Pagination, error) {
	path := "/organizations"
	if params != nil {
		query, err := addQueryParams(path, params)
		if err != nil {
			return nil, nil, err
		}
		path = query
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

// Create creates a new organization
func (s *organizationService) Create(ctx context.Context, org *OrganizationCreateParams) (*Organization, error) {
	path := "/organizations"
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, org)
	if err != nil {
		return nil, err
	}

	newOrg := new(Organization)
	_, err = s.client.Do(ctx, req, newOrg)
	if err != nil {
		return nil, err
	}

	return newOrg, nil
}

// Update updates an existing organization
func (s *organizationService) Update(ctx context.Context, id string, org *OrganizationUpdateParams) (*Organization, error) {
	path := fmt.Sprintf("/organizations/%s", id)
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, org)
	if err != nil {
		return nil, err
	}

	updatedOrg := new(Organization)
	_, err = s.client.Do(ctx, req, updatedOrg)
	if err != nil {
		return nil, err
	}

	return updatedOrg, nil
}

// Delete removes an organization
func (s *organizationService) Delete(ctx context.Context, id string) error {
	path := fmt.Sprintf("/organizations/%s", id)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// ListUsers returns all users in an organization with optional filtering
func (s *organizationService) ListUsers(ctx context.Context, orgID string, params *ListParams) ([]*User, *Pagination, error) {
	path := fmt.Sprintf("/organizations/%s/users", orgID)
	if params != nil {
		query, err := addQueryParams(path, params)
		if err != nil {
			return nil, nil, err
		}
		path = query
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

// AddUser adds a user to an organization
func (s *organizationService) AddUser(ctx context.Context, orgID string, user *UserCreateParams) (*User, error) {
	path := fmt.Sprintf("/organizations/%s/users", orgID)
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, user)
	if err != nil {
		return nil, err
	}

	newUser := new(User)
	_, err = s.client.Do(ctx, req, newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// RemoveUser removes a user from an organization
func (s *organizationService) RemoveUser(ctx context.Context, orgID string, userID string) error {
	path := fmt.Sprintf("/organizations/%s/users/%s", orgID, userID)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// URL parameter handling functions have been moved to utils.go
// for better code organization and to avoid duplication.

// parseTag splits a struct field's url tag into its name and options.
// These utility functions have been moved to utils.go
// to maintain DRY principle and avoid code duplication.
