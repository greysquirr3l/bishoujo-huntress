package huntress

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/go-querystring/query"
)

// accountService implements the AccountService interface
type accountService struct {
	client *Client
}

// Get retrieves account details by ID
func (s *accountService) Get(ctx context.Context, id int) (*Account, error) {
	path := fmt.Sprintf("/accounts/%d", id)
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

// GetCurrent retrieves the current authenticated account
func (s *accountService) GetCurrent(ctx context.Context) (*Account, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "/accounts/current", nil)
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

// Update updates account settings
func (s *accountService) Update(ctx context.Context, id int, input *AccountUpdateInput) (*Account, error) {
	path := fmt.Sprintf("/accounts/%d", id)
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, input)
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

// ListUsers lists users associated with an account
func (s *accountService) ListUsers(ctx context.Context, id int, opts *ListOptions) ([]*User, *Pagination, error) {
	path := fmt.Sprintf("/accounts/%d/users", id)
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

// GetStatistics retrieves account statistics
func (s *accountService) GetStatistics(ctx context.Context, id int) (*AccountStatistics, error) {
	path := fmt.Sprintf("/accounts/%d/statistics", id)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	stats := new(AccountStatistics)
	_, err = s.client.Do(ctx, req, stats)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// addOptions adds the parameters in opts as URL query parameters to s.
// opts must be a struct whose fields may contain "url" tags.
func addOptions(s string, opts interface{}) (string, error) {
	if opts == nil {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// extractPagination extracts pagination information from HTTP response headers
func extractPagination(resp *http.Response) *Pagination {
	if resp == nil {
		return nil
	}

	pagination := &Pagination{}

	if currentPage := resp.Header.Get("X-Page"); currentPage != "" {
		pagination.CurrentPage, _ = strconv.Atoi(currentPage)
	}

	if perPage := resp.Header.Get("X-Per-Page"); perPage != "" {
		pagination.PerPage, _ = strconv.Atoi(perPage)
	}

	if totalPages := resp.Header.Get("X-Total-Pages"); totalPages != "" {
		pagination.TotalPages, _ = strconv.Atoi(totalPages)
	}

	if totalItems := resp.Header.Get("X-Total-Items"); totalItems != "" {
		pagination.TotalItems, _ = strconv.Atoi(totalItems)
	}

	return pagination
}
