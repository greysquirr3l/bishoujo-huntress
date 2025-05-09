package webhook

import (
	"testing"
)

func TestParseEvent(t *testing.T) {
	good := []byte(`{"id":"1","type":"foo","timestamp":"2024-01-01T00:00:00Z","data":{}}`)
	evt, err := ParseEvent(good)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if evt.Type != "foo" {
		t.Error("type mismatch")
	}
	bad := []byte(`{"id":"1","timestamp":"2024-01-01T00:00:00Z","data":{}}`)
	_, err = ParseEvent(bad)
	if err == nil {
		t.Error("expected error for missing type")
	}
	invalid := []byte(`notjson`)
	_, err = ParseEvent(invalid)
	if err == nil {
		t.Error("expected error for invalid json")
	}
}

func TestValidateEvent(t *testing.T) {
	if err := ValidateEvent(nil); err == nil {
		t.Error("expected error for nil event")
	}
	evt := &Event{ID: "", Type: "foo"}
	if err := ValidateEvent(evt); err == nil {
		t.Error("expected error for missing id")
	}
	evt.ID = "1"
	evt.Type = ""
	if err := ValidateEvent(evt); err == nil {
		t.Error("expected error for missing type")
	}
	evt.Type = "foo"
	if err := ValidateEvent(evt); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
