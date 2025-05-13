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

func TestE2E_AgentService_List_Concurrent(t *testing.T) {
	agents := []*huntress.Agent{{ID: "a1"}, {ID: "a2"}}
	body, _ := json.Marshal(agents)
	handler := func(_ *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header:     make(http.Header),
		}
	}
	client := huntress.New(
		huntress.WithCredentials("e2e", "e2e"),
		huntress.WithHTTPClient(&http.Client{Transport: roundTripFunc(handler)}),
	)
	var wg sync.WaitGroup
	const n = 5
	wg.Add(n)
	errs := make([]error, n)
	for i := 0; i < n; i++ {
		go func(idx int) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			_, _, err := client.Agent.List(ctx, nil)
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

func TestE2E_AgentService_List_ErrorPropagation(t *testing.T) {
	handler := func(_ *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewReader([]byte(`{"message":"bad request"}`))),
			Header:     make(http.Header),
		}
	}
	client := huntress.New(
		huntress.WithCredentials("e2e", "e2e"),
		huntress.WithHTTPClient(&http.Client{Transport: roundTripFunc(handler)}),
	)
	_, _, err := client.Agent.List(context.Background(), nil)
	if err == nil || err.Error() == "" {
		t.Error("expected error for HTTP 400")
	}
}
