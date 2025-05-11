package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	httpClient "github.com/greysquirr3l/bishoujo-huntress/internal/infrastructure/http"
)

// Use roundTripFunc from client_test.go (do not redeclare here)

func TestOrganizationRepository_List_HTTPError(t *testing.T) {
	stdClient := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{StatusCode: 500, Body: io.NopCloser(strings.NewReader("fail"))}
		}),
	}
	hc, _ := httpClient.NewClient("http://x", "k", "s")
	hc.HTTPClient = stdClient
	repo := &OrganizationRepository{httpClient: hc, baseURL: "http://x"}
	ctx := context.Background()
	_, _, err := repo.List(ctx, nil)
	if err == nil {
		t.Error("expected error for HTTP 500")
	}
}

func TestOrganizationRepository_List_BadJSON(t *testing.T) {
	stdClient := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("notjson"))}
		}),
	}
	hc, _ := httpClient.NewClient("http://x", "k", "s")
	hc.HTTPClient = stdClient
	repo := &OrganizationRepository{httpClient: hc, baseURL: "http://x"}
	ctx := context.Background()
	_, _, err := repo.List(ctx, nil)
	if err == nil {
		t.Error("expected decode error")
	}
}

func TestOrganizationRepository_Delete_HTTPError(t *testing.T) {
	stdClient := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{StatusCode: 404, Body: io.NopCloser(strings.NewReader("not found"))}
		}),
	}
	hc, _ := httpClient.NewClient("http://x", "k", "s")
	hc.HTTPClient = stdClient
	repo := &OrganizationRepository{httpClient: hc, baseURL: "http://x"}
	ctx := context.Background()
	err := repo.Delete(ctx, "org1")
	if err == nil {
		t.Error("expected error for HTTP 404")
	}
}

func TestOrganizationRepository_Delete_Success(t *testing.T) {
	stdClient := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{StatusCode: 204, Body: io.NopCloser(strings.NewReader(""))}
		}),
	}
	hc, _ := httpClient.NewClient("http://x", "k", "s")
	hc.HTTPClient = stdClient
	repo := &OrganizationRepository{httpClient: hc, baseURL: "http://x"}
	ctx := context.Background()
	err := repo.Delete(ctx, "org1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOrganizationRepository_Get_Success(t *testing.T) {
	want := map[string]interface{}{"id": "org1", "name": "Test Org"}
	body, _ := json.Marshal(want)
	stdClient := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}
		}),
	}
	hc, _ := httpClient.NewClient("http://x", "k", "s")
	hc.HTTPClient = stdClient
	repo := &OrganizationRepository{httpClient: hc, baseURL: "http://x"}
	ctx := context.Background()
	got, err := repo.Get(ctx, "org1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.ID != "org1" || got.Name != "Test Org" {
		t.Errorf("unexpected result: %+v", got)
	}
}

func TestOrganizationRepository_Get_HTTPError(t *testing.T) {
	stdClient := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{StatusCode: 404, Body: io.NopCloser(strings.NewReader("not found"))}
		}),
	}
	hc, _ := httpClient.NewClient("http://x", "k", "s")
	hc.HTTPClient = stdClient
	repo := &OrganizationRepository{httpClient: hc, baseURL: "http://x"}
	ctx := context.Background()
	_, err := repo.Get(ctx, "org1")
	if err == nil || !strings.Contains(err.Error(), "API error") {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestOrganizationRepository_Get_BadJSON(t *testing.T) {
	stdClient := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("notjson"))}
		}),
	}
	hc, _ := httpClient.NewClient("http://x", "k", "s")
	hc.HTTPClient = stdClient
	repo := &OrganizationRepository{httpClient: hc, baseURL: "http://x"}
	ctx := context.Background()
	_, err := repo.Get(ctx, "org1")
	if err == nil || !strings.Contains(err.Error(), "parsing response body") {
		t.Errorf("expected parse error, got %v", err)
	}
}
