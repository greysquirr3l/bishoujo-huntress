package common

import (
	"errors"
	"testing"
)

func TestDomainError_ErrorAndUnwrap(t *testing.T) {
	base := errors.New("base")
	d := &DomainError{Code: "C", Message: "msg", Err: base}
	if d.Error() == "" {
		t.Error("expected non-empty error string")
	}
	if !errors.Is(d, base) {
		t.Error("unwrap did not work")
	}
}

func TestNewDomainError(t *testing.T) {
	err := NewDomainError("C", "msg", nil)
	if err.Code != "C" || err.Message != "msg" || err.Err != nil {
		t.Error("unexpected domain error fields")
	}
}

func TestValidationErrors(t *testing.T) {
	var errs ValidationErrors
	err := errs.Error()
	if err != "validation errors: 0 errors" {
		t.Errorf("got %q", err)
	}
	errs.AddError("f", "bad")
	if !errs.HasErrors() {
		t.Error("expected HasErrors true")
	}
	if errs.Error() == "" {
		t.Error("expected error string")
	}
}
