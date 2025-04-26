// Package api implements the WebhookRepository using the Huntress API.
package api

import (
	"context"
	"fmt"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/webhook"
	"github.com/greysquirr3l/bishoujo-huntress/internal/ports/repository"
)

// WebhookRepositoryImpl implements repository.WebhookRepository using the Huntress API.
type WebhookRepositoryImpl struct {
	// TODO: Inject HTTP client, base URL, auth, etc.
}

// Get retrieves a webhook by ID.
func (r *WebhookRepositoryImpl) Get(_ context.Context, _ int64) (*webhook.Webhook, error) {
	// TODO: Implement API call
	return nil, fmt.Errorf("not implemented")
}

// List returns all webhooks matching the filter.
func (r *WebhookRepositoryImpl) List(_ context.Context, _ repository.WebhookFilter) ([]*webhook.Webhook, error) {
	// TODO: Implement API call
	return nil, fmt.Errorf("not implemented")
}

// Create creates a new webhook.
func (r *WebhookRepositoryImpl) Create(_ context.Context, _ *webhook.Webhook) (*webhook.Webhook, error) {
	// TODO: Implement API call
	return nil, fmt.Errorf("not implemented")
}

// Update updates an existing webhook.
func (r *WebhookRepositoryImpl) Update(_ context.Context, _ int64, _ *webhook.Webhook) (*webhook.Webhook, error) {
	// TODO: Implement API call
	return nil, fmt.Errorf("not implemented")
}

// Delete removes a webhook by ID.
func (r *WebhookRepositoryImpl) Delete(_ context.Context, _ int64) error {
	// TODO: Implement API call
	return fmt.Errorf("not implemented")
}
