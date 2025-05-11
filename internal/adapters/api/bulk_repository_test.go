package api

import (
	"context"
	"net/http"
	"testing"
)

// Use roundTripFunc from client_test.go (do not redeclare here)

func TestBulkRepository_BulkAgentAction_Error(t *testing.T) {
	repo := &BulkRepository{
		Client: &http.Client{Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return nil // Simulate network error by returning nil
		})},
		BaseURL:   "https://api.example.com",
		APIKey:    "key",
		APISecret: "secret",
	}
	_, err := repo.BulkAgentAction(context.Background(), "disable", []string{"a1"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestBulkRepository_BulkAgentAction_Success(t *testing.T) {
	repo := &BulkRepository{
		Client: &http.Client{Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{
				StatusCode: 200,
				Body:       http.NoBody,
			}
		})},
		BaseURL:   "https://api.example.com",
		APIKey:    "key",
		APISecret: "secret",
	}
	// This will fail to decode, but should not error on HTTP
	_, err := repo.BulkAgentAction(context.Background(), "disable", []string{"a1"}, nil)
	if err == nil {
		t.Error("expected decode error")
	}
}

func TestBulkRepository_BulkOrgAction_Error(t *testing.T) {
	repo := &BulkRepository{
		Client: &http.Client{Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return nil // Simulate network error by returning nil
		})},
		BaseURL:   "https://api.example.com",
		APIKey:    "key",
		APISecret: "secret",
	}
	_, err := repo.BulkOrgAction(context.Background(), "disable", []string{"o1"}, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestBulkRepository_BulkOrgAction_Success(t *testing.T) {
	repo := &BulkRepository{
		Client: &http.Client{Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{
				StatusCode: 200,
				Body:       http.NoBody,
			}
		})},
		BaseURL:   "https://api.example.com",
		APIKey:    "key",
		APISecret: "secret",
	}
	// This will fail to decode, but should not error on HTTP
	_, err := repo.BulkOrgAction(context.Background(), "disable", []string{"o1"}, nil)
	if err == nil {
		t.Error("expected decode error")
	}
}
