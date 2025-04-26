package huntress_test

import (
	"context"
	"errors"
	"testing"

	"github.com/greysquirr3l/bishoujo-huntress/pkg/huntress"
)

func TestWebhookServiceStubbedMethods(t *testing.T) {
	client := huntress.New(huntress.WithCredentials("test", "test"))
	wh := client.Webhook
	_, err := wh.Get(context.Background(), 1)
	if !errors.Is(err, huntress.ErrNotImplemented) {
		t.Error("expected ErrNotImplemented for Get")
	}
	_, err = wh.List(context.Background(), nil)
	if !errors.Is(err, huntress.ErrNotImplemented) {
		t.Error("expected ErrNotImplemented for List")
	}
	_, err = wh.Create(context.Background(), nil)
	if !errors.Is(err, huntress.ErrNotImplemented) {
		t.Error("expected ErrNotImplemented for Create")
	}
	_, err = wh.Update(context.Background(), 1, nil)
	if !errors.Is(err, huntress.ErrNotImplemented) {
		t.Error("expected ErrNotImplemented for Update")
	}
	err = wh.Delete(context.Background(), 1)
	if !errors.Is(err, huntress.ErrNotImplemented) {
		t.Error("expected ErrNotImplemented for Delete")
	}
}
