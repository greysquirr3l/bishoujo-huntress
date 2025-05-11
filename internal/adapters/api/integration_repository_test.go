package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

// Use roundTripFunc from client_test.go (do not redeclare here)

func TestIntegrationRepository_Update_Success(t *testing.T) {
	want := map[string]interface{}{"id": "123", "name": "updated"}
	body, _ := json.Marshal(want)
	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}
		}),
	}
	repo := &IntegrationRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	got, err := repo.Update(ctx, "123", map[string]interface{}{"name": "updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["id"] != want["id"] || got["name"] != want["name"] {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestIntegrationRepository_Update_HTTPError(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{StatusCode: 400, Body: io.NopCloser(strings.NewReader("bad req"))}
		}),
	}
	repo := &IntegrationRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	_, err := repo.Update(ctx, "123", map[string]interface{}{"name": "fail"})
	if err == nil || !strings.Contains(err.Error(), "unexpected status") {
		t.Errorf("expected status error, got %v", err)
	}
}

func TestIntegrationRepository_Update_BadJSON(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("notjson"))}
		}),
	}
	repo := &IntegrationRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	_, err := repo.Update(ctx, "123", map[string]interface{}{"name": "fail"})
	if err == nil || !strings.Contains(err.Error(), "decode") {
		t.Errorf("expected decode error, got %v", err)
	}
}

func TestIntegrationRepository_Delete_Success(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{StatusCode: 204, Body: io.NopCloser(strings.NewReader(""))}
		}),
	}
	repo := &IntegrationRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	err := repo.Delete(ctx, "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIntegrationRepository_Delete_HTTPError(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{StatusCode: 404, Body: io.NopCloser(strings.NewReader("not found"))}
		}),
	}
	repo := &IntegrationRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	err := repo.Delete(ctx, "123")
	if err == nil || !strings.Contains(err.Error(), "unexpected status") {
		t.Errorf("expected status error, got %v", err)
	}
}

func TestIntegrationRepository_Delete_CloseError(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{StatusCode: 204, Body: &errReadCloser{data: []byte("")}}
		}),
	}
	repo := &IntegrationRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	err := repo.Delete(ctx, "123")
	if err == nil || !strings.Contains(err.Error(), "closing response body") {
		t.Errorf("expected close error, got %v", err)
	}
}

func TestIntegrationRepository_Delete_WrongStatus(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(""))}
		}),
	}
	repo := &IntegrationRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	err := repo.Delete(ctx, "123")
	if err == nil || !strings.Contains(err.Error(), "unexpected status") {
		t.Errorf("expected status error, got %v", err)
	}
}

func TestIntegrationRepository_List_Success(t *testing.T) {
	want := []map[string]interface{}{{"id": "1"}, {"id": "2"}}
	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			body, _ := json.Marshal(want)
			return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(body))}
		}),
	}
	repo := &IntegrationRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	got, err := repo.List(ctx, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0]["id"] != "1" || got[1]["id"] != "2" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestIntegrationRepository_List_Error(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{StatusCode: 500, Body: io.NopCloser(strings.NewReader("fail"))}
		}),
	}
	repo := &IntegrationRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	_, err := repo.List(ctx, nil)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

// Use roundTripFunc from client_test.go for http.Client mocking

func TestIntegrationRepository_Get_Success(t *testing.T) {
	want := map[string]interface{}{"id": "123", "name": "test"}
	body, _ := json.Marshal(want)
	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}
		}),
	}
	repo := &IntegrationRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	got, err := repo.Get(ctx, "123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["id"] != want["id"] || got["name"] != want["name"] {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestIntegrationRepository_Get_HTTPError(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{StatusCode: 404, Body: io.NopCloser(strings.NewReader("not found"))}
		}),
	}
	repo := &IntegrationRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	_, err := repo.Get(ctx, "123")
	if err == nil || !strings.Contains(err.Error(), "unexpected status") {
		t.Errorf("expected status error, got %v", err)
	}
}

func TestIntegrationRepository_Get_BadJSON(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("notjson"))}
		}),
	}
	repo := &IntegrationRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	_, err := repo.Get(ctx, "123")
	if err == nil || !strings.Contains(err.Error(), "decode") {
		t.Errorf("expected decode error, got %v", err)
	}
}

func TestIntegrationRepository_Create_Success(t *testing.T) {
	want := map[string]interface{}{"id": "123", "name": "test"}
	body, _ := json.Marshal(want)
	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{
				StatusCode: 201,
				Body:       io.NopCloser(bytes.NewReader(body)),
			}
		}),
	}
	repo := &IntegrationRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	got, err := repo.Create(ctx, map[string]interface{}{"name": "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["id"] != want["id"] || got["name"] != want["name"] {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestIntegrationRepository_Create_HTTPError(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{StatusCode: 400, Body: io.NopCloser(strings.NewReader("bad req"))}
		}),
	}
	repo := &IntegrationRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	_, err := repo.Create(ctx, map[string]interface{}{"name": "test"})
	if err == nil || !strings.Contains(err.Error(), "unexpected status") {
		t.Errorf("expected status error, got %v", err)
	}
}

func TestIntegrationRepository_Create_BadJSON(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response {
			return &http.Response{StatusCode: 201, Body: io.NopCloser(strings.NewReader("notjson"))}
		}),
	}
	repo := &IntegrationRepository{Client: client, BaseURL: "http://x", APIKey: "k", APISecret: "s"}
	ctx := context.Background()
	_, err := repo.Create(ctx, map[string]interface{}{"name": "test"})
	if err == nil || !strings.Contains(err.Error(), "decode") {
		t.Errorf("expected decode error, got %v", err)
	}
}
