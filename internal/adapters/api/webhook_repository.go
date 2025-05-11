// Package api provides the Webhook API adapter for Huntress.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/webhook"
)

// WebhookRepository provides CRUD operations for Huntress webhooks.
//
// NOTE: Only one definition of WebhookRepository is allowed in this file.
type WebhookRepository struct {
	Client    *http.Client
	BaseURL   string
	APIKey    string
	APISecret string
}

// List returns all webhooks registered in Huntress.
// It returns a slice of webhook.Webhook and an error if the request fails.
func (r *WebhookRepository) List(ctx context.Context) ([]*webhook.Webhook, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", r.BaseURL+"/webhooks", nil)
	if err != nil {
		return nil, fmt.Errorf("webhook list: %w", err)
	}
	req.SetBasicAuth(r.APIKey, r.APISecret)
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("webhook list: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		errClose := resp.Body.Close()
		if errClose != nil {
			return nil, fmt.Errorf("webhook list: error closing response body: %w", errClose)
		}
		return nil, fmt.Errorf("webhook list: unexpected status: %d", resp.StatusCode)
	}
	var out []*webhook.Webhook
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		errClose := resp.Body.Close()
		if errClose != nil {
			return nil, fmt.Errorf("webhook list: error closing response body: %w", errClose)
		}
		return nil, fmt.Errorf("webhook list: decode: %w", err)
	}
	errClose := resp.Body.Close()
	if errClose != nil {
		return nil, fmt.Errorf("webhook list: error closing response body: %w", errClose)
	}
	return out, nil
}

// Create creates a new webhook in Huntress.
// It returns the created webhook.Webhook and an error if the request fails.
func (r *WebhookRepository) Create(ctx context.Context, wh *webhook.Webhook) (*webhook.Webhook, error) {
	body, err := json.Marshal(wh)
	if err != nil {
		return nil, fmt.Errorf("webhook create: marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", r.BaseURL+"/webhooks", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("webhook create: %w", err)
	}
	req.SetBasicAuth(r.APIKey, r.APISecret)
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("webhook create: %w", err)
	}
	errClose := resp.Body.Close()
	if errClose != nil {
		return nil, fmt.Errorf("webhook create: error closing response body: %w", errClose)
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("webhook create: unexpected status: %d", resp.StatusCode)
	}
	var out webhook.Webhook
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("webhook create: decode: %w", err)
	}
	return &out, nil
}

// Get retrieves a webhook by its ID from Huntress.
// It returns the webhook.Webhook and an error if the request fails.
func (r *WebhookRepository) Get(ctx context.Context, id string) (*webhook.Webhook, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", r.BaseURL+"/webhooks/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("webhook get: %w", err)
	}
	req.SetBasicAuth(r.APIKey, r.APISecret)
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("webhook get: %w", err)
	}
	errClose := resp.Body.Close()
	if errClose != nil {
		return nil, fmt.Errorf("webhook get: error closing response body: %w", errClose)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("webhook get: unexpected status: %d", resp.StatusCode)
	}
	var out webhook.Webhook
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("webhook get: decode: %w", err)
	}
	return &out, nil
}

// Update updates an existing webhook in Huntress.
// It returns the updated webhook.Webhook and an error if the request fails.
func (r *WebhookRepository) Update(ctx context.Context, id string, wh *webhook.Webhook) (*webhook.Webhook, error) {
	body, err := json.Marshal(wh)
	if err != nil {
		return nil, fmt.Errorf("webhook update: marshal: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", r.BaseURL+"/webhooks/"+id, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("webhook update: %w", err)
	}
	req.SetBasicAuth(r.APIKey, r.APISecret)
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("webhook update: %w", err)
	}
	errClose := resp.Body.Close()
	if errClose != nil {
		return nil, fmt.Errorf("webhook update: error closing response body: %w", errClose)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("webhook update: unexpected status: %d", resp.StatusCode)
	}
	var out webhook.Webhook
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("webhook update: decode: %w", err)
	}
	return &out, nil
}

// Delete removes a webhook by its ID from Huntress.
// It returns an error if the request fails.
func (r *WebhookRepository) Delete(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", r.BaseURL+"/webhooks/"+id, nil)
	if err != nil {
		return fmt.Errorf("webhook delete: %w", err)
	}
	req.SetBasicAuth(r.APIKey, r.APISecret)
	resp, err := r.Client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook delete: %w", err)
	}
	errClose := resp.Body.Close()
	if errClose != nil {
		return fmt.Errorf("webhook delete: error closing response body: %w", errClose)
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("webhook delete: unexpected status: %d", resp.StatusCode)
	}
	return nil
}

// (Removed duplicate WebhookRepository struct and List method)
