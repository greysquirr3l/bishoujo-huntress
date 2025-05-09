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
		return nil, fmt.Errorf("failed to create request for Get: %w", err)
	}

	org := new(Organization)
	resp, err := s.client.Do(ctx, req, org)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for Get: %w", err)
	}
	if resp != nil {
		errClose := resp.Body.Close()
		if errClose != nil {
			return nil, fmt.Errorf("organization get: error closing response body: %w", errClose)
		}
	}

	return org, nil
}

// List returns all organizations with optional filtering
func (s *organizationService) List(ctx context.Context, params *ListOrganizationsParams) ([]*Organization, *Pagination, error) {
	var orgs []*Organization
	pagination, err := listResource(ctx, s.client, "/organizations", params, &orgs)
	if err != nil {
		return nil, nil, err
	}
	return orgs, pagination, nil
}

// Create creates a new organization
func (s *organizationService) Create(ctx context.Context, org *OrganizationCreateParams) (*Organization, error) {
	if err := org.Validate(); err != nil {
		return nil, fmt.Errorf("invalid organization params: %w", err)
	}
	path := "/organizations"
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, org)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for Create: %w", err)
	}

	createdOrg := new(Organization)
	resp, err := s.client.Do(ctx, req, createdOrg)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for Create: %w", err)
	}
	if resp != nil {
		errClose := resp.Body.Close()
		if errClose != nil {
			return nil, fmt.Errorf("organization create: error closing response body: %w", errClose)
		}
	}

	return createdOrg, nil
}

// Update updates an existing organization
func (s *organizationService) Update(ctx context.Context, id string, org *OrganizationUpdateParams) (*Organization, error) {
	if err := org.Validate(); err != nil {
		return nil, fmt.Errorf("invalid organization params: %w", err)
	}
	path := fmt.Sprintf("/organizations/%s", id)
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, org)
	if err != nil {
		return nil, err
	}

	updatedOrg := new(Organization)
	resp, err := s.client.Do(ctx, req, updatedOrg)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		errClose := resp.Body.Close()
		if errClose != nil {
			return nil, fmt.Errorf("organization update: error closing response body: %w", errClose)
		}
	}

	return updatedOrg, nil
}

// Delete removes an organization
func (s *organizationService) Delete(ctx context.Context, orgID string) error {
	path := fmt.Sprintf("/organizations/%s", orgID)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request for Delete: %w", err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return fmt.Errorf("failed to execute request for Delete: %w", err)
	}
	if resp != nil {
		errClose := resp.Body.Close()
		if errClose != nil {
			return fmt.Errorf("organization delete: error closing response body: %w", errClose)
		}
	}
	return nil
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
	if resp != nil {
		errClose := resp.Body.Close()
		if errClose != nil {
			return nil, nil, fmt.Errorf("organization list: error closing response body: %w", errClose)
		}
	}

	pagination := extractPagination(resp)
	return users, pagination, nil
}

// AddUser adds a user to an organization
func (s *organizationService) AddUser(ctx context.Context, orgID string, user *UserCreateParams) (*User, error) {
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("invalid user params: %w", err)
	}
	path := fmt.Sprintf("/organizations/%s/users", orgID)
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, user)
	if err != nil {
		return nil, err
	}

	newUser := new(User)
	resp, err := s.client.Do(ctx, req, newUser)
	if err != nil {
		return nil, err
	}
	if resp != nil {
		errClose := resp.Body.Close()
		if errClose != nil {
			return nil, fmt.Errorf("organization add user: error closing response body: %w", errClose)
		}
	}

	return newUser, nil
}

// RemoveUser removes a user from an organization
func (s *organizationService) RemoveUser(ctx context.Context, orgID string, userID string) error {
	path := fmt.Sprintf("/organizations/%s/users/%s", orgID, userID)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request for RemoveUser: %w", err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return fmt.Errorf("failed to execute request for RemoveUser: %w", err)
	}
	if resp != nil {
		errClose := resp.Body.Close()
		if errClose != nil {
			return fmt.Errorf("organization remove user: error closing response body: %w", errClose)
		}
	}
	return nil
}

// InviteUser sends an invitation to a user for an organization (if supported by the Huntress API)
func (s *organizationService) InviteUser(ctx context.Context, orgID string, params *UserInviteParams) (*User, error) {
	if params == nil {
		return nil, fmt.Errorf("invite params required")
	}
	// The endpoint is assumed to be /organizations/{orgID}/users/invite
	path := fmt.Sprintf("/organizations/%s/users/invite", orgID)
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for InviteUser: %w", err)
	}

	newUser := new(User)
	resp, err := s.client.Do(ctx, req, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for InviteUser: %w", err)
	}
	if resp != nil {
		errClose := resp.Body.Close()
		if errClose != nil {
			return nil, fmt.Errorf("organization invite user: error closing response body: %w", errClose)
		}
	}

	return newUser, nil
}

// URL parameter handling functions have been moved to utils.go
// for better code organization and to avoid duplication.

// parseTag splits a struct field's url tag into its name and options.
// These utility functions have been moved to utils.go
// to maintain DRY principle and avoid code duplication.
