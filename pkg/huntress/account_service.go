// Package huntress provides a client for the Huntress API
package huntress

import (
	"context"
	"fmt"
	"net/http"
)

// accountService implements the AccountService interface
type accountService struct {
	client *Client
}

// Get retrieves current account details
func (s *accountService) Get(ctx context.Context) (*Account, error) {
	path := "/account"
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for Get: %w", err)
	}

	account := new(Account)
	resp, err := s.client.Do(ctx, req, account)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for Get: %w", err)
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}

	return account, nil
}

// GetCurrent is an alias for Get to maintain backward compatibility
func (s *accountService) GetCurrent(ctx context.Context) (*Account, error) {
	return s.Get(ctx)
}

// Update updates account settings
func (s *accountService) Update(ctx context.Context, account *AccountUpdateParams) (*Account, error) {
	path := "/account"
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, account)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for Update: %w", err)
	}

	updatedAccount := new(Account)
	resp, err := s.client.Do(ctx, req, updatedAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for Update: %w", err)
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}

	return updatedAccount, nil
}

// GetStats retrieves account statistics
func (s *accountService) GetStats(ctx context.Context) (*AccountStats, error) {
	path := "/account/stats"
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for GetStats: %w", err)
	}

	stats := new(AccountStats)
	resp, err := s.client.Do(ctx, req, stats)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for GetStats: %w", err)
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}

	return stats, nil
}

// GetStatistics is an alias for GetStats to maintain backward compatibility
func (s *accountService) GetStatistics(ctx context.Context) (*AccountStats, error) {
	return s.GetStats(ctx)
}

// ListUsers lists users associated with an account
func (s *accountService) ListUsers(ctx context.Context, params *ListParams) ([]*User, *Pagination, error) {
	var users []*User
	pagination, err := listResource(ctx, s.client, "/account/users", params, &users)
	if err != nil {
		return nil, nil, err
	}
	return users, pagination, nil
}

// AddUser adds a new user to the account
func (s *accountService) AddUser(ctx context.Context, user *UserCreateParams) (*User, error) {
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("invalid user params: %w", err)
	}
	path := "/users"
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for AddUser: %w", err)
	}

	newUser := new(User)
	resp, err := s.client.Do(ctx, req, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for AddUser: %w", err)
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}

	return newUser, nil
}

// UpdateUser updates an existing user
func (s *accountService) UpdateUser(ctx context.Context, userID string, user *UserUpdateParams) (*User, error) {
	if err := user.Validate(); err != nil {
		return nil, fmt.Errorf("invalid user params: %w", err)
	}
	path := fmt.Sprintf("/users/%s", userID)
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for UpdateUser: %w", err)
	}

	updatedUser := new(User)
	resp, err := s.client.Do(ctx, req, updatedUser)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request for UpdateUser: %w", err)
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}

	return updatedUser, nil
}

// RemoveUser removes a user from the account
func (s *accountService) RemoveUser(ctx context.Context, userID string) error {
	path := fmt.Sprintf("/users/%s", userID)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request for RemoveUser: %w", err)
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return fmt.Errorf("failed to execute request for RemoveUser: %w", err)
	}
	if resp != nil {
		defer func() { _ = resp.Body.Close() }()
	}
	return nil
}

// Note: The utility functions (addQueryParams, extractPagination, parseInt)
// have been moved to utils.go to prevent redeclaration errors
