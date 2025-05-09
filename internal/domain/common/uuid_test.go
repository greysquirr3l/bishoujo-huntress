package common

import (
	"testing"
)

func TestUUID_Basics(t *testing.T) {
	u := NewUUID()
	if u.IsZero() {
		t.Error("new uuid should not be zero")
	}
	s := u.String()
	if s == "" {
		t.Error("string should not be empty")
	}
	parsed, err := ParseUUID(s)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if !u.Equal(parsed) {
		t.Error("parsed uuid not equal to original")
	}
	zero := UUID{}
	if !zero.IsZero() {
		t.Error("zero uuid should be zero")
	}
}
