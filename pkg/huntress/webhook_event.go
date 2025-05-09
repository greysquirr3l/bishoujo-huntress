// Package huntress provides helpers for validating and parsing webhook event payloads.
package huntress

import (
	"encoding/json"
	"fmt"
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
		return nil, fmt.Errorf("unmarshaling webhook event: %w", err)
	}
	if evt.Type == "" {
		return nil, fmt.Errorf("missing event type in webhook payload")
	}
	return &evt, nil
}

// ValidateWebhookEvent checks if the event payload is valid.
func ValidateWebhookEvent(evt *WebhookEvent) error {
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
