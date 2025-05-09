// Package webhook provides helpers for validating and parsing webhook payloads.
package webhook

import (
	"encoding/json"
	"fmt"
)

// Event represents a generic Huntress webhook event payload.
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// ParseEvent parses a webhook payload into an Event struct.
func ParseEvent(payload []byte) (*Event, error) {
	var evt Event
	if err := json.Unmarshal(payload, &evt); err != nil {
		return nil, fmt.Errorf("unmarshaling webhook event: %w", err)
	}
	if evt.Type == "" {
		return nil, fmt.Errorf("missing event type in webhook payload")
	}
	return &evt, nil
}

// ValidateEvent checks if the event payload is valid.
func ValidateEvent(evt *Event) error {
	if evt == nil {
		return fmt.Errorf("event is nil")
	}
	if evt.Type == "" {
		return fmt.Errorf("event type is required")
	}
	if evt.ID == "" {
		return fmt.Errorf("event id is required")
	}
	return nil
}
