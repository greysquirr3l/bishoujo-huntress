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
		return nil, err
	}

	account := new(Account)
	_, err = s.client.Do(ctx, req, account)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	updatedAccount := new(Account)
	_, err = s.client.Do(ctx, req, updatedAccount)
	if err != nil {
		return nil, err
	}

	return updatedAccount, nil
}

// GetStats retrieves account statistics
func (s *accountService) GetStats(ctx context.Context) (*AccountStats, error) {
	path := "/account/stats"
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	stats := new(AccountStats)
	_, err = s.client.Do(ctx, req, stats)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// GetStatistics is an alias for GetStats to maintain backward compatibility
func (s *accountService) GetStatistics(ctx context.Context) (*AccountStats, error) {
	return s.GetStats(ctx)
}

// ListUsers lists users associated with an account
func (s *accountService) ListUsers(ctx context.Context, params *ListParams) ([]*User, *Pagination, error) {
	path := "/account/users"
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

// AddUser adds a new user to the account
func (s *accountService) AddUser(ctx context.Context, user *UserCreateParams) (*User, error) {
	path := "/users"
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

// UpdateUser updates an existing user
func (s *accountService) UpdateUser(ctx context.Context, userID string, user *UserUpdateParams) (*User, error) {
	path := fmt.Sprintf("/users/%s", userID)
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, user)
	if err != nil {
		return nil, err
	}

	updatedUser := new(User)
	_, err = s.client.Do(ctx, req, updatedUser)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// RemoveUser removes a user from the account
func (s *accountService) RemoveUser(ctx context.Context, userID string) error {
	path := fmt.Sprintf("/users/%s", userID)
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// Note: The utility functions (addQueryParams, extractPagination, parseInt)
// have been moved to utils.go to prevent redeclaration errors
