package repository

import (
	"context"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/account"
)

// AccountFilters defines filters for account queries
type AccountFilters struct {
	Search    string
	Status    *account.Status
	Page      int
	Limit     int
	OrderBy   []OrderBy
	TimeRange *TimeRange
}

// AccountRepository defines the repository interface for accounts
type AccountRepository interface {
	// Get retrieves an account by ID
	Get(ctx context.Context, id int) (*account.Account, error)

	// List retrieves multiple accounts based on filters
	List(ctx context.Context, filters AccountFilters) ([]*account.Account, Pagination, error)

	// Create creates a new account
	Create(ctx context.Context, account *account.Account) error

	// Update updates an existing account
	Update(ctx context.Context, account *account.Account) error

	// Delete deletes an account by ID
	Delete(ctx context.Context, id int) error

	// GetByName retrieves an account by name
	GetByName(ctx context.Context, name string) (*account.Account, error)

	// GetStatistics retrieves account statistics
	GetStatistics(ctx context.Context, id int) (map[string]interface{}, error)
}
