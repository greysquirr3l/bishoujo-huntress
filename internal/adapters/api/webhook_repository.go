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

// (Removed duplicate struct and methods below)

func (r *WebhookRepository) List(ctx context.Context) ([]*webhook.Webhook, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", r.BaseURL+"/webhooks", nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(r.APIKey, r.APISecret)
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	var out []*webhook.Webhook
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

// Create creates a new webhook.
func (r *WebhookRepository) Create(ctx context.Context, wh *webhook.Webhook) (*webhook.Webhook, error) {
	body, err := json.Marshal(wh)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", r.BaseURL+"/webhooks", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(r.APIKey, r.APISecret)
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	var out webhook.Webhook
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Get retrieves a webhook by ID.
func (r *WebhookRepository) Get(ctx context.Context, id string) (*webhook.Webhook, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", r.BaseURL+"/webhooks/"+id, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(r.APIKey, r.APISecret)
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	var out webhook.Webhook
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update updates an existing webhook.
func (r *WebhookRepository) Update(ctx context.Context, id string, wh *webhook.Webhook) (*webhook.Webhook, error) {
	body, err := json.Marshal(wh)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", r.BaseURL+"/webhooks/"+id, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(r.APIKey, r.APISecret)
	req.Header.Set("Content-Type", "application/json")
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	var out webhook.Webhook
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Delete removes a webhook by ID.
func (r *WebhookRepository) Delete(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", r.BaseURL+"/webhooks/"+id, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(r.APIKey, r.APISecret)
	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	return nil
}

// (Removed duplicate WebhookRepository struct and List method)
