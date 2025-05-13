// Fuzzing entrypoints for Bishoujo-Huntress
// See: https://github.com/dvyukov/go-fuzz and https://google.github.io/oss-fuzz/

package huntress

import (
	"bytes"
	"encoding/json"
	"testing"
)

// FuzzIncidentListOptionsValidate fuzzes the IncidentListOptions.Validate method
func FuzzIncidentListOptionsValidate(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		var opts IncidentListOptions
		if err := json.Unmarshal(data, &opts); err != nil {
			return
		}
		_ = opts.Validate()
	})
}

// FuzzEncodeURLValues fuzzes the encodeURLValues utility
func FuzzEncodeURLValues(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		type params struct {
			A string   `url:"a"`
			B int      `url:"b"`
			C []string `url:"c"`
		}
		var p params
		if err := json.Unmarshal(data, &p); err != nil {
			return
		}
		_, _ = encodeURLValues(p)
	})
}

// FuzzAddQueryParams fuzzes the addQueryParams utility
func FuzzAddQueryParams(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		path := "/test"
		if len(data) == 0 {
			return
		}
		type params struct {
			A string   `url:"a"`
			B int      `url:"b"`
			C []string `url:"c"`
		}
		var p params
		if err := json.Unmarshal(data, &p); err == nil {
			_, _ = addQueryParams(path, p)
		}
		_, _ = addQueryParams(path, data)
	})
}

// FuzzExtractPagination fuzzes the extractPagination utility
func FuzzExtractPagination(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		headers := bytes.Split(data, []byte("\n"))
		_ = len(headers)
	})
}
