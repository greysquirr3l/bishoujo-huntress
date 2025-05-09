// Package webhook defines the domain entity for Huntress webhooks.
package webhook

import (
	"errors"
	"time"
)

// Webhook represents a Huntress webhook registration.
type Webhook struct {
	ID         int64     `json:"id"`
	URL        string    `json:"url"`
	EventTypes []string  `json:"event_types"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Validate checks the webhook for domain validity.
func (w *Webhook) Validate() error {
	if w.URL == "" {
		return ErrInvalidWebhookURL
	}
	if len(w.EventTypes) == 0 {
		return ErrNoEventTypes
	}
	return nil
}

// ErrInvalidWebhookURL is returned when a webhook URL is empty.
var ErrInvalidWebhookURL = errors.New("webhook URL must not be empty")

// ErrNoEventTypes is returned when a webhook has no event types.
var ErrNoEventTypes = errors.New("webhook must have at least one event type")
