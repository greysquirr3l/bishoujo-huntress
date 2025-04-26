// Package repository defines repository interfaces for Huntress domain entities.
package repository

import (
	"context"
	"time"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/agent"
)

// AgentFilters defines filters for agent queries
type AgentFilters struct {
	OrganizationID *int
	Search         string
	Status         *agent.Status
	Platform       *agent.Platform
	Hostname       string
	IPAddress      string
	Version        string
	LastSeenSince  *time.Time
	Tags           []string
	Page           int
	Limit          int
	OrderBy        []OrderBy
}

// AgentRepository defines the repository interface for agents
type AgentRepository interface {
	// Get retrieves an agent by ID
	Get(ctx context.Context, id string) (*agent.Agent, error)

	// List retrieves multiple agents based on filters
	List(ctx context.Context, filters AgentFilters) ([]*agent.Agent, Pagination, error)

	// Update updates an existing agent
	Update(ctx context.Context, agent *agent.Agent) error

	// Delete deletes an agent by ID
	Delete(ctx context.Context, id string) error

	// GetByHostname retrieves an agent by hostname within an organization
	GetByHostname(ctx context.Context, organizationID int, hostname string) (*agent.Agent, error)

	// GetStatistics retrieves agent statistics
	GetStatistics(ctx context.Context, id string) (map[string]interface{}, error)

	// UpdateTags updates the tags for an agent
	UpdateTags(ctx context.Context, id string, tags []string) error

	// UpdateStatus updates the status of an agent
	UpdateStatus(ctx context.Context, id string, status agent.Status) error

	// ListByOrganization retrieves all agents for an organization
	ListByOrganization(ctx context.Context, organizationID int, filters AgentFilters) ([]*agent.Agent, Pagination, error)
}
