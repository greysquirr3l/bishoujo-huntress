package huntress_test

import (
	"testing"

	"github.com/greysquirr3l/bishoujo-huntress/pkg/huntress"
)

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
