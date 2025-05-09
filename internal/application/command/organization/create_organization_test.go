package organization

import (
	"context"
	"errors"
	"testing"

	orgdomain "github.com/greysquirr3l/bishoujo-huntress/internal/domain/organization"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrganizationHandler_Handle_Success(t *testing.T) {
	ctx := context.Background()
	cmd := CreateOrganizationCommand{
		AccountID: 1,
		Name:      "TestOrg",
		ContactInfo: struct {
			Name        string
			Email       string
			PhoneNumber string
			Title       string
		}{
			Name:  "Contact",
			Email: "contact@example.com",
		},
	}
	org := &orgdomain.Organization{AccountID: 1, Name: "TestOrg", ContactInfo: orgdomain.ContactInfo{Name: "Contact", Email: "contact@example.com"}, Status: orgdomain.StatusActive}

	repo := &FakeOrganizationRepository{
		CreateFunc: func(_ context.Context, _ *orgdomain.Organization) (*orgdomain.Organization, error) {
			return org, nil
		},
	}
	h := NewCreateOrganizationHandler(repo)
	res, err := h.Handle(ctx, cmd)
	assert.NoError(t, err)
	assert.Equal(t, org, res)
}

func TestCreateOrganizationHandler_Handle_ValidationError(t *testing.T) {
	ctx := context.Background()
	repo := &FakeOrganizationRepository{}
	h := NewCreateOrganizationHandler(repo)
	badCmd := CreateOrganizationCommand{AccountID: 1, Name: ""}
	_, err := h.Handle(ctx, badCmd)
	assert.Error(t, err)
}

func TestCreateOrganizationHandler_Handle_RepoError(t *testing.T) {
	ctx := context.Background()
	cmd := CreateOrganizationCommand{AccountID: 1, Name: "TestOrg"}
	repo := &FakeOrganizationRepository{
		CreateFunc: func(_ context.Context, _ *orgdomain.Organization) (*orgdomain.Organization, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewCreateOrganizationHandler(repo)
	_, err := h.Handle(ctx, cmd)
	assert.Error(t, err)
}
