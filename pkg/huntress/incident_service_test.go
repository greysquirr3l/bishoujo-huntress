package huntress

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestIncidentService_Get(t *testing.T) {
	respIncident := &Incident{ID: "i1"}
	body, _ := json.Marshal(respIncident)
	client := newTestClient(roundTripFunc(func(r *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header:     make(http.Header),
		}
	}))
	svc := &incidentService{client: &Client{httpClient: client.Client.httpClient}}
	inc, err := svc.Get(context.Background(), "i1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inc.ID != "i1" {
		t.Errorf("expected incident ID i1, got %v", inc.ID)
	}
}

func TestIncidentService_Get_HTTPError(t *testing.T) {
	client := newTestClient(roundTripFunc(func(r *http.Request) *http.Response {
		return &http.Response{StatusCode: 500, Body: io.NopCloser(bytes.NewReader([]byte("fail")))}
	}))
	svc := &incidentService{client: &Client{httpClient: client.Client.httpClient}}
	_, err := svc.Get(context.Background(), "bad")
	if err == nil {
		t.Error("expected error for HTTP 500")
	}
}
