// Package webhook provides helpers for validating and parsing webhook payloads.
package webhook

import (
	"encoding/json"
	"errors"
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
		return nil, err
	}
	if evt.Type == "" {
		return nil, errors.New("missing event type in webhook payload")
	}
	return &evt, nil
}

// ValidateEvent checks if the event payload is valid.
func ValidateEvent(evt *Event) error {
	if evt == nil {
		return errors.New("event is nil")
	}
	if evt.Type == "" {
		return errors.New("event type is required")
	}
	if evt.ID == "" {
		return errors.New("event id is required")
	}
	return nil
}
