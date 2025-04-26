// Package huntress provides the public Webhook service interface.
package huntress

import (
	"context"
	"errors"
)

// WebhookListParams contains parameters for listing webhooks.
type WebhookListParams struct {
	Enabled   *bool  // If set, filters webhooks by enabled status.
	EventType string // If set, filters webhooks by event type.
}

// WebhookService defines the public interface for managing Huntress webhooks.
// It provides CRUD operations for webhooks.
type WebhookService interface {
	// Get retrieves a webhook by ID.
	Get(ctx context.Context, id int64) (*Webhook, error)
	// List returns all webhooks matching the filter.
	List(ctx context.Context, params *WebhookListParams) ([]*Webhook, error)
	// Create creates a new webhook.
	Create(ctx context.Context, w *Webhook) (*Webhook, error)
	// Update updates an existing webhook.
	Update(ctx context.Context, id int64, params *WebhookUpdateParams) (*Webhook, error)
	// Delete removes a webhook by ID.
	Delete(ctx context.Context, id int64) error
}

// webhookService implements the WebhookService interface.
type webhookService struct {
	client *Client
}

// NewWebhookService returns a new WebhookService instance.
func NewWebhookService(client *Client) WebhookService {
	return &webhookService{client: client}
}

// Get returns a webhook by ID
func (s *webhookService) Get(_ context.Context, _ int64) (*Webhook, error) {
	// TODO: Implement API call
	return nil, ErrNotImplemented
}

// List returns all webhooks
func (s *webhookService) List(_ context.Context, _ *WebhookListParams) ([]*Webhook, error) {
	// TODO: Implement API call
	return nil, ErrNotImplemented
}

// Create creates a new webhook
func (s *webhookService) Create(_ context.Context, _ *Webhook) (*Webhook, error) {
	// TODO: Implement API call
	return nil, ErrNotImplemented
}

// Update updates a webhook
func (s *webhookService) Update(_ context.Context, _ int64, _ *WebhookUpdateParams) (*Webhook, error) {
	// TODO: Implement API call
	return nil, ErrNotImplemented
}

// Delete removes a webhook
func (s *webhookService) Delete(_ context.Context, _ int64) error {
	// TODO: Implement API call
	return ErrNotImplemented
}

// ErrNotImplemented is returned for stubbed methods.
var ErrNotImplemented = errors.New("not implemented")
