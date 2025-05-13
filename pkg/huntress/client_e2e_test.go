package huntress_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/greysquirr3l/bishoujo-huntress/pkg/huntress"
)

// roundTripFunc is used for HTTP mocking (copied from testhelpers_test.go)
type roundTripFunc func(*http.Request) *http.Response

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r), nil
}

// e2eTestServer returns a client with a handler for E2E flows.
func e2eTestClient(handler func(*http.Request) *http.Response) *huntress.Client {
	return huntress.New(
		huntress.WithCredentials("e2e", "e2e"),
		huntress.WithHTTPClient(&http.Client{Transport: roundTripFunc(handler)}),
	)
}

func TestE2E_AccountGet_Concurrent(t *testing.T) {
	// Simulate a real account response
	account := &huntress.Account{ID: "acct-1", Name: "E2E Test"}
	body, _ := json.Marshal(account)
	handler := func(_ *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header:     make(http.Header),
		}
	}
	client := e2eTestClient(handler)
	var wg sync.WaitGroup
	const n = 10
	errs := make([]error, n)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(idx int) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			_, err := client.Account.Get(ctx)
			errs[idx] = err
		}(i)
	}
	wg.Wait()
	for i, err := range errs {
		if err != nil {
			t.Errorf("goroutine %d: unexpected error: %v", i, err)
		}
	}
}

func TestE2E_ErrorPropagation_AccountGet(t *testing.T) {
	handler := func(_ *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(bytes.NewReader([]byte(`{"message":"internal error"}`))),
			Header:     make(http.Header),
		}
	}
	client := e2eTestClient(handler)
	_, err := client.Account.Get(context.Background())
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
	if err.Error() == "" || !contains(err.Error(), "internal error") {
		t.Errorf("expected error message to propagate, got: %v", err)
	}
}

// contains is a helper for substring search.
func contains(haystack, needle string) bool {
	return len(needle) > 0 && len(haystack) > 0 && (len(haystack) >= len(needle)) && (stringIndex(haystack, needle) >= 0)
}

func stringIndex(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
