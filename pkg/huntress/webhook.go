// Package huntress provides the public Webhook type and params.
package huntress

import "time"

// Webhook represents a Huntress webhook registration.
type Webhook struct {
	ID         int64     `json:"id"`
	URL        string    `json:"url"`
	EventTypes []string  `json:"event_types"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// WebhookCreateParams contains parameters for creating a webhook.
type WebhookCreateParams struct {
	URL        string   `json:"url"`
	EventTypes []string `json:"event_types"`
	Enabled    bool     `json:"enabled"`
}

// WebhookUpdateParams contains parameters for updating a webhook.
type WebhookUpdateParams struct {
	// URL is the new webhook URL (optional).
	URL *string `json:"url,omitempty"`
	// EventTypes is the new set of event types (optional).
	EventTypes *[]string `json:"event_types,omitempty"`
	// Enabled indicates if the webhook should be enabled (optional).
	Enabled *bool `json:"enabled,omitempty"`
}
