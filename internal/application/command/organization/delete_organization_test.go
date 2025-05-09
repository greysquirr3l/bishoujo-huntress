package organization

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteOrganizationHandler_Handle_Success(t *testing.T) {
	ctx := context.Background()
	cmd := DeleteOrganizationCommand{ID: "org-123"}
	repo := &FakeOrganizationRepository{
		DeleteFunc: func(_ context.Context, id string) error {
			if id == "org-123" {
				return nil
			}
			return errors.New("unexpected id")
		},
	}
	h := NewDeleteOrganizationHandler(repo)
	err := h.Handle(ctx, cmd)
	assert.NoError(t, err)
}

func TestDeleteOrganizationHandler_Handle_MissingID(t *testing.T) {
	ctx := context.Background()
	repo := &FakeOrganizationRepository{}
	h := NewDeleteOrganizationHandler(repo)
	badCmd := DeleteOrganizationCommand{ID: ""}
	err := h.Handle(ctx, badCmd)
	assert.Error(t, err)
}

func TestDeleteOrganizationHandler_Handle_RepoError(t *testing.T) {
	ctx := context.Background()
	cmd := DeleteOrganizationCommand{ID: "org-123"}
	repo := &FakeOrganizationRepository{
		DeleteFunc: func(_ context.Context, id string) error {
			return errors.New("db error")
		},
	}
	h := NewDeleteOrganizationHandler(repo)
	err := h.Handle(ctx, cmd)
	assert.Error(t, err)
}
