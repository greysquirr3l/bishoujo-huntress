package organization

import (
	"context"
	"errors"
	"testing"

	orgdomain "github.com/greysquirr3l/bishoujo-huntress/internal/domain/organization"
	"github.com/stretchr/testify/assert"
)

func TestUpdateOrganizationHandler_Handle_Success(t *testing.T) {
	ctx := context.Background()
	const orgID = "org-123"
	cmd := UpdateOrganizationCommand{
		ID:   orgID,
		Name: "UpdatedOrg",
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
	org := &orgdomain.Organization{ID: orgID, AccountID: 1, Name: "UpdatedOrg", ContactInfo: orgdomain.ContactInfo{Name: "Contact", Email: "contact@example.com"}, Status: orgdomain.StatusActive}

	repo := &FakeOrganizationRepository{
		GetFunc: func(_ context.Context, id string) (*orgdomain.Organization, error) {
			if id == orgID {
				return &orgdomain.Organization{ID: orgID, AccountID: 1, Name: "OldOrg", ContactInfo: orgdomain.ContactInfo{}}, nil
			}
			return nil, errors.New("not found")
		},
		UpdateFunc: func(_ context.Context, o *orgdomain.Organization) (*orgdomain.Organization, error) {
			if o.ID == orgID && o.Name == "UpdatedOrg" && o.ContactInfo.Name == "Contact" && o.ContactInfo.Email == "contact@example.com" && o.AccountID == 1 {
				return org, nil
			}
			return nil, errors.New("update mismatch")
		},
	}
	h := NewUpdateOrganizationHandler(repo)
	res, err := h.Handle(ctx, cmd)
	assert.NoError(t, err)
	assert.Equal(t, org, res)
}

func TestUpdateOrganizationHandler_Handle_ValidationError(t *testing.T) {
	ctx := context.Background()
	repo := &FakeOrganizationRepository{}
	h := NewUpdateOrganizationHandler(repo)
	badCmd := UpdateOrganizationCommand{ID: "", Name: ""}
	_, err := h.Handle(ctx, badCmd)
	assert.Error(t, err)
}

func TestUpdateOrganizationHandler_Handle_RepoError(t *testing.T) {
	ctx := context.Background()
	cmd := UpdateOrganizationCommand{ID: "org-123", Name: "UpdatedOrg"}
	repo := &FakeOrganizationRepository{
		GetFunc: func(_ context.Context, id string) (*orgdomain.Organization, error) {
			return &orgdomain.Organization{ID: id, AccountID: 1, Name: "OldOrg"}, nil
		},
		UpdateFunc: func(_ context.Context, _ *orgdomain.Organization) (*orgdomain.Organization, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewUpdateOrganizationHandler(repo)
	_, err := h.Handle(ctx, cmd)
	assert.Error(t, err)
}
