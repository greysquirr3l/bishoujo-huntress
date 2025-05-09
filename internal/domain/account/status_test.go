package account

import "testing"

func TestStatus_String(t *testing.T) {
	var s Status = "active"
	if s.String() != "active" {
		t.Errorf("expected 'active', got %q", s.String())
	}
}

func TestStatus_IsValid(t *testing.T) {
	var s Status = "active"
	if !s.IsValid() {
		t.Error("expected valid status")
	}
	s = "invalid"
	if s.IsValid() {
		t.Error("expected invalid status")
	}
}
