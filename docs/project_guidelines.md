# Project Guidelines: DDD, CQRS, and Clean Architecture in Go

## Introduction

This document provides guidelines for developing our app using Domain-Driven
Design (DDD), Command Query Responsibility Segregation (CQRS), and Clean Architecture
principles in Go. These approaches combine to create maintainable, testable, and
business-focused code.

## Core Concepts

### Domain-Driven Design (DDD)

DDD focuses on creating a software model that reflects the business domain:

- **Ubiquitous Language**: Use consistent terminology between developers and domain experts
- **Bounded Contexts**: Divide the domain into distinct areas with clear boundaries
- **Aggregates**: Treat related entities as a single unit with a root entity
- **Domain Events**: Model significant occurrences within the domain
- **Value Objects**: Immutable objects defined by their attributes
- **Entities**: Objects with identity that persists across state changes
- **Repositories**: Abstractions for persisting and retrieving domain objects

### CQRS (Command Query Responsibility Segregation)

CQRS separates operations that read data from operations that write data:

- **Commands**: Operations that change state but don't return data
- **Queries**: Operations that return data but don't change state
- **Command Handlers**: Process commands and update the domain model
- **Query Handlers**: Process queries and return data representations

### Clean Architecture

Clean Architecture organizes code in layers, with dependencies pointing inward:

- **Domain Layer**: Core business logic and entities
- **Application Layer**: Use cases that orchestrate domain objects
- **Infrastructure Layer**: External concerns like databases, frameworks, and I/O
- **Interface Layer**: API endpoints, UI components, and other user interfaces

## Project Structure

```sh
/cmd                  # Application entry points
  /api                # API server
  /worker             # Background workers
/internal
  /domain             # Domain model, entities, value objects, domain services
    /{boundedcontext} # Specific bounded contexts
  /application        # Application services, commands, queries, handlers
    /command          # Command definitions and handlers
    /query            # Query definitions and handlers
  /ports              # Ports (interfaces) required by the application
  /adapters           # Adapters implementing the ports
    /primary          # Input adapters (REST API, gRPC, etc.)
    /secondary        # Output adapters (repositories, external services)
  /infrastructure     # Infrastructure concerns, frameworks, DB
/pkg                  # Public packages
/docs                 # Documentation
```

## Implementation Guidelines

### Domain Layer

1. **Start with the domain model**:
   - Define entities and value objects
   - Implement domain logic within aggregates
   - Define domain events

```go
// Example of a domain entity
type User struct {
    ID       UserID
    Email    Email     // Value object
    Password Password  // Value object
    Role     UserRole  // Enum
    Active   bool
}

// Domain logic within entity
func (u *User) ChangePassword(currentPassword, newPassword Password) error {
    if !u.Password.Matches(currentPassword) {
        return ErrInvalidPassword
    }
    u.Password = newPassword
    return nil
}
```

2. **Use value objects for validation**:

```go
type Email string

func NewEmail(email string) (Email, error) {
    // Validation logic
    if !validEmail(email) {
        return "", ErrInvalidEmail
    }
    return Email(email), nil
}
```

### Application Layer

1. **Define commands and queries**:

```go
// Command
type CreateUser struct {
    Email    string
    Password string
}

// Query
type GetUserByID struct {
    ID string
}
```

2. **Implement command/query handlers**:

```go
// Command handler
type CreateUserHandler struct {
    userRepo UserRepository
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUser) error {
    email, err := NewEmail(cmd.Email)
    if err != nil {
        return err
    }

    password, err := NewPassword(cmd.Password)
    if err != nil {
        return err
    }

    user := NewUser(email, password)
    return h.userRepo.Save(ctx, user)
}
```

### Infrastructure Layer

1. **Implement repositories**:

```go
type PostgresUserRepository struct {
    db *sql.DB
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *User) error {
    // Implementation
}
```

2. **Use dependency injection**:

```go
func NewCreateUserHandler(repo UserRepository) *CreateUserHandler {
    return &CreateUserHandler{userRepo: repo}
}
```

## Testing

1. **Unit tests for domain logic**:
   - Test entities, value objects, and domain services in isolation

2. **Application layer tests**:
   - Test command and query handlers with mocked dependencies

3. **Integration tests**:
   - Test repositories and external services with real or containerized dependencies

4. **End-to-end tests**:
   - Test complete workflows through the entire system

## Error Handling

1. **Domain errors**:
   - Define semantic errors in the domain
   - Use custom error types with meaningful messages

2. **Application errors**:
   - Wrap domain errors with context
   - Return appropriate error codes

3. **Infrastructure errors**:
   - Translate technical errors to application errors
   - Avoid leaking infrastructure details

## Additional Considerations

1. **Event Sourcing**:
   - Consider using event sourcing for complex domains with audit requirements
   - Store the sequence of events rather than just the current state

2. **Read Models**:
   - Optimize read operations with denormalized models
   - Update read models asynchronously via domain events

3. **Performance**:
   - Profile and optimize critical paths
   - Consider caching for frequently accessed data

## Conclusion

These guidelines aim to provide a balance between architecture purity and practical implementation. The goal is to create maintainable, testable code that accurately models our business domain while leveraging Go's strengths.

Remember that these patterns should serve our needs, not constrain us. Adapt them as necessary for your specific context.
