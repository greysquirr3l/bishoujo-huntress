# Bishoujo-Huntress: Project Structure

This document outlines the project structure for the Bishoujo-Huntress Go client library for the Huntress API. The structure follows Go best practices and implements Clean Architecture with Domain-Driven Design (DDD) principles.

## Directory Structure

```sh
bishoujo-huntress/
├── cmd/                          # Command-line tools
│   └── examples/                 # Example applications using the library
├── docs/                         # Documentation
│   ├── architecture.md           # Architecture design document
│   ├── project_structure.md      # This document
│   ├── project_guidelines.md     # DDD, CQRS, and Clean Architecture guidelines
│   ├── OSSF_SECURITY_BASELINES.md # Security baseline documentation
│   └── img/                      # Images for documentation
├── internal/                     # Private implementation packages
│   ├── domain/                   # Domain model
│   │   ├── account/              # Account domain entities
│   │   ├── organization/         # Organization domain entities
│   │   ├── agent/                # Agent domain entities
│   │   ├── incident/             # Incident report domain entities
│   │   ├── report/               # Reports domain entities
│   │   ├── billing/              # Billing domain entities
│   │   └── common/               # Shared domain objects
│   ├── application/              # Application services
│   │   ├── command/              # Command handlers
│   │   │   ├── account/          # Account commands
│   │   │   ├── organization/     # Organization commands
│   │   │   └── ...               # Other resource commands
│   │   └── query/                # Query handlers
│   │       ├── account/          # Account queries
│   │       ├── organization/     # Organization queries
│   │       └── ...               # Other resource queries
│   ├── ports/                    # Interface definitions
│   │   ├── api/                  # API client interfaces
│   │   ├── repository/           # Repository interfaces
│   │   └── service/              # Service interfaces
│   └── adapters/                 # Implementation adapters
│       ├── api/                  # API client implementations
│       ├── http/                 # HTTP client implementation
│       └── repository/           # Repository implementations
├── pkg/                          # Public packages
│   ├── huntress/                 # Main client package
│   │   ├── client.go             # Client API
│   │   ├── options.go            # Client configuration
│   │   ├── errors.go             # Error definitions
│   │   └── models.go             # Public models
│   ├── account/                  # Account service
│   ├── organization/             # Organization service
│   ├── agent/                    # Agent service
│   ├── incident/                 # Incident report service
│   ├── report/                   # Reports service
│   └── billing/                  # Billing service
├── test/                         # Integration tests
│   ├── integration/              # Integration test suites
│   └── fixtures/                 # Test fixtures
├── Makefile                      # Build and development tasks
├── go.mod                        # Go module definition
├── go.sum                        # Go module checksums
├── LICENSE                       # License file
└── README.md                     # Project README
```

## Package Details

### Public API (`pkg/huntress`)

The public API provides a clean, simple interface for consumers:

```go
// Client creation
client := huntress.New(
    huntress.WithCredentials("API_KEY", "API_SECRET"),
    huntress.WithTimeout(30 * time.Second),
)

// Working with organizations
orgs, err := client.Organization.List(ctx, &organization.ListParams{
    Page:  1,
    Limit: 100,
})
```

### Domain Layer (`internal/domain`)

The domain layer contains the core business logic and entities:

- **Entities**: Core business objects (Agent, Organization, etc.)
- **Value Objects**: Immutable objects defined by their attributes (Email, ID, etc.)
- **Domain Services**: Core business operations that don't belong to a single entity
- **Repositories**: Interfaces for data access (defined in `internal/ports/repository`)

Example domain entity:

```go
// internal/domain/agent/agent.go
package agent

type Agent struct {
    ID            string
    Version       string
    Hostname      string
    IPV4Address   string
    Platform      string
    OS            string
    OrganizationID int
    // other fields...
}

// Domain methods
func (a *Agent) IsWindows() bool {
    return a.Platform == "windows"
}
```

### Application Layer (`internal/application`)

The application layer contains use cases and orchestration logic:

- **Commands**: Operations that change state
- **Queries**: Operations that retrieve data
- **Command Handlers**: Process commands
- **Query Handlers**: Process queries

Example query:

```go
// internal/application/query/agent/list_agents.go
package agent

type ListAgentsQuery struct {
    OrganizationID *int
    Platform       *string
    Page           int
    Limit          int
}

type ListAgentsHandler struct {
    repo agent.Repository
}

func (h *ListAgentsHandler) Handle(ctx context.Context, query ListAgentsQuery) ([]*agent.Agent, error) {
    // Implementation
}
```

### Ports and Adapters (`internal/ports` and `internal/adapters`)

The ports define interfaces, and adapters implement them:

- **Ports**: Define what functionality is required
- **Adapters**: Provide implementations for ports

Example port:

```go
// internal/ports/repository/agent_repository.go
package repository

import "github.com/greysquirr3l/bishoujo-huntress/internal/domain/agent"

type AgentRepository interface {
    Get(ctx context.Context, id string) (*agent.Agent, error)
    List(ctx context.Context, filters AgentFilters) ([]*agent.Agent, Pagination, error)
    // other methods...
}
```

Example adapter:

```go
// internal/adapters/api/agent_repository.go
package api

import "github.com/greysquirr3l/bishoujo-huntress/internal/ports/repository"

type AgentRepositoryImpl struct {
    httpClient *http.Client
    baseURL    string
    auth       AuthProvider
}

func (r *AgentRepositoryImpl) Get(ctx context.Context, id string) (*agent.Agent, error) {
    // Implementation to make API call
}
```

## Testing Strategy

1. **Unit Tests**: Located alongside the code they test
   - Test domain logic
   - Test command/query handlers
   - Mock external dependencies

2. **Integration Tests**: Located in `test/integration`
   - Test against a real API or mocked server
   - Test end-to-end workflows

3. **Examples**: Located in `cmd/examples`
   - Demonstrate library usage
   - Serve as documentation

## Build and Development

The `Makefile` defines common tasks:

- `make build`: Build the library
- `make test`: Run tests
- `make lint`: Run linters
- `make coverage`: Generate test coverage report

## Dependency Management

Dependencies are managed via Go modules:

- Minimal external dependencies
- Pin versions for reproducible builds
- Regular updates for security patches

## Documentation

Documentation follows standard Go practices:

- GoDoc comments for all exported symbols
- README.md with getting started guide
- Comprehensive examples in `cmd/examples`
- Architecture and design docs in `docs/`
