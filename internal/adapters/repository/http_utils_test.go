package repository

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestExtractPagination(t *testing.T) {
	headers := http.Header{}
	headers.Set("X-Page", "3")
	headers.Set("X-Per-Page", "50")
	headers.Set("X-Total-Pages", "7")
	headers.Set("X-Total-Count", "123")
	p := extractPagination(headers)
	if p.Page != 3 {
		t.Errorf("expected Page=3, got %d", p.Page)
	}
	if p.PerPage != 50 {
		t.Errorf("expected PerPage=50, got %d", p.PerPage)
	}
	if p.TotalPages != 7 {
		t.Errorf("expected TotalPages=7, got %d", p.TotalPages)
	}
	if p.TotalItems != 123 {
		t.Errorf("expected TotalItems=123, got %d", p.TotalItems)
	}

	// Test defaults when headers are missing or invalid
	empty := http.Header{}
	p2 := extractPagination(empty)
	if p2.Page != 1 || p2.PerPage != 20 || p2.TotalPages != 1 || p2.TotalItems != 0 {
		t.Errorf("expected defaults, got %+v", p2)
	}

	bad := http.Header{}
	bad.Set("X-Page", "notanint")
	bad.Set("X-Per-Page", "bad")
	bad.Set("X-Total-Pages", "bad")
	bad.Set("X-Total-Count", "bad")
	p3 := extractPagination(bad)
	if p3.Page != 1 || p3.PerPage != 20 || p3.TotalPages != 1 || p3.TotalItems != 0 {
		t.Errorf("expected fallback to defaults, got %+v", p3)
	}
}

func TestParseTime(t *testing.T) {
	tstr := "2023-05-09T12:34:56Z"
	tm, err := parseTime(tstr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tm.Year() != 2023 || tm.Month() != 5 || tm.Day() != 9 {
		t.Errorf("unexpected time parsed: %v", tm)
	}

	// Empty string returns zero time
	tm2, err := parseTime("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !tm2.IsZero() {
		t.Errorf("expected zero time, got %v", tm2)
	}

	// Bad format returns error
	_, err = parseTime("not-a-time")
	if err == nil {
		t.Errorf("expected error for bad time string")
	}
}

func TestBuildQueryParams(t *testing.T) {
	filters := map[string]interface{}{
		"str":   "foo",
		"int":   42,
		"strs":  []string{"a", "b"},
		"bool":  true,
		"pstr":  func() *string { s := "bar"; return &s }(),
		"pint":  func() *int { i := 7; return &i }(),
		"pbool": func() *bool { b := false; return &b }(),
	}
	q := buildQueryParams(filters)
	if q.Get("str") != "foo" {
		t.Errorf("expected str=foo, got %s", q.Get("str"))
	}
	if q.Get("int") != "42" {
		t.Errorf("expected int=42, got %s", q.Get("int"))
	}
	if q.Get("strs") != "a" || len(q["strs"]) != 2 {
		t.Errorf("expected strs to have two values, got %v", q["strs"])
	}
	if q.Get("bool") != "true" {
		t.Errorf("expected bool=true, got %s", q.Get("bool"))
	}
	if q.Get("pstr") != "bar" {
		t.Errorf("expected pstr=bar, got %s", q.Get("pstr"))
	}
	if q.Get("pint") != "7" {
		t.Errorf("expected pint=7, got %s", q.Get("pint"))
	}
	if q.Get("pbool") != "false" {
		t.Errorf("expected pbool=false, got %s", q.Get("pbool"))
	}
}

func TestCreateRequest(t *testing.T) {
	ctx := context.Background()
	body := []byte(`{"foo":"bar"}`)
	headers := map[string]string{"X-Test": "1"}
	req, err := createRequest(ctx, http.MethodPost, "http://example.com", body, headers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != http.MethodPost {
		t.Errorf("expected POST, got %s", req.Method)
	}
	bjson := "application/json"
	if req.Header.Get("Content-Type") != bjson {
		t.Errorf("expected Content-Type=application/json")
	}
	if req.Header.Get("X-Test") != "1" {
		t.Errorf("expected X-Test=1")
	}
	if req.Header.Get("Accept") != bjson {
		t.Errorf("expected Accept=application/json")
	}
	b, _ := io.ReadAll(req.Body)
	if !bytes.Equal(b, body) {
		t.Errorf("body mismatch: got %s", string(b))
	}
}

func TestCreateRequest_NoBody(t *testing.T) {
	ctx := context.Background()
	req, err := createRequest(ctx, http.MethodGet, "http://example.com", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != http.MethodGet {
		t.Errorf("expected GET, got %s", req.Method)
	}
	if req.Header.Get("Content-Type") != "" {
		t.Errorf("expected no Content-Type header")
	}
	if req.Header.Get("Accept") != "application/json" {
		t.Errorf("expected Accept=application/json")
	}
}

func TestHandleErrorResponse_JSON(t *testing.T) {
	resp := &http.Response{
		StatusCode: 400,
		Body:       io.NopCloser(strings.NewReader(`{"code":"ERR","message":"fail","details":"bad","request_id":"req-1"}`)),
	}
	err := handleErrorResponse(resp)
	if err == nil || !strings.Contains(err.Error(), "fail") {
		t.Errorf("expected error containing 'fail', got %v", err)
	}
}

func TestHandleErrorResponse_BadJSON(t *testing.T) {
	resp := &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(strings.NewReader("not-json")),
	}
	err := handleErrorResponse(resp)
	if err == nil || !strings.Contains(err.Error(), "HTTP error: 500") {
		t.Errorf("expected generic error, got %v", err)
	}
}
