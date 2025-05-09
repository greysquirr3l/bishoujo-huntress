package errors

import (
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	err := &APIError{StatusCode: 400, Code: "CODE", Message: "msg", Details: "details"}
	if err.Error() == "" {
		t.Error("expected non-empty error string")
	}
}

func TestNewAPIError(t *testing.T) {
	err := NewAPIError(404, "NOT_FOUND", "not found", "")
	if err.StatusCode != 404 || err.Code != "NOT_FOUND" || err.Message != "not found" {
		t.Error("unexpected api error fields")
	}
}
