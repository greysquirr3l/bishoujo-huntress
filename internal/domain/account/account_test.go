package account

import (
	"testing"

	"github.com/google/uuid"
)

func TestAccount_Validate(t *testing.T) {
	tests := []struct {
		name    string
		acct    Account
		wantErr error
	}{
		{
			name:    "valid account",
			acct:    Account{ID: uuid.New(), Name: "Test", PrimaryContact: Contact{Name: "A", Email: "a@b.com"}, BillingContact: Contact{Name: "B", Email: "b@b.com"}},
			wantErr: nil,
		},
		{
			name:    "empty name",
			acct:    Account{ID: uuid.New(), Name: "", PrimaryContact: Contact{Name: "A", Email: "a@b.com"}, BillingContact: Contact{Name: "B", Email: "b@b.com"}},
			wantErr: ErrEmptyName,
		},
		{
			name:    "invalid primary contact",
			acct:    Account{ID: uuid.New(), Name: "Test", PrimaryContact: Contact{Name: "", Email: "a@b.com"}, BillingContact: Contact{Name: "B", Email: "b@b.com"}},
			wantErr: ErrEmptyContactName,
		},
		{
			name:    "invalid billing contact",
			acct:    Account{ID: uuid.New(), Name: "Test", PrimaryContact: Contact{Name: "A", Email: "a@b.com"}, BillingContact: Contact{Name: "", Email: "b@b.com"}},
			wantErr: ErrEmptyContactName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.acct.Validate()
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestContact_Validate(t *testing.T) {
	tests := []struct {
		name    string
		c       Contact
		wantErr error
	}{
		{"valid", Contact{Name: "A", Email: "a@b.com"}, nil},
		{"empty name", Contact{Name: "", Email: "a@b.com"}, ErrEmptyContactName},
		{"empty email", Contact{Name: "A", Email: ""}, ErrEmptyContactEmail},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.Validate()
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestAccount_IsActive(t *testing.T) {
	a := Account{Status: StatusActive}
	if !a.IsActive() {
		t.Error("expected active account to be active")
	}
	a.Status = StatusTrialing
	if a.IsActive() {
		t.Error("expected non-active account to not be active")
	}
}
