# Bishoujo Huntress

<img src="docs/img/bishoujo-huntress_crop.png" alt="Bishoujo Huntress Logo" width="400">

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
