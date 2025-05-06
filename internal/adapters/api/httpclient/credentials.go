// Package httpclient provides credential helpers for Huntress API adapters.
package httpclient

import (
	"encoding/base64"
	"net/http"
)

// BasicAuthHeader returns the value for an HTTP Basic Authorization header.
func BasicAuthHeader(apiKey, apiSecret string) string {
	creds := apiKey + ":" + apiSecret
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(creds))
}

// AddBasicAuth sets the Authorization header for the request.
func AddBasicAuth(req *http.Request, apiKey, apiSecret string) {
	req.Header.Set("Authorization", BasicAuthHeader(apiKey, apiSecret))
}
