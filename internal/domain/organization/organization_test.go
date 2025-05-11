package organization

import (
	"testing"
)

func TestOrganization_Validate(t *testing.T) {
	tests := []struct {
		name    string
		o       Organization
		wantErr bool
	}{
		{"valid", Organization{Name: "Org", AccountID: 1, Status: StatusActive, ContactInfo: ContactInfo{Email: "a@b.com"}}, false},
		{"empty name", Organization{Name: "", AccountID: 1}, true},
		{"invalid account id", Organization{Name: "Org", AccountID: 0}, true},
		{"invalid status", Organization{Name: "Org", AccountID: 1, Status: "bad"}, true},
		{"invalid contact email", Organization{Name: "Org", AccountID: 1, ContactInfo: ContactInfo{Email: "bad"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.o.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("got error %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOrganization_IsActive(t *testing.T) {
	o := Organization{Status: StatusActive}
	if !o.IsActive() {
		t.Error("expected IsActive true for active status")
	}
	o.Status = StatusInactive
	if o.IsActive() {
		t.Error("expected IsActive false for inactive status")
	}
}

func TestOrganization_CanHaveAgents(t *testing.T) {
	o := Organization{Status: StatusActive}
	if !o.CanHaveAgents() {
		t.Error("expected CanHaveAgents true for active org")
	}
	o.Status = StatusInactive
	if o.CanHaveAgents() {
		t.Error("expected CanHaveAgents false for inactive org")
	}
}

func TestUser_FullName(t *testing.T) {
	u := User{FirstName: "Jane", LastName: "Doe"}
	if u.FullName() != "Jane Doe" {
		t.Errorf("expected 'Jane Doe', got %q", u.FullName())
	}
	u = User{FirstName: "Jane", LastName: ""}
	if u.FullName() != "Jane" {
		t.Errorf("expected 'Jane', got %q", u.FullName())
	}
}

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		u       User
		wantErr bool
	}{
		{"valid", User{Email: "a@b.com", Role: RoleAdmin}, false},
		{"empty email", User{Email: "", Role: RoleAdmin}, true},
		{"invalid email", User{Email: "bad", Role: RoleAdmin}, true},
		{"invalid role", User{Email: "a@b.com", Role: "bad"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.u.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("got error %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_HasAdminAccess(t *testing.T) {
	u := User{Role: RoleAdmin}
	if !u.HasAdminAccess() {
		t.Error("expected HasAdminAccess true for admin role")
	}
	u.Role = RoleViewer
	if u.HasAdminAccess() {
		t.Error("expected HasAdminAccess false for non-admin role")
	}
}

func TestUser_HasManagerAccess(t *testing.T) {
	u := User{Role: RoleAdmin}
	if !u.HasManagerAccess() {
		t.Error("expected HasManagerAccess true for admin role")
	}
	u.Role = RoleManager
	if !u.HasManagerAccess() {
		t.Error("expected HasManagerAccess true for manager role")
	}
	u.Role = RoleViewer
	if u.HasManagerAccess() {
		t.Error("expected HasManagerAccess false for viewer role")
	}
}

func Test_isValidStatus(t *testing.T) {
	if !isValidStatus(StatusActive) {
		t.Error("expected valid status")
	}
	if isValidStatus("bad") {
		t.Error("expected invalid status")
	}
}

func Test_isValidRole(t *testing.T) {
	if !isValidRole(RoleAdmin) {
		t.Error("expected valid role")
	}
	if isValidRole("bad") {
		t.Error("expected invalid role")
	}
}

func Test_isValidEmail(t *testing.T) {
	if !isValidEmail("a@b.com") {
		t.Error("expected valid email")
	}
	if isValidEmail("bad") {
		t.Error("expected invalid email")
	}
	if isValidEmail("a@bcom") {
		t.Error("expected invalid email")
	}
}
