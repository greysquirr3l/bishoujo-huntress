package httpclient

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

func TestBasicAuthHeader(t *testing.T) {
	h := BasicAuthHeader("foo", "bar")
	if !strings.HasPrefix(h, "Basic ") {
		t.Errorf("expected prefix 'Basic ', got %s", h)
	}
}

func TestAddBasicAuth(t *testing.T) {
	req, _ := http.NewRequestWithContext(context.Background(), "GET", "http://example.com", nil)
	AddBasicAuth(req, "foo", "bar")
	h := req.Header.Get("Authorization")
	if !strings.HasPrefix(h, "Basic ") {
		t.Errorf("expected prefix 'Basic ', got %s", h)
	}
}
