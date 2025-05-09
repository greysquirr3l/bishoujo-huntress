package incident

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestIndicatorOfCompromise_Struct(t *testing.T) {
	id := uuid.New()
	ioc := IndicatorOfCompromise{
		ID:          id,
		IncidentID:  id,
		Type:        "file",
		Value:       "malicious.exe",
		Description: "desc",
		Source:      "sensor",
		Timestamp:   time.Now(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if ioc.ID != id {
		t.Error("ID mismatch")
	}
}

func TestArtifact_Struct(t *testing.T) {
	id := uuid.New()
	a := Artifact{
		ID:          id,
		IncidentID:  id,
		Name:        "artifact",
		Type:        "exe",
		Size:        123,
		Hash:        "abc",
		Path:        "/tmp",
		Description: "desc",
		ContentHash: "def",
		StoragePath: "/store",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if a.ID != id {
		t.Error("ID mismatch")
	}
}

func TestNote_Struct(t *testing.T) {
	id := uuid.New()
	n := Note{
		ID:         id,
		IncidentID: id,
		Content:    "note",
		CreatedBy:  "user",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if n.ID != id {
		t.Error("ID mismatch")
	}
}

func TestIncident_Struct(t *testing.T) {
	id := uuid.New()
	inc := Incident{
		ID:             id,
		OrganizationID: id,
		AgentID:        "agent",
		Title:          "incident",
		Description:    "desc",
		Status:         StatusNew,
		Severity:       SeverityHigh,
		Type:           TypeMalware,
		AssignedTo:     "user",
		Reporter:       "reporter",
		Notes:          []Note{},
		IOCs:           []IndicatorOfCompromise{},
		Artifacts:      []Artifact{},
		Tags:           []string{"tag1"},
		DetectedAt:     time.Now(),
		CreatedAt:      time.Now(),
	}
	if inc.ID != id {
		t.Error("ID mismatch")
	}
}
