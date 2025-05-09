package huntress

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestReportService_Generate(t *testing.T) {
	respReport := &Report{ID: "r1"}
	body, _ := json.Marshal(respReport)
	client := newTestClient(roundTripFunc(func(r *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header:     make(http.Header),
		}
	}))
	svc := &reportService{client: &Client{httpClient: client.Client.httpClient}}
	// Pass a valid ReportGenerateInput as the body
	input := &ReportGenerateInput{}
	rep, err := svc.Generate(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rep.ID != "r1" {
		t.Errorf("expected report ID r1, got %v", rep.ID)
	}
}

func TestReportService_Generate_HTTPError(t *testing.T) {
	client := newTestClient(roundTripFunc(func(r *http.Request) *http.Response {
		return &http.Response{StatusCode: 500, Body: io.NopCloser(bytes.NewReader([]byte("fail")))}
	}))
	svc := &reportService{client: &Client{httpClient: client.Client.httpClient}}
	_, err := svc.Generate(context.Background(), &ReportGenerateInput{})
	if err == nil {
		t.Error("expected error for HTTP 500")
	}
}

// Additional tests for Get, List, Download, GetSummary, GetDetails, Schedule, Export would follow the same pattern.
