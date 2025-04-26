# Bishoujo Huntress

<img src="docs/img/bishoujo-huntress_crop.png" alt="Bishoujo Huntress Logo" width="600">

[![Go Reference](https://pkg.go.dev/badge/github.com/greysquirr3l/bishoujo-huntress.svg)](https://pkg.go.dev/github.com/greysquirr3l/bishoujo-huntress)
[![Go Report Card](https://goreportcard.com/badge/github.com/greysquirr3l/bishoujo-huntress)](https://goreportcard.com/report/github.com/greysquirr3l/bishoujo-huntress)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/greysquirr3l/bishoujo-huntress/badge)](https://securityscorecards.dev/viewer/?uri=github.com/greysquirr3l/bishoujo-huntress)

A comprehensive Go client library for the Huntress API, designed with Domain-Driven Design and Clean Architecture principles.

## 🚀 Features

- Complete coverage of Huntress API endpoints
- Domain-driven design with clear separation of concerns
- Thread-safe operations for concurrent use
- Robust error handling and reporting
- Rate limiting compliance (60 requests/minute)
- Comprehensive documentation and examples
- Minimal external dependencies

## 📋 Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Architecture](#architecture)
- [Documentation](#documentation)
- [Examples](#examples)
- [Security](#security)
- [Contributing](#contributing)
- [License](#license)

## 📥 Installation

```bash
go get github.com/greysquirr3l/bishoujo-huntress
```

Ensure you're using Go 1.20 or later for best compatibility.

## 🚀 Quick Start

### Client Initialization

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/greysquirr3l/bishoujo-huntress/pkg/huntress"
	"github.com/greysquirr3l/bishoujo-huntress/pkg/organization"
)

func main() {
	// Initialize client with API credentials
	client := huntress.New(
		huntress.WithCredentials("YOUR_API_KEY", "YOUR_API_SECRET"),
		huntress.WithTimeout(30 * time.Second),
	)

	ctx := context.Background()

	// List organizations
	orgs, err := client.Organization.List(ctx, &organization.ListParams{
		Page:  1,
		Limit: 10,
	})

	if err != nil {
		log.Fatalf("Failed to list organizations: %v", err)
	}

	// Process results
	for _, org := range orgs {
		fmt.Printf("Organization: %s (ID: %d)\n", org.Name, org.ID)
	}
}
```

### Working with Agents

```go
// Get agents for a specific organization
agents, err := client.Agent.List(ctx, &agent.ListParams{
	OrganizationID: orgID,
	Page:           1,
	Limit:          100,
})

if err != nil {
	log.Fatalf("Failed to list agents: %v", err)
}

// Get a specific agent by ID
agentDetails, err := client.Agent.Get(ctx, "agent-id-here")
```

### Handling Incidents

```go
// Get all incident reports
incidents, err := client.Incident.List(ctx, &incident.ListParams{
	Status:     incident.StatusOpen,
	StartDate:  time.Now().AddDate(0, -1, 0),  // Last month
	EndDate:    time.Now(),
	Page:       1,
	Limit:      25,
})
```

## 🏛 Architecture

Bishoujo-Huntress is built following Domain-Driven Design (DDD), Command Query Responsibility Segregation (CQRS), and Clean Architecture principles. This ensures:

- Clear separation of concerns
- Testable and maintainable code
- Business-focused domain model
- Independence from external frameworks

The architecture consists of four main layers:

1. **Domain Layer**: Core business entities and logic
2. **Application Layer**: Use cases implemented as commands and queries
3. **Infrastructure Layer**: External concerns like HTTP communication
4. **Interface Layer**: Public API for end users

For more details, see the [Architecture Documentation](docs/architecture.md).

## 📚 Documentation

### API Reference

Full API documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/greysquirr3l/bishoujo-huntress).

### Project Documentation

- [Architecture Design](docs/architecture.md)
- [Project Structure](docs/project_structure.md)
- [Project Guidelines](docs/project_guidelines.md)
- [Security Baselines](docs/OSSF_SECURITY_BASELINES.md)

## 🧑‍💻 Example Usage

A full example is available in [`cmd/examples/basic/main.go`](cmd/examples/basic/main.go):

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/greysquirr3l/bishoujo-huntress/pkg/huntress"
)

func main() {
	apiKey := os.Getenv("HUNTRESS_API_KEY")
	apiSecret := os.Getenv("HUNTRESS_API_SECRET")
	baseURL := os.Getenv("HUNTRESS_BASE_URL") // Optional

	if apiKey == "" || apiSecret == "" {
		log.Fatal("HUNTRESS_API_KEY and HUNTRESS_API_SECRET environment variables must be set")
	}

	client := huntress.New(
		huntress.WithCredentials(apiKey, apiSecret),
		huntress.WithTimeout(60*time.Second),
	)
	if baseURL != "" {
		client = huntress.New(
			huntress.WithCredentials(apiKey, apiSecret),
			huntress.WithTimeout(60*time.Second),
			huntress.WithBaseURL(baseURL),
		)
	}

	ctx := context.Background()

	// Get current account details
	fmt.Println("Fetching account details...")
	account, err := client.Account.Get(ctx)
	if err != nil {
		log.Fatalf("Error fetching account details: %v", err)
	}
	fmt.Printf("Account Name: %s (ID: %s)\n", account.Name, account.ID)

	// List organizations
	fmt.Println("\nFetching organizations...")
	orgs, pagination, err := client.Organization.List(ctx, &huntress.OrganizationListOptions{
		ListParams: huntress.ListParams{
			Page:     1,
			PerPage:  10,
			SortBy:   "name",
			SortDesc: false,
		},
	})
	if err != nil {
		log.Fatalf("Error listing organizations: %v", err)
	}

	fmt.Printf("Found %d organizations (page %d of %d, total items: %d):\n",
		len(orgs), pagination.CurrentPage, pagination.TotalPages, pagination.TotalItems)

	for i, org := range orgs {
		fmt.Printf("%d. %s (ID: %s, Status: %s)\n", i+1, org.Name, org.ID, org.Status)

		// For the first organization, get its agents and incidents
		if i == 0 && len(orgs) > 0 {
			fmt.Printf("\nFetching agents for organization '%s'...\n", org.Name)

			orgIDInt, err := strconv.Atoi(org.ID)
			if err != nil {
				fmt.Printf("Error converting organization ID to int: %v\n", err)
				continue
			}

			agents, _, err := client.Agent.List(ctx, &huntress.AgentListOptions{
				ListParams: huntress.ListParams{
					Page:    1,
					PerPage: 5,
				},
				OrganizationID: orgIDInt,
			})
			if err != nil {
				fmt.Printf("Error listing agents: %v\n", err)
				continue
			}

			fmt.Printf("Found %d agents:\n", len(agents))
			for j, agent := range agents {
				fmt.Printf("  %d. %s (ID: %s, OS: %s, Status: %s)\n",
					j+1, agent.Hostname, agent.ID, agent.OS, agent.Status)
			}

			fmt.Printf("\nFetching incidents for organization '%s'...\n", org.Name)

			incidents, _, err := client.Incident.List(ctx, &huntress.IncidentListOptions{
				ListOptions: huntress.ListOptions{
					Page:    1,
					PerPage: 5,
				},
				Organization: orgIDInt,
			})
			if err != nil {
				fmt.Printf("Error listing incidents: %v\n", err)
				continue
			}

			fmt.Printf("Found %d incidents:\n", len(incidents))
			for j, incident := range incidents {
				fmt.Printf("  %d. %s (ID: %s, Status: %s, Severity: %s)\n",
					j+1, incident.Title, incident.ID, incident.Status, incident.Severity)
			}
		}
	}

	// Get account statistics
	fmt.Println("\nFetching account statistics...")
	stats, err := client.Account.GetStats(ctx)
	if err != nil {
		log.Fatalf("Error fetching account statistics: %v", err)
	}

	fmt.Printf("Account statistics:\n")
	fmt.Printf("  Organization Count: %d\n", stats.OrganizationCount)
	fmt.Printf("  Agent Count: %d\n", stats.AgentCount)
	fmt.Printf("  Incident Count: %d\n", stats.IncidentCount)
	fmt.Printf("  User Count: %d\n", stats.UserCount)
}
```

## 🧩 API Coverage

This client covers all major Huntress API resources:

- **Accounts**: Get, update, list users, statistics
- **Organizations**: CRUD, list, manage users
- **Agents**: Get, list (with filters), update, delete, statistics
- **Incidents**: Get, list (with filters), update status, assign
- **Reports**: Generate, get, list, download, export, schedule
- **Billing**: Get summary, list/get invoices, usage statistics
- **Webhooks**: CRUD (scaffolded, see docs)

See [docs/todo.md](docs/todo.md) for implementation status and roadmap.

## 🛡️ Error Handling

- All API errors are mapped to Go error types with context.
- Domain, application, and infrastructure errors are separated for clarity.
- See [`pkg/huntress/errors.go`](pkg/huntress/errors.go) for details.

## 🧪 Testing

- Unit and integration tests are provided for all major services.
- Run `make test` to execute the full test suite.
- Example usage and test fixtures are in [`cmd/examples`](cmd/examples) and [`test/fixtures`](test/fixtures).

## 🧪 Examples

Working examples can be found in the [cmd/examples](cmd/examples) directory, demonstrating:

- Client initialization
- Authentication
- Working with organizations
- Managing agents
- Handling incident reports
- Generating and working with reports

## 🔒 Security

This project follows the Open Source Security Foundation (OSSF) security baselines:

- No hardcoded credentials
- Input validation for all API parameters
- Secure error handling without leaking sensitive information
- Minimal external dependencies

For details on our security practices, see the [Security Baselines](docs/OSSF_SECURITY_BASELINES.md) documentation.

To report a security vulnerability, please see our [Security Policy](SECURITY.md).

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Setup

1. Clone the repository

   ```bash
   git clone https://github.com/greysquirr3l/bishoujo-huntress.git
   cd bishoujo-huntress
   ```

2. Install dependencies

   ```bash
   go mod download
   ```

3. Run tests

   ```bash
   make test
   ```

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgements

- The [Huntress Team](https://www.huntress.com) for their API
- [Go](https://golang.org) and its contributors
- The Open Source Security Foundation for security best practices
