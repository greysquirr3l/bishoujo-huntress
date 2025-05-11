package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// Use roundTripFunc from client_test.go (do not redeclare here)

// errReadCloser is an io.ReadCloser that returns an error on Close.
type errReadCloser struct {
	data []byte
}

func (e *errReadCloser) Read(p []byte) (int, error) {
	if len(e.data) == 0 {
		return 0, io.EOF
	}
	n := copy(p, e.data)
	e.data = e.data[n:]
	return n, nil
}
func (e *errReadCloser) Close() error { return assertCloseErr }

var assertCloseErr = &closeError{"close error"}

type closeError struct{ msg string }

func (e *closeError) Error() string { return e.msg }

func TestDoGetWithQueryAndDecode_ErrorStatus_CloseError(t *testing.T) {
	resp := &http.Response{
		StatusCode: 404,
		Body:       &errReadCloser{data: []byte("irrelevant")},
	}
	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response { return resp }),
	}
	ctx := context.Background()
	_, err := doGetWithQueryAndDecode(ctx, client, "http://x", "", "user", "pass", nil)
	if err == nil || !strings.Contains(err.Error(), "closing response body") {
		t.Errorf("expected close error, got %v", err)
	}
}

func TestDoGetWithQueryAndDecode_BadJSON_CloseError(t *testing.T) {
	resp := &http.Response{
		StatusCode: 200,
		Body:       &errReadCloser{data: []byte("notjson")},
	}
	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response { return resp }),
	}
	ctx := context.Background()
	_, err := doGetWithQueryAndDecode(ctx, client, "http://x", "", "user", "pass", nil)
	if err == nil || !strings.Contains(err.Error(), "closing response body") {
		t.Errorf("expected close error, got %v", err)
	}
}

func TestDoPostBulkActionAndDecode_ErrorStatus_CloseError(t *testing.T) {
	resp := &http.Response{
		StatusCode: 400,
		Body:       &errReadCloser{data: []byte("irrelevant")},
	}
	client := &http.Client{
		Transport: roundTripFunc(func(_ *http.Request) *http.Response { return resp }),
	}
	ctx := context.Background()
	_, err := doPostBulkActionAndDecode(ctx, client, "http://x", "", "user", "pass", "ids", []string{"1"}, nil)
	if err == nil || !strings.Contains(err.Error(), "closing response body") {
		t.Errorf("expected close error, got %v", err)
	}
}

func TestDoGetWithQueryAndDecode_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		if _, err := w.Write([]byte(`[{"foo":"bar"},{"baz":"qux"}]`)); err != nil {
			t.Fatalf("error writing response: %v", err)
		}
	}))
	defer srv.Close()
	client := srv.Client()
	ctx := context.Background()
	out, err := doGetWithQueryAndDecode(ctx, client, srv.URL, "", "user", "pass", map[string]string{"a": "b"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(out) != 2 || out[0]["foo"] != "bar" || out[1]["baz"] != "qux" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestDoGetWithQueryAndDecode_ErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(404)
	}))
	defer srv.Close()
	client := srv.Client()
	ctx := context.Background()
	_, err := doGetWithQueryAndDecode(ctx, client, srv.URL, "", "user", "pass", nil)
	if err == nil || !strings.Contains(err.Error(), "404") {
		t.Errorf("expected 404 error, got %v", err)
	}
}

func TestDoGetWithQueryAndDecode_BadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		if _, err := w.Write([]byte("notjson")); err != nil {
			t.Fatalf("error writing response: %v", err)
		}
	}))
	defer srv.Close()
	client := srv.Client()
	ctx := context.Background()
	_, err := doGetWithQueryAndDecode(ctx, client, srv.URL, "", "user", "pass", nil)
	if err == nil || !strings.Contains(err.Error(), "decoding response") {
		t.Errorf("expected decode error, got %v", err)
	}
}

func TestDoPostBulkActionAndDecode_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		if _, err := w.Write([]byte(`{"ok":true}`)); err != nil {
			t.Fatalf("error writing response: %v", err)
		}
	}))
	defer srv.Close()
	client := srv.Client()
	ctx := context.Background()
	out, err := doPostBulkActionAndDecode(ctx, client, srv.URL, "", "user", "pass", "ids", []string{"1", "2"}, map[string]interface{}{"foo": "bar"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out["ok"] != true {
		t.Errorf("expected ok=true, got %v", out)
	}
}

func TestDoPostBulkActionAndDecode_ErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(400)
	}))
	defer srv.Close()
	client := srv.Client()
	ctx := context.Background()
	_, err := doPostBulkActionAndDecode(ctx, client, srv.URL, "", "user", "pass", "ids", []string{"1"}, nil)
	if err == nil || !strings.Contains(err.Error(), "400") {
		t.Errorf("expected 400 error, got %v", err)
	}
}

func TestDoPostBulkActionAndDecode_BadJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(200)
		if _, err := w.Write([]byte("notjson")); err != nil {
			t.Fatalf("error writing response: %v", err)
		}
	}))
	defer srv.Close()
	client := srv.Client()
	ctx := context.Background()
	_, err := doPostBulkActionAndDecode(ctx, client, srv.URL, "", "user", "pass", "ids", []string{"1"}, nil)
	if err == nil || !strings.Contains(err.Error(), "decoding response") {
		t.Errorf("expected decode error, got %v", err)
	}
}

func TestBuildQueryParams_Empty(t *testing.T) {
	v := struct{}{}
	params := buildQueryParams(v)
	if params != "" {
		t.Errorf("expected empty params, got %q", params)
	}
}

func TestBuildQueryParams_Basic(t *testing.T) {
	v := struct {
		Foo string `url:"foo"`
		Bar int    `url:"bar"`
	}{Foo: "abc", Bar: 123}
	params := buildQueryParams(v)
	u, err := url.ParseQuery(params)
	if err != nil {
		t.Fatalf("invalid query: %v", err)
	}
	if u.Get("foo") != "abc" || u.Get("bar") != "123" {
		t.Errorf("unexpected params: %v", u)
	}
}

func TestBuildQueryParams_IgnoreEmpty(t *testing.T) {
	v := struct {
		Foo string `url:"foo"`
		Bar int    `url:"bar"`
	}{Foo: "", Bar: 0}
	params := buildQueryParams(v)
	if params != "" {
		t.Errorf("expected empty params, got %q", params)
	}
}
