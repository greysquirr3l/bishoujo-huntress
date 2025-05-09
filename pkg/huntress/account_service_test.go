package huntress_test

import (
	"context"
	"testing"

	"github.com/greysquirr3l/bishoujo-huntress/pkg/huntress"
)

func TestAccountService_InviteUser(t *testing.T) {
	ctx := context.Background()
	client := huntress.New(huntress.WithBaseURL("http://localhost:12345")) // or use a test client
	svc := client.Account

	params := &huntress.UserInviteParams{
		Email:     "invitee@example.com",
		FirstName: "Invitee",
		LastName:  "User",
		Role:      huntress.UserRoleViewer,
	}

	// This will fail unless you have a test server running!
	user, err := svc.InviteUser(ctx, params)
	if err == nil {
		t.Logf("user: %+v", user)
	}
}

func TestUserCreateParams_Validate(t *testing.T) {
	p := &huntress.UserCreateParams{Email: ""}
	if err := p.Validate(); err == nil {
		t.Error("expected error for missing email")
	}
	p.Email = "user@example.com"
	p.Role = "bad"
	if err := p.Validate(); err == nil {
		t.Error("expected error for invalid role")
	}
	p.Role = huntress.UserRoleAdmin
	p.Roles = []huntress.UserRole{"bad"}
	if err := p.Validate(); err == nil {
		t.Error("expected error for invalid role in roles")
	}
	p.Roles = []huntress.UserRole{huntress.UserRoleViewer}
	if err := p.Validate(); err != nil {
		t.Errorf("expected no error for valid params, got %v", err)
	}
}

func TestUserUpdateParams_Validate(t *testing.T) {
	p := &huntress.UserUpdateParams{Role: "bad"}
	if err := p.Validate(); err == nil {
		t.Error("expected error for invalid role")
	}
	p.Role = huntress.UserRoleManager
	p.Roles = []huntress.UserRole{"bad"}
	if err := p.Validate(); err == nil {
		t.Error("expected error for invalid role in roles")
	}
	p.Roles = []huntress.UserRole{huntress.UserRoleAdmin}
	if err := p.Validate(); err != nil {
		t.Errorf("expected no error for valid params, got %v", err)
	}
}
