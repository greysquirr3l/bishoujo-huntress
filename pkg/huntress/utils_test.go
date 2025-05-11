package huntress

import (
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

// Additional edge case tests for utils.go

func TestEncodeURLValues_EdgeCases(t *testing.T) {
	// Test with slice of non-string type (should marshal as JSON)
	type paramsIntSlice struct {
		A []int `url:"a"`
	}
	pInt := paramsIntSlice{A: []int{1, 2, 3}}
	q, err := encodeURLValues(pInt)
	if err != nil {
		t.Fatalf("unexpected error for int slice: %v", err)
	}
	if !strings.Contains(q, "a=%5B1%2C2%2C3%5D") {
		t.Errorf("expected JSON-encoded int slice, got: %s", q)
	}

	// Test with pointer to non-string type (should marshal as JSON)
	type paramsPtr struct {
		A *int `url:"a"`
	}
	pi := 42
	pPtr := paramsPtr{A: &pi}
	q, err = encodeURLValues(pPtr)
	if err != nil {
		t.Fatalf("unexpected error for int pointer: %v", err)
	}
	if !strings.Contains(q, "a=42") && !strings.Contains(q, "a=%2242%22") {
		t.Errorf("expected int pointer encoded, got: %s", q)
	}

	// Test with custom struct field (default case)
	type custom struct{ X int }
	type paramsCustom struct {
		A custom `url:"a"`
	}
	pCustom := paramsCustom{A: custom{X: 5}}
	q, err = encodeURLValues(pCustom)
	if err != nil {
		t.Fatalf("unexpected error for custom struct: %v", err)
	}
	if !strings.Contains(q, "a=%7B%5C%22X%5C%22%3A5%7D") && !strings.Contains(q, "a=%7B%22X%22%3A5%7D") {
		t.Errorf("expected JSON-encoded custom struct, got: %s", q)
	}
}

func TestAddQueryParams(t *testing.T) {
	type params struct {
		A string `url:"a"`
		B int    `url:"b"`
		C bool   `url:"c,omitempty"`
	}
	p := params{A: "foo", B: 42, C: true}
	q, err := addQueryParams("/test", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(q, "a=foo") || !strings.Contains(q, "b=42") || !strings.Contains(q, "c=true") {
		t.Errorf("missing params in query: %s", q)
	}

	// Test with nil/empty params
	_, err = addQueryParams("/test", nil)
	if err != nil {
		t.Errorf("expected no error for nil params, got %v", err)
	}
}

func TestExtractPagination(t *testing.T) {
	r := &http.Response{Header: make(http.Header)}
	r.Header.Set("X-Page", "2")
	r.Header.Set("X-Per-Page", "10")
	r.Header.Set("X-Total-Pages", "5")
	r.Header.Set("X-Total-Items", "100")
	p := extractPagination(r)
	if p.CurrentPage != 2 || p.PerPage != 10 || p.TotalPages != 5 || p.TotalItems != 100 {
		t.Errorf("pagination parse failed: %+v", p)
	}

	// Test with nil response
	if extractPagination(nil) != nil {
		t.Error("expected nil for nil response")
	}
}

func TestParseInt(t *testing.T) {
	v, err := parseInt("123")
	if err != nil || v != 123 {
		t.Errorf("parseInt failed: %v, %d", err, v)
	}
	_, err = parseInt("notanint")
	if err == nil {
		t.Error("expected error for bad int")
	}
}

func TestEncodeURLValues(t *testing.T) {
	type params struct {
		A string   `url:"a"`
		B int      `url:"b"`
		C []string `url:"c"`
		D *string  `url:"d,omitempty"`
		E *int     `url:"e,omitempty"`
	}
	d := "bar"
	e := 7
	p := params{A: "foo", B: 42, C: []string{"x", "y"}, D: &d, E: &e}
	q, err := encodeURLValues(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("encoded query: %s", q)
	if !strings.Contains(q, "a=foo") || !strings.Contains(q, "b=42") || !strings.Contains(q, "c=x") || !strings.Contains(q, "c=y") || !strings.Contains(q, "d=%22bar%22") || !strings.Contains(q, "e=7") {
		t.Errorf("missing params in encoded: %s", q)
	}

	// Test with nil pointer
	var p2 *params
	q, err = encodeURLValues(p2)
	if err != nil || q != "" {
		t.Errorf("expected empty string for nil pointer, got %q, %v", q, err)
	}

	// Test with non-struct
	_, err = encodeURLValues(123)
	if err == nil {
		t.Error("expected error for non-struct")
	}
}

func TestAddValues_EmbeddedStruct(t *testing.T) {
	type inner struct {
		A string `url:"a"`
	}
	type outer struct {
		inner
		B int `url:"b"`
	}
	p := outer{inner: inner{A: "foo"}, B: 42}
	v := reflect.ValueOf(p)
	vals := make(url.Values)
	if err := addValues(vals, v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vals.Get("a") != "foo" || vals.Get("b") != "42" {
		t.Errorf("embedded struct values missing: %v", vals)
	}
}

func TestParseTag(t *testing.T) {
	name, opts := parseTag("foo,omitempty")
	if name != "foo" || !opts.Contains("omitempty") {
		t.Errorf("parseTag failed: %s, %v", name, opts)
	}
	name, opts = parseTag("bar")
	if name != "bar" || opts.Contains("omitempty") {
		t.Errorf("parseTag failed: %s, %v", name, opts)
	}
}

func TestTagOptions_Contains(t *testing.T) {
	o := tagOptions("a,b,c")
	if !o.Contains("b") || o.Contains("z") {
		t.Errorf("tagOptions.Contains failed")
	}
}

func TestIsEmptyValue(t *testing.T) {
	var s string
	if !isEmptyValue(reflect.ValueOf(s)) {
		t.Error("empty string should be empty")
	}
	var i int
	if !isEmptyValue(reflect.ValueOf(i)) {
		t.Error("zero int should be empty")
	}
	var b bool
	if !isEmptyValue(reflect.ValueOf(b)) {
		t.Error("false bool should be empty")
	}
	var arr []int
	if !isEmptyValue(reflect.ValueOf(arr)) {
		t.Error("empty slice should be empty")
	}
	var m map[string]int
	if !isEmptyValue(reflect.ValueOf(m)) {
		t.Error("empty map should be empty")
	}
	var p *int
	if !isEmptyValue(reflect.ValueOf(p)) {
		t.Error("nil pointer should be empty")
	}
	if isEmptyValue(reflect.ValueOf(1)) {
		t.Error("non-zero int should not be empty")
	}
}
