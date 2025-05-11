package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuditLogRepository_Get_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"id":"1","message":"test"}`)); err != nil {
			t.Fatalf("error writing response: %v", err)
		}
	}))
	defer srv.Close()
	repo := &AuditLogRepository{
		Client:    srv.Client(),
		BaseURL:   srv.URL,
		APIKey:    "key",
		APISecret: "secret",
	}
	log, err := repo.Get(context.Background(), "1")
	if err != nil || log == nil || log.ID != "1" {
		t.Fatalf("unexpected error or result: %v, %+v", err, log)
	}
}
