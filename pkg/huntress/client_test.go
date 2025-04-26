package huntress_test

import (
	"context"
	"testing"

	"github.com/greysquirr3l/bishoujo-huntress/pkg/huntress"
)

func TestClient_New(t *testing.T) {
	client := huntress.New(
		huntress.WithCredentials("test-key", "test-secret"),
	)
	if client == nil {
		t.Fatal("expected client to be non-nil")
	}
}

func TestAccountService_Get_NotImplemented(t *testing.T) {
	client := huntress.New(huntress.WithCredentials("test", "test"))
	_, err := client.Account.Get(context.Background())
	if err == nil {
		t.Error("expected error for Get, got nil")
	}
}
