package huntress

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// Use roundTripFunc from testhelpers_test.go (do not redeclare here)

func TestAgentService_Get(t *testing.T) {
	respAgent := &Agent{ID: "a1"}
	body, _ := json.Marshal(respAgent)
	client := newTestClient(roundTripFunc(func(_ *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header:     make(http.Header),
		}
	}))
	svc := &agentService{client: &Client{httpClient: client.httpClient}}
	agent, err := svc.Get(context.Background(), "a1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if agent.ID != "a1" {
		t.Errorf("expected ID a1, got %v", agent.ID)
	}
}

func TestAgentService_Get_HTTPError(t *testing.T) {
	client := newTestClient(roundTripFunc(func(_ *http.Request) *http.Response {
		return &http.Response{StatusCode: 500, Body: io.NopCloser(bytes.NewReader([]byte("fail")))}
	}))
	svc := &agentService{client: &Client{httpClient: client.httpClient}}
	_, err := svc.Get(context.Background(), "bad")
	if err == nil {
		t.Error("expected error for HTTP 500")
	}
}

func TestAgentService_List(t *testing.T) {
	agents := []*Agent{{ID: "a1"}, {ID: "a2"}}
	body, _ := json.Marshal(agents)
	client := newTestClient(roundTripFunc(func(_ *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header:     make(http.Header),
		}
	}))
	svc := &agentService{client: &Client{httpClient: client.httpClient}}
	result, _, err := svc.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 agents, got %d", len(result))
	}
}

func TestAgentService_List_HTTPError(t *testing.T) {
	client := newTestClient(roundTripFunc(func(_ *http.Request) *http.Response {
		return &http.Response{StatusCode: 500, Body: io.NopCloser(bytes.NewReader([]byte("fail")))}
	}))
	svc := &agentService{client: &Client{httpClient: client.httpClient}}
	_, _, err := svc.List(context.Background(), nil)
	if err == nil {
		t.Error("expected error for HTTP 500")
	}
}

func TestAgentService_List_InvalidParams(t *testing.T) {
	svc := &agentService{client: &Client{httpClient: http.DefaultClient}}
	bad := &AgentListOptions{Status: "not-a-status"}
	_, _, err := svc.List(context.Background(), bad)
	if err == nil {
		t.Error("expected error for invalid params")
	}
}

func TestAgentService_GetStats(t *testing.T) {
	stats := &AgentStatistics{TotalDetections: 42, UpTime: 99.9}
	body, _ := json.Marshal(stats)
	client := newTestClient(roundTripFunc(func(_ *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header:     make(http.Header),
		}
	}))
	svc := &agentService{client: &Client{httpClient: client.httpClient}}
	result, err := svc.GetStats(context.Background(), "a1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TotalDetections != 42 || result.UpTime != 99.9 {
		t.Errorf("unexpected stats: %+v", result)
	}
}

func TestAgentService_GetStats_HTTPError(t *testing.T) {
	client := newTestClient(roundTripFunc(func(_ *http.Request) *http.Response {
		return &http.Response{StatusCode: 500, Body: io.NopCloser(bytes.NewReader([]byte("fail")))}
	}))
	svc := &agentService{client: &Client{httpClient: client.httpClient}}
	_, err := svc.GetStats(context.Background(), "bad")
	if err == nil {
		t.Error("expected error for HTTP 500")
	}
}

func TestAgentService_Update(t *testing.T) {
	respAgent := &Agent{ID: "a1"}
	body, _ := json.Marshal(respAgent)
	client := newTestClient(roundTripFunc(func(_ *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header:     make(http.Header),
		}
	}))
	svc := &agentService{client: &Client{httpClient: client.httpClient}}
	result, err := svc.Update(context.Background(), "a1", map[string]interface{}{"foo": "bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "a1" {
		t.Errorf("expected ID a1, got %v", result.ID)
	}
}

func TestAgentService_Update_HTTPError(t *testing.T) {
	client := newTestClient(roundTripFunc(func(_ *http.Request) *http.Response {
		return &http.Response{StatusCode: 500, Body: io.NopCloser(bytes.NewReader([]byte("fail")))}
	}))
	svc := &agentService{client: &Client{httpClient: client.httpClient}}
	_, err := svc.Update(context.Background(), "bad", map[string]interface{}{"foo": "bar"})
	if err == nil {
		t.Error("expected error for HTTP 500")
	}
}

func TestAgentService_Delete(t *testing.T) {
	client := newTestClient(roundTripFunc(func(_ *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 204,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
			Header:     make(http.Header),
		}
	}))
	svc := &agentService{client: &Client{httpClient: client.httpClient}}
	err := svc.Delete(context.Background(), "a1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAgentService_Delete_HTTPError(t *testing.T) {
	client := newTestClient(roundTripFunc(func(_ *http.Request) *http.Response {
		return &http.Response{StatusCode: 500, Body: io.NopCloser(bytes.NewReader([]byte("fail")))}
	}))
	svc := &agentService{client: &Client{httpClient: client.httpClient}}
	err := svc.Delete(context.Background(), "bad")
	if err == nil {
		t.Error("expected error for HTTP 500")
	}
}
