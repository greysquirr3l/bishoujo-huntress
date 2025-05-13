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

- [Development Environment Setup](#linting-and-development-environment-setup)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Architecture](#architecture)
- [Documentation](#documentation)
- [Examples](#examples)
- [Security](#security)
- [Contributing](#contributing)
- [License](#license)

## Linting and Development Environment Setup

### Automated golangci-lint Installation

The Makefile and CI workflows will automatically install `golangci-lint` (v1.56.2) if it is not present in your `$PATH`.

**Manual Installation (Recommended for Consistency):**

To ensure you have the correct version and avoid issues with Homebrew or system package managers, run:

```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.56.2
```

Or, to install to `/usr/local/bin` (requires sudo):

```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.56.2
```

After installation, ensure your `$PATH` includes the install location (e.g., `$GOPATH/bin` or `/usr/local/bin`).

**Verify installation:**

```bash
golangci-lint --version
```

The Makefile target `deps` will also install `golangci-lint` if missing. This ensures that both local and CI environments use the same linter version for consistent results.

**Note:** If you previously installed `golangci-lint` via Homebrew or another package manager, you may need to remove or update it to avoid version conflicts. The official install script is preferred for reproducibility.

### Linting and Testing

To lint and test the codebase:

```bash
make lint   # Runs golangci-lint with project config
make test   # Runs all unit and integration tests
```

For full code quality checks (format, vet, lint, test):

```bash
make check
```

### Static Application Security Testing (SAST) with Semgrep

This project uses [Semgrep](https://semgrep.dev/) for SAST. The recommended version is **v1.119.0** (pinned in CI and `ossf-attest.sh`).

To run Semgrep locally:

```bash
brew install semgrep  # or pipx install semgrep==1.119.0
semgrep --config p/owasp-top-ten .
```

Semgrep is run automatically in CI and as part of the OSSF attestation script (`ossf-attest.sh`).

#### SBOM Generation

This project uses [syft](https://github.com/anchore/syft) **v1.23.1** for SBOM generation. Please use this version for reproducibility and OSSF Scorecard compliance.

Install syft v1.23.1 (recommended):

```bash
curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh -s -- -b /usr/local/bin v1.23.1
syft --version  # should print syft 1.23.1
syft . -o cyclonedx-json > sbom.json
```

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

### Working with Audit Logs

```go
// List audit logs
logs, pagination, err := client.AuditLog.List(ctx, &huntress.AuditLogListParams{
	StartTime:    nil, // or &startTime
	EndTime:      nil, // or &endTime
	Actor:        nil, // or &actor
	Action:       nil, // or &action
	ResourceType: nil, // or &resourceType
	ResourceID:   nil, // or &resourceID
	Page:         1,
	Limit:        20,
})

if err != nil {
	log.Fatalf("Failed to list audit logs: %v", err)
}

for _, entry := range logs {
	fmt.Printf("AuditLog: %s %s %s\n", entry.Timestamp, entry.Actor, entry.Action)
}

// Get a specific audit log entry
logEntry, err := client.AuditLog.Get(ctx, "auditlog-id")

if err != nil {
	log.Fatalf("Failed to get audit log: %v", err)
}

fmt.Printf("AuditLog: %+v\n", logEntry)
```

### Working with Reports

```go
// Generate a report
report, err := client.Report.Generate(ctx, &huntress.ReportGenerateOptions{
	OrganizationID: orgID,
	ReportType:     huntress.ReportTypeIncident,
	StartDate:      time.Now().AddDate(0, -1, 0),
	EndDate:        time.Now(),
})

if err != nil {
	log.Fatalf("Failed to generate report: %v", err)
}

fmt.Printf("Report ID: %s\n", report.ID)

// Get report details
reportDetails, err := client.Report.Get(ctx, report.ID)

if err != nil {
	log.Fatalf("Failed to get report: %v", err)
}

fmt.Printf("Report Details: %+v\n", reportDetails)

// Download report
reportData, err := client.Report.Download(ctx, report.ID)

if err != nil {
	log.Fatalf("Failed to download report: %v", err)
}

fmt.Printf("Report Data: %s\n", string(reportData))
```

### Working with Webhooks

```go
// Create a new webhook
webhook, err := client.Webhook.Create(ctx, &huntress.WebhookCreateOptions{
	URL:         "https://example.com/webhook",
	Description: "My Webhook",
	Events:      []string{"incident.created", "incident.updated"},
})

if err != nil {
	log.Fatalf("Failed to create webhook: %v", err)
}

fmt.Printf("Webhook ID: %s\n", webhook.ID)

// List webhooks
webhooks, err := client.Webhook.List(ctx)

if err != nil {
	log.Fatalf("Failed to list webhooks: %v", err)
}

for _, wh := range webhooks {
	fmt.Printf("Webhook: %s (ID: %s)\n", wh.Description, wh.ID)
}

// Get a specific webhook
wh, err := client.Webhook.Get(ctx, webhook.ID)

if err != nil {
	log.Fatalf("Failed to get webhook: %v", err)
}

fmt.Printf("Webhook Details: %+v\n", wh)

// Update a webhook
wh, err = client.Webhook.Update(ctx, webhook.ID, &huntress.WebhookUpdateOptions{
	Description: "Updated Webhook",
})

if err != nil {
	log.Fatalf("Failed to update webhook: %v", err)
}

fmt.Printf("Updated Webhook: %+v\n", wh)

// Delete a webhook
err = client.Webhook.Delete(ctx, webhook.ID)

if err != nil {
	log.Fatalf("Failed to delete webhook: %v", err)
}

fmt.Println("Webhook deleted successfully")
```

### Working with Billing

```go
// Get billing summary
billingSummary, err := client.Billing.GetSummary(ctx)

if err != nil {
	log.Fatalf("Failed to get billing summary: %v", err)
}

fmt.Printf("Billing Summary: %+v\n", billingSummary)

// List invoices
invoices, err := client.Billing.ListInvoices(ctx, &huntress.BillingListInvoicesOptions{
	Page:    1,
	PerPage: 10,
})

if err != nil {
	log.Fatalf("Failed to list invoices: %v", err)
}

for _, invoice := range invoices {
	fmt.Printf("Invoice: %s (ID: %s)\n", invoice.Description, invoice.ID)
}

// Get a specific invoice
invoice, err := client.Billing.GetInvoice(ctx, "invoice-id")

if err != nil {
	log.Fatalf("Failed to get invoice: %v", err)
}

fmt.Printf("Invoice Details: %+v\n", invoice)

// Get usage statistics
usageStats, err := client.Billing.GetUsageStatistics(ctx)

if err != nil {
	log.Fatalf("Failed to get usage statistics: %v", err)
}

fmt.Printf("Usage Statistics: %+v\n", usageStats)
```

### Working with Users

```go
// List users
users, err := client.User.List(ctx, &huntress.UserListOptions{
	Page:    1,
	PerPage: 10,
})

if err != nil {
	log.Fatalf("Failed to list users: %v", err)
}

for _, user := range users {
	fmt.Printf("User: %s (ID: %s)\n", user.Name, user.ID)
}

// Get a specific user
user, err := client.User.Get(ctx, "user-id")

if err != nil {
	log.Fatalf("Failed to get user: %v", err)
}

fmt.Printf("User Details: %+v\n", user)

// Create a new user
newUser, err := client.User.Create(ctx, &huntress.UserCreateOptions{
	Email:    "user@example.com",
	Name:     "New User",
	Password:  "password123",
})

if err != nil {
	log.Fatalf("Failed to create user: %v", err)
}

fmt.Printf("New User: %+v\n", newUser)

// Update a user
updatedUser, err := client.User.Update(ctx, newUser.ID, &huntress.UserUpdateOptions{
	Name: "Updated User",
})

if err != nil {
	log.Fatalf("Failed to update user: %v", err)
}

fmt.Printf("Updated User: %+v\n", updatedUser)

// Delete a user
err = client.User.Delete(ctx, newUser.ID)

if err != nil {
	log.Fatalf("Failed to delete user: %v", err)
}

fmt.Println("User deleted successfully")
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
