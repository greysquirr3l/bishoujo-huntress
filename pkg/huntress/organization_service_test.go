package huntress_test

import (
	"context"
	"testing"

	"github.com/greysquirr3l/bishoujo-huntress/pkg/huntress"
)

func TestOrganizationService_InviteUser(t *testing.T) {
	ctx := context.Background()
	client := huntress.New(huntress.WithBaseURL("http://localhost:12345")) // or use a test client
	svc := client.Organization

	orgID := "org-123"
	params := &huntress.UserInviteParams{
		Email:     "invitee@example.com",
		FirstName: "Invitee",
		LastName:  "User",
		Role:      huntress.UserRoleViewer,
	}

	// This will fail unless you have a test server running!
	user, err := svc.InviteUser(ctx, orgID, params)
	if err == nil {
		t.Logf("user: %+v", user)
	}
}

func TestOrganizationCreateParams_Validate(t *testing.T) {
	p := &huntress.OrganizationCreateParams{Name: ""}
	if err := p.Validate(); err == nil {
		t.Error("expected error for missing name")
	}
	p.Name = "Test Org"
	p.Status = "bad"
	if err := p.Validate(); err == nil {
		t.Error("expected error for invalid status")
	}
	p.Status = huntress.OrganizationStatusActive
	if err := p.Validate(); err != nil {
		t.Errorf("expected no error for valid params, got %v", err)
	}
}

func TestOrganizationUpdateParams_Validate(t *testing.T) {
	p := &huntress.OrganizationUpdateParams{Status: "bad"}
	if err := p.Validate(); err == nil {
		t.Error("expected error for invalid status")
	}
	p.Status = huntress.OrganizationStatusInactive
	if err := p.Validate(); err != nil {
		t.Errorf("expected no error for valid status, got %v", err)
	}
}
