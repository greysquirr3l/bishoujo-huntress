package webhook

import (
	"testing"
	"time"
)

func TestWebhook_Validate(t *testing.T) {
	w := &Webhook{URL: "https://foo", EventTypes: []string{"a"}, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := w.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	w.URL = ""
	if err := w.Validate(); err != ErrInvalidWebhookURL {
		t.Errorf("expected ErrInvalidWebhookURL, got %v", err)
	}
	w.URL = "https://foo"
	w.EventTypes = nil
	if err := w.Validate(); err != ErrNoEventTypes {
		t.Errorf("expected ErrNoEventTypes, got %v", err)
	}
}
