// org_fake_test.go centralizes a flexible fake for OrganizationRepository for all command handler tests.
package organization

import (
	"context"

	orgdomain "github.com/greysquirr3l/bishoujo-huntress/internal/domain/organization"
	repo "github.com/greysquirr3l/bishoujo-huntress/internal/ports/repository"
)

type FakeOrganizationRepository struct {
	CreateFunc func(ctx context.Context, org *orgdomain.Organization) (*orgdomain.Organization, error)
	GetFunc    func(ctx context.Context, id string) (*orgdomain.Organization, error)
	UpdateFunc func(ctx context.Context, org *orgdomain.Organization) (*orgdomain.Organization, error)
	DeleteFunc func(ctx context.Context, id string) error
	ListFunc   func(ctx context.Context, filters map[string]interface{}) ([]*orgdomain.Organization, *repo.Pagination, error)
}

func (f *FakeOrganizationRepository) Create(ctx context.Context, org *orgdomain.Organization) (*orgdomain.Organization, error) {
	if f.CreateFunc != nil {
		return f.CreateFunc(ctx, org)
	}
	return nil, nil
}
func (f *FakeOrganizationRepository) Get(ctx context.Context, id string) (*orgdomain.Organization, error) {
	if f.GetFunc != nil {
		return f.GetFunc(ctx, id)
	}
	return nil, nil
}
func (f *FakeOrganizationRepository) Update(ctx context.Context, org *orgdomain.Organization) (*orgdomain.Organization, error) {
	if f.UpdateFunc != nil {
		return f.UpdateFunc(ctx, org)
	}
	return nil, nil
}
func (f *FakeOrganizationRepository) Delete(ctx context.Context, id string) error {
	if f.DeleteFunc != nil {
		return f.DeleteFunc(ctx, id)
	}
	return nil
}
func (f *FakeOrganizationRepository) List(ctx context.Context, filters map[string]interface{}) ([]*orgdomain.Organization, *repo.Pagination, error) {
	if f.ListFunc != nil {
		return f.ListFunc(ctx, filters)
	}
	return nil, nil, nil
}
