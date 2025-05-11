package api

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/webhook"
)

// Use roundTripFunc from client_test.go (do not redeclare here)

func TestWebhookRepository_List_HTTPError(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(_ *http.Request) *http.Response {
		return &http.Response{StatusCode: 500, Body: io.NopCloser(strings.NewReader("fail"))}
	})}
	repo := &WebhookRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	_, err := repo.List(ctx)
	if err == nil {
		t.Error("expected error for HTTP 500")
	}
}

func TestWebhookRepository_List_BadJSON(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(_ *http.Request) *http.Response {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("notjson"))}
	})}
	repo := &WebhookRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	_, err := repo.List(ctx)
	if err == nil {
		t.Error("expected decode error")
	}
}

func TestWebhookRepository_Delete_HTTPError(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(_ *http.Request) *http.Response {
		return &http.Response{StatusCode: 404, Body: io.NopCloser(strings.NewReader("not found"))}
	})}
	repo := &WebhookRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	err := repo.Delete(ctx, "webhook1")
	if err == nil {
		t.Error("expected error for HTTP 404")
	}
}

func TestWebhookRepository_Delete_Success(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(_ *http.Request) *http.Response {
		return &http.Response{StatusCode: 204, Body: io.NopCloser(strings.NewReader(""))}
	})}
	repo := &WebhookRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	err := repo.Delete(ctx, "webhook1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWebhookRepository_List_Success(t *testing.T) {
	// Return a valid JSON array of webhooks
	resp := `[{"id":1,"url":"https://example.com","event_types":["incident.created"],"enabled":true,"created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"}]`
	client := &http.Client{Transport: roundTripFunc(func(_ *http.Request) *http.Response {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(resp))}
	})}
	repo := &WebhookRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	webhooks, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(webhooks) != 1 || webhooks[0].URL != "https://example.com" {
		t.Errorf("unexpected result: %+v", webhooks)
	}
}

func TestWebhookRepository_Create_HTTPError(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(_ *http.Request) *http.Response {
		return &http.Response{StatusCode: 500, Body: io.NopCloser(strings.NewReader("fail"))}
	})}
	repo := &WebhookRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	wh := &webhook.Webhook{URL: "https://example.com", EventTypes: []string{"incident.created"}, Enabled: true}
	_, err := repo.Create(ctx, wh)
	if err == nil {
		t.Error("expected error for HTTP 500")
	}
}

func TestWebhookRepository_Create_BadJSON(t *testing.T) {
	client := &http.Client{Transport: roundTripFunc(func(_ *http.Request) *http.Response {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("notjson"))}
	})}
	repo := &WebhookRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	wh := &webhook.Webhook{URL: "https://example.com", EventTypes: []string{"incident.created"}, Enabled: true}
	_, err := repo.Create(ctx, wh)
	if err == nil {
		t.Error("expected decode error")
	}
}

func TestWebhookRepository_Create_Success(t *testing.T) {
	resp := `{"id":2,"url":"https://example.com","event_types":["incident.created"],"enabled":true,"created_at":"2023-01-01T00:00:00Z","updated_at":"2023-01-01T00:00:00Z"}`
	client := &http.Client{Transport: roundTripFunc(func(_ *http.Request) *http.Response {
		return &http.Response{StatusCode: 201, Body: io.NopCloser(strings.NewReader(resp))}
	})}
	repo := &WebhookRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	wh := &webhook.Webhook{URL: "https://example.com", EventTypes: []string{"incident.created"}, Enabled: true}
	got, err := repo.Create(ctx, wh)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.URL != "https://example.com" {
		t.Errorf("unexpected result: %+v", got)
	}
}
