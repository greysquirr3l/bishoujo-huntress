// Package huntress provides helpers for validating and parsing webhook event payloads.
package huntress

import (
	"encoding/json"
	"errors"
)

// WebhookEvent represents a generic Huntress webhook event payload.
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp string                 `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// ParseWebhookEvent parses a webhook payload into a WebhookEvent struct.
func ParseWebhookEvent(payload []byte) (*WebhookEvent, error) {
	var evt WebhookEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		return nil, err
	}
	if evt.Type == "" {
		return nil, errors.New("missing event type in webhook payload")
	}
	return &evt, nil
}

// ValidateWebhookEvent checks if the event payload is valid.
func ValidateWebhookEvent(evt *WebhookEvent) error {
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
