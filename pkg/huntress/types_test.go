package huntress_test

import (
	"testing"

	"github.com/greysquirr3l/bishoujo-huntress/pkg/huntress"
)

func TestAgentListOptions_Validate(t *testing.T) {
	opt := &huntress.AgentListOptions{Status: "invalid"}
	err := opt.Validate()
	if err == nil {
		t.Error("expected error for invalid status")
	}
	opt.Status = huntress.AgentStatusOnline
	err = opt.Validate()
	if err != nil {
		t.Errorf("expected no error for valid status, got %v", err)
	}
}

func TestIncidentListOptions_Validate(t *testing.T) {
	opt := &huntress.IncidentListOptions{Status: "bad"}
	err := opt.Validate()
	if err == nil {
		t.Error("expected error for invalid status")
	}
	opt.Status = huntress.IncidentStatusNew
	err = opt.Validate()
	if err != nil {
		t.Errorf("expected no error for valid status, got %v", err)
	}
}
