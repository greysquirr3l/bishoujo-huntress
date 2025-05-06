package huntress_test

import (
	"testing"

	"github.com/greysquirr3l/bishoujo-huntress/pkg/huntress"
)

func TestAccountService_InviteUser(t *testing.T) {
	client := newMockClient()
	svc := &huntress.AccountService{client}

	params := &huntress.UserInviteParams{
		Email:     "invitee@example.com",
		FirstName: "Invitee",
		LastName:  "User",
		Role:      huntress.UserRoleViewer,
	}

	user, err := svc.InviteUser(ctx, params)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil || user.Email != params.Email {
		t.Errorf("expected invited user with email %s, got %+v", params.Email, user)
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
