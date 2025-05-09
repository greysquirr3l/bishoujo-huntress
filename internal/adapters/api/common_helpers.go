// Package api provides shared helpers for API adapters.
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// doGetWithQueryAndDecode performs a GET request with query params and decodes the JSON array response.
func doGetWithQueryAndDecode(ctx context.Context, client *http.Client, baseURL, endpoint, apiKey, apiSecret string, params map[string]string) ([]map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating GET request: %w", err)
	}
	q := req.URL.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()
	req.SetBasicAuth(apiKey, apiSecret)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing GET request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		errClose := resp.Body.Close()
		if errClose != nil {
			return nil, fmt.Errorf("GET %s failed: %d; closing response body: %w", endpoint, resp.StatusCode, errClose)
		}
		return nil, fmt.Errorf("GET %s failed: %d", endpoint, resp.StatusCode)
	}
	var out []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		errClose := resp.Body.Close()
		if errClose != nil {
			return nil, fmt.Errorf("decoding response: %w; closing response body: %w", err, errClose)
		}
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	errClose := resp.Body.Close()
	if errClose != nil {
		return nil, fmt.Errorf("closing response body: %w", errClose)
	}
	return out, nil
}

// doPostBulkActionAndDecode performs a POST to a bulk endpoint and decodes the JSON object response.
func doPostBulkActionAndDecode(ctx context.Context, client *http.Client, baseURL, endpoint, apiKey, apiSecret, idsKey string, ids []string, payload interface{}) (map[string]interface{}, error) {
	reqBody := map[string]interface{}{
		idsKey: ids,
		"data": payload,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshaling request body: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", baseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating POST request: %w", err)
	}
	req.SetBasicAuth(apiKey, apiSecret)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing POST request: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errClose := resp.Body.Close()
		if errClose != nil {
			return nil, fmt.Errorf("POST %s failed: %d; closing response body: %w", endpoint, resp.StatusCode, errClose)
		}
		return nil, fmt.Errorf("POST %s failed: %d", endpoint, resp.StatusCode)
	}
	var out map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		errClose := resp.Body.Close()
		if errClose != nil {
			return nil, fmt.Errorf("decoding response: %w; closing response body: %w", err, errClose)
		}
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	errClose := resp.Body.Close()
	if errClose != nil {
		return nil, fmt.Errorf("closing response body: %w", errClose)
	}
	return out, nil
}
