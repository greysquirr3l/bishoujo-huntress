// Package huntress provides a client for the Huntress API
package huntress

// AgentSettings represents configuration settings for an agent
type AgentSettings struct {
	AutoUpdate       bool              `json:"auto_update"`
	MonitorProcesses bool              `json:"monitor_processes"`
	MonitorServices  bool              `json:"monitor_services"`
	CustomSettings   map[string]string `json:"custom_settings,omitempty"`
}

// Agent types have been moved to types_consolidated.go.
// This file now only contains types specific to the agent implementation
// that aren't shared with other parts of the codebase.
