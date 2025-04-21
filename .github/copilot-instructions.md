# Bishoujo-Huntress: Copilot Instructions

## Project Overview

Bishoujo-Huntress is a Go client library for the Huntress API, implementing Domain-Driven Design (DDD) and Clean Architecture principles. This library provides a comprehensive, type-safe interface for interacting with all Huntress API endpoints.

Do not create Jupyter Notebooks. This is a Go library, and all code should be written in Go. The library is designed to be used in a variety of applications, including web services, command-line tools, and other Go applications.
The library is structured to support both command and query operations, following the Command Query Responsibility Segregation (CQRS) pattern. The goal is to provide a clean, maintainable, and extensible codebase that adheres to best practices in Go development.

## Key Resources

- **Project Documentation**

  - [Architecture Design](../docs/architecture.md)
  - [Project Structure](../docs/project_structure.md)
  - [Project Guidelines](../docs/project_guidelines.md)
  - [OSSF Security Baselines](../docs/OSSF_SECURITY_BASELINES.md)
  - [Implementation Todo List](../docs/todo.md) (tracks implementation progress)

- **API Resources**
  - Account: Represents a Huntress account
  - Organizations: Customer organizations within an account
  - Agents: Installed Huntress agents on endpoints
  - Incident Reports: Security incidents detected by Huntress
  - Summary Reports: Generated reports for accounts/organizations
  - Billing Reports: Billing information for accounts

## Development Guidelines

### Code Structure

1. **Domain Entities**: Define core entities in `internal/domain/{resource}` packages

   - Follow value object patterns for immutable data
   - Implement domain logic within entities
   - Use proper validation and error handling

2. **Repository Interfaces**: Define in `internal/ports/repository`

   - One repository interface per domain entity
   - Methods should reflect domain operations

3. **API Implementation**: Implement in `internal/adapters/api`

   - Handle HTTP communication details
   - Convert between API DTOs and domain entities
   - Manage rate limiting (60 requests/minute)

4. **Public API**: Expose in `pkg/huntress` and resource-specific packages
   - Simple, consistent interfaces
   - Comprehensive documentation
   - Hide implementation details

### Error Handling

1. Define domain-specific errors in domain packages
2. Wrap lower-level errors with context
3. Create public error types in `pkg/huntress/errors.go`
4. Use sentinel errors for common failure modes

### Authentication

1. Implement Basic Authentication as required by Huntress API
2. Securely handle API credentials
3. No hardcoded credentials in code

### Testing

1. Unit test all domain logic
2. Mock external dependencies for testing
3. Implement integration tests for API communication

### Security

1. Follow [OSSF Security Baselines](../docs/OSSF_SECURITY_BASELINES.md)
2. Properly handle credentials and sensitive data
3. Validate all inputs and outputs
4. Use secure dependencies
5. Implement rate limiting to prevent API abuse
6. Secure error messages to avoid information disclosure

## API Structure

### Client Configuration

```go
client := huntress.New(
    huntress.WithCredentials("API_KEY", "API_SECRET"),
    huntress.WithTimeout(30 * time.Second),
)
```

### Resource Operations

Each resource follows consistent patterns:

```go
// Get a single entity by ID
entity, err := client.Resource.Get(ctx, id)

// List entities with optional filters
entities, pagination, err := client.Resource.List(ctx, &resource.ListParams{
    Page: 1,
    Limit: 100,
    // Resource-specific filters
})

// Other resource-specific operations
result, err := client.Resource.CustomOperation(ctx, params)
```

## Code Style

1. Follow Go best practices and idioms
2. Use context for cancellation and timeouts
3. Prefer composition over inheritance
4. Keep functions focused and small
5. Comprehensive comments on all exported symbols
6. Use descriptive variable and function names

## Implementation Status and Priorities

The implementation status is tracked in detail in the [Todo List](../docs/todo.md). Below is a summary:

### Completed

- ✅ Project structure and documentation setup
- ✅ Basic domain entity structure for Account, Organization, Agent, and Incident
- ✅ Repository interface definitions for primary entities
- ✅ Basic application query structure
- ✅ Basic HTTP client and repository scaffolding
- ✅ Client interface scaffold and basic options

### In Progress

- 🔄 Domain entity implementation (need to complete properties and validation)
- 🔄 Repository implementations
- 🔄 Query handlers implementation
- 🔄 HTTP client implementation (authentication, rate limiting, etc.)

### Next Priorities

1. Complete HTTP client implementation (authentication, rate limiting, error handling, pagination)
2. Finish core domain entities with validation
3. Implement complete repository implementations for Account and Organization
4. Complete command and query handlers
5. Build out public API services
6. Implement test suite
7. Add security features and error handling
8. Complete Billing and Reports resources

## Current Development Focus

Based on the project's current state as tracked in [todo.md](../docs/todo.md), focus should be on:

1. **HTTP Client Implementation**

   - Authentication mechanism implementation
   - Rate limiting (60 requests/minute)
   - Robust error handling with proper context
   - Pagination handling for list endpoints

2. **Domain Entities Completion**

   - Adding all required properties to existing entities
   - Implementing validation logic
   - Creating value objects for complex attributes
   - Completing entity relationships

3. **Query and Command Handlers**
   - Implementing CQRS pattern completely
   - Building out query handlers first (higher priority)
   - Following with command handlers for state-changing operations

## Best Practices for Current Stage

At the current implementation stage, follow these practices:

1. **Domain First Development**

   - Complete domain models before repository implementations
   - Ensure domain validation logic is thorough
   - Keep domain entities free from infrastructure concerns

2. **Interface Stability**

   - Public interfaces in `pkg/` should remain stable
   - Internal implementations can change as needed
   - Document breaking changes if unavoidable

3. **Test Coverage**

   - Add tests as you implement features
   - Focus on testing domain logic thoroughly
   - Use mocks for external dependencies
   - Consider writing tests before implementation (TDD)

4. **Implementation Tracking**
   - Update [todo.md](../docs/todo.md) as items are completed
   - Maintain up-to-date implementation status
   - Document technical decisions that deviate from initial architecture
   - Ensure all changes are reflected in the project documentation

## Additional Resources

- [Huntress API Documentation](https://api.huntress.io/docs#api-overview)
- [Huntress API Preview (swagger)](https://api.huntress.io/docs/preview)
