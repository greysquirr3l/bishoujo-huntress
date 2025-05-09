package huntress

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type roundTripFunc func(*http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

type testClient struct {
	*Client
}

func newTestClient(rt http.RoundTripper) *testClient {
	c := &Client{
		httpClient: &http.Client{Transport: rt},
		baseURL:    "http://test",
		apiKey:     "k",
		apiSecret:  "s",
		userAgent:  "test",
	}
	return &testClient{Client: c}
}

// NewRequest marshals struct bodies to JSON if needed, mimicking the real client
func (c *testClient) NewRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	url := c.baseURL + path
	var bodyReader io.Reader
	if body != nil {
		if rdr, ok := body.(io.Reader); ok {
			bodyReader = rdr
		} else {
			b, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			bodyReader = bytes.NewReader(b)
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic ZDpk") // dummy
	req.Header.Set("User-Agent", c.userAgent)
	return req, nil
}
