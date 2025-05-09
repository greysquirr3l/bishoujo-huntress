package agent

import (
	"testing"
	"time"
)

func TestAgent_Validate(t *testing.T) {
	tests := []struct {
		name    string
		a       Agent
		wantErr error
	}{
		{"valid", Agent{ID: "1", Hostname: "host", OrganizationID: 1}, nil},
		{"empty id", Agent{ID: "", Hostname: "host", OrganizationID: 1}, ErrInvalidID},
		{"empty hostname", Agent{ID: "1", Hostname: "", OrganizationID: 1}, ErrEmptyHostname},
		{"invalid org", Agent{ID: "1", Hostname: "host", OrganizationID: 0}, ErrInvalidOrganization},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.a.Validate()
			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestAgent_IsWindows(t *testing.T) {
	a := Agent{Platform: PlatformWindows}
	if !a.IsWindows() {
		t.Error("expected IsWindows true for windows platform")
	}
	a.Platform = PlatformLinux
	if a.IsWindows() {
		t.Error("expected IsWindows false for non-windows platform")
	}
}

func TestAgent_IsMac(t *testing.T) {
	a := Agent{Platform: PlatformMac}
	if !a.IsMac() {
		t.Error("expected IsMac true for mac platform")
	}
	a.Platform = PlatformLinux
	if a.IsMac() {
		t.Error("expected IsMac false for non-mac platform")
	}
}

func TestAgent_IsLinux(t *testing.T) {
	a := Agent{Platform: PlatformLinux}
	if !a.IsLinux() {
		t.Error("expected IsLinux true for linux platform")
	}
	a.Platform = PlatformWindows
	if a.IsLinux() {
		t.Error("expected IsLinux false for non-linux platform")
	}
}

func TestAgent_IsOnline(t *testing.T) {
	a := Agent{Status: AgentStatusOnline}
	if !a.IsOnline() {
		t.Error("expected IsOnline true for online status")
	}
	a.Status = AgentStatusOffline
	if a.IsOnline() {
		t.Error("expected IsOnline false for non-online status")
	}
}

func TestAgent_IsRecentlySeen(t *testing.T) {
	now := time.Now()
	a := Agent{LastSeen: now.Add(-time.Minute)}
	if !a.IsRecentlySeen(2 * time.Minute) {
		t.Error("expected recently seen true")
	}
	if a.IsRecentlySeen(30 * time.Second) {
		t.Error("expected recently seen false")
	}
}

func TestAgent_HasTag(t *testing.T) {
	a := Agent{Tags: []string{"foo", "bar"}}
	if !a.HasTag("foo") {
		t.Error("expected HasTag true for present tag")
	}
	if a.HasTag("baz") {
		t.Error("expected HasTag false for absent tag")
	}
}

func TestAgent_AddTag(t *testing.T) {
	a := Agent{Tags: []string{"foo"}}
	a.AddTag("bar")
	if !a.HasTag("bar") {
		t.Error("expected AddTag to add tag")
	}
	prev := len(a.Tags)
	a.AddTag("foo") // duplicate
	if len(a.Tags) != prev {
		t.Error("expected AddTag to not add duplicate")
	}
	a.AddTag("")
	if len(a.Tags) != prev {
		t.Error("expected AddTag to not add empty tag")
	}
}

func TestAgent_RemoveTag(t *testing.T) {
	a := Agent{Tags: []string{"foo", "bar", "baz"}}
	a.RemoveTag("bar")
	if a.HasTag("bar") {
		t.Error("expected RemoveTag to remove tag")
	}
	prev := len(a.Tags)
	a.RemoveTag("notfound")
	if len(a.Tags) != prev {
		t.Error("expected RemoveTag to not change tags for missing tag")
	}
	a.RemoveTag("")
	if len(a.Tags) != prev {
		t.Error("expected RemoveTag to not change tags for empty tag")
	}
}
