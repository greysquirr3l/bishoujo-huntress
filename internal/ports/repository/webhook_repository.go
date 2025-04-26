// Package repository defines the interface for webhook persistence.
package repository

import (
	"context"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/webhook"
)

// WebhookRepository defines the interface for persisting webhooks.
type WebhookRepository interface {
	// Get retrieves a webhook by ID.
	Get(ctx context.Context, id int64) (*webhook.Webhook, error)
	// List returns all webhooks matching the filter.
	List(ctx context.Context, filter WebhookFilter) ([]*webhook.Webhook, error)
	// Create creates a new webhook.
	Create(ctx context.Context, w *webhook.Webhook) (*webhook.Webhook, error)
	// Update updates an existing webhook by ID.
	Update(ctx context.Context, id int64, params *webhook.Webhook) (*webhook.Webhook, error)
	// Delete removes a webhook by ID.
	Delete(ctx context.Context, id int64) error
}

// WebhookFilter is used to filter webhook queries.
type WebhookFilter struct {
	Enabled   *bool
	EventType string
}
