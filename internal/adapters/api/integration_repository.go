// Package api provides the Integrations API adapter for Huntress.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// IntegrationRepository provides access to Huntress integrations.
type IntegrationRepository struct {
	Client    *http.Client
	BaseURL   string
	APIKey    string
	APISecret string
}

// Get retrieves a specific integration by ID.
func (r *IntegrationRepository) Get(ctx context.Context, id string) (map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", r.BaseURL+"/integrations/"+id, nil)
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
		return nil, fmt.Errorf("integration get failed: %d", resp.StatusCode)
	}
	var out map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

// Create creates a new integration.
func (r *IntegrationRepository) Create(ctx context.Context, integration map[string]interface{}) (map[string]interface{}, error) {
	body, err := json.Marshal(integration)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", r.BaseURL+"/integrations", bytes.NewReader(body))
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
		return nil, fmt.Errorf("integration create failed: %d", resp.StatusCode)
	}
	var out map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

// Update updates an existing integration.
func (r *IntegrationRepository) Update(ctx context.Context, id string, integration map[string]interface{}) (map[string]interface{}, error) {
	body, err := json.Marshal(integration)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "PUT", r.BaseURL+"/integrations/"+id, bytes.NewReader(body))
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
		return nil, fmt.Errorf("integration update failed: %d", resp.StatusCode)
	}
	var out map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

// Delete removes an integration by ID.
func (r *IntegrationRepository) Delete(ctx context.Context, id string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", r.BaseURL+"/integrations/"+id, nil)
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
		return fmt.Errorf("integration delete failed: %d", resp.StatusCode)
	}
	return nil
}

// List returns all integrations.
func (r *IntegrationRepository) List(ctx context.Context, params map[string]string) ([]map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", r.BaseURL+"/integrations", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()
	req.SetBasicAuth(r.APIKey, r.APISecret)
	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("integration list failed: %d", resp.StatusCode)
	}
	var out []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}
