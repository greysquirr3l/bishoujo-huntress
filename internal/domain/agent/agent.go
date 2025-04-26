// Package agent contains the domain model for Huntress agents.
package agent

import (
	"errors"
	"time"
)

// Common errors for agent domain
var (
	ErrInvalidID           = errors.New("invalid agent ID")
	ErrEmptyHostname       = errors.New("hostname cannot be empty")
	ErrInvalidOrganization = errors.New("invalid organization ID")
)

// Status represents the current status of an agent
type Status string

const (
	// AgentStatusOnline indicates an online agent
	AgentStatusOnline Status = "online"
	// AgentStatusOffline indicates an offline agent
	AgentStatusOffline Status = "offline"
	// AgentStatusPending indicates a pending agent
	AgentStatusPending Status = "pending"
	// AgentStatusUnknown indicates an agent with unknown status
	AgentStatusUnknown Status = "unknown"
)

// Platform represents the platform on which an agent is running
type Platform string

const (
	// PlatformWindows indicates a Windows platform
	PlatformWindows Platform = "windows"
	// PlatformMac indicates a macOS platform
	PlatformMac Platform = "mac"
	// PlatformLinux indicates a Linux platform
	PlatformLinux Platform = "linux"
)

// Agent represents a Huntress agent installed on an endpoint
type Agent struct {
	ID               string
	OrganizationID   int
	Version          string
	Hostname         string
	IPV4Address      string
	MACAddress       string
	Platform         Platform
	OS               string
	OSVersion        string
	Status           Status
	LastSeen         time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	ExternalID       string
	UserInfo         UserInfo
	SystemInfo       SystemInfo
	EncryptionStatus string
	Tags             []string
}

// UserInfo contains user-related information for the agent
type UserInfo struct {
	Username  string
	Domain    string
	IsAdmin   bool
	LastLogon time.Time
}

// SystemInfo contains system-related information for the agent
type SystemInfo struct {
	Manufacturer  string
	Model         string
	TotalRAM      int64 // in bytes
	DiskSize      int64 // in bytes
	ProcessorInfo string
	BIOSVersion   string
}

// Validate checks if the agent has valid data
func (a *Agent) Validate() error {
	if a.ID == "" {
		return ErrInvalidID
	}
	if a.Hostname == "" {
		return ErrEmptyHostname
	}
	if a.OrganizationID <= 0 {
		return ErrInvalidOrganization
	}
	return nil
}

// IsWindows returns true if the agent is running on Windows
func (a *Agent) IsWindows() bool {
	return a.Platform == PlatformWindows
}

// IsMac returns true if the agent is running on macOS
func (a *Agent) IsMac() bool {
	return a.Platform == PlatformMac
}

// IsLinux returns true if the agent is running on Linux
func (a *Agent) IsLinux() bool {
	return a.Platform == PlatformLinux
}

// IsOnline returns true if the agent is currently online
func (a *Agent) IsOnline() bool {
	return a.Status == AgentStatusOnline
}

// IsRecentlySeen returns true if the agent was seen within the given duration
func (a *Agent) IsRecentlySeen(within time.Duration) bool {
	return time.Since(a.LastSeen) <= within
}

// HasTag checks if the agent has a specific tag
func (a *Agent) HasTag(tag string) bool {
	for _, t := range a.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// AddTag adds a tag to the agent if it doesn't already exist
func (a *Agent) AddTag(tag string) {
	if tag == "" || a.HasTag(tag) {
		return
	}
	a.Tags = append(a.Tags, tag)
}

// RemoveTag removes a tag from the agent
func (a *Agent) RemoveTag(tag string) {
	if tag == "" {
		return
	}

	for i, t := range a.Tags {
		if t == tag {
			// Remove the tag by replacing it with the last element
			// and then truncating the slice
			a.Tags[i] = a.Tags[len(a.Tags)-1]
			a.Tags = a.Tags[:len(a.Tags)-1]
			return
		}
	}
}
