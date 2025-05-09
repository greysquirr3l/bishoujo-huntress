package auditlog

import (
	"testing"
	"time"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/common"
)

func TestAuditLog_Validate(t *testing.T) {
	now := time.Now()
	valid := AuditLog{
		ID:           "id",
		Timestamp:    now,
		Actor:        "actor",
		Action:       "action",
		ResourceType: "type",
		ResourceID:   "rid",
	}
	if err := valid.Validate(); err != nil {
		t.Errorf("expected valid, got %v", err)
	}

	cases := []struct {
		name string
		log  AuditLog
		err  error
	}{
		{"empty id", AuditLog{Timestamp: now, Actor: "a", Action: "b", ResourceType: "c", ResourceID: "d"}, common.ErrInvalidID},
		{"zero timestamp", AuditLog{ID: "id", Actor: "a", Action: "b", ResourceType: "c", ResourceID: "d"}, common.ErrInvalidTimestamp},
		{"empty actor", AuditLog{ID: "id", Timestamp: now, Action: "b", ResourceType: "c", ResourceID: "d"}, common.ErrEmptyActor},
		{"empty action", AuditLog{ID: "id", Timestamp: now, Actor: "a", ResourceType: "c", ResourceID: "d"}, common.ErrEmptyAction},
		{"empty resource type", AuditLog{ID: "id", Timestamp: now, Actor: "a", Action: "b", ResourceID: "d"}, common.ErrEmptyResourceType},
		{"empty resource id", AuditLog{ID: "id", Timestamp: now, Actor: "a", Action: "b", ResourceType: "c"}, common.ErrEmptyResourceID},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.log.Validate()
			if err != tc.err {
				t.Errorf("got %v, want %v", err, tc.err)
			}
		})
	}
}
