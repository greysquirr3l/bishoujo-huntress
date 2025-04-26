# Bishoujo-Huntress: Copilot Instructions

## Project Overview

Bishoujo-Huntress is a Go client library for the Huntress API, implementing Domain-Driven Design (DDD), CQRS, and Clean Architecture principles. The library provides a comprehensive, type-safe interface for all Huntress API endpoints, with a focus on maintainability, testability, and security.

- **Status:** Core API resources (Accounts, Organizations, Agents, Incidents, Reports, Billing) are fully implemented and tested. Webhooks are scaffolded. Bulk operations, audit logs, and integrations are planned.
- **API Coverage:** See [docs/todo.md](../docs/todo.md) for detailed status and roadmap.
- **Security:** Follows [OSSF Security Baselines](../docs/OSSF_SECURITY_BASELINES.md) and Go secure coding best practices.

## Key Resources

- [Architecture Design](../docs/architecture.md)
- [Project Structure](../docs/project_structure.md)
- [Project Guidelines](../docs/project_guidelines.md)
- [OSSF Security Baselines](../docs/OSSF_SECURITY_BASELINES.md)
- [Implementation Todo List](../docs/todo.md)

## Development Guidelines

### Code Structure

- **Domain Entities:** `internal/domain/{resource}`
- **Repository Interfaces:** `internal/ports/repository`
- **API Implementation:** `internal/adapters/api`, `internal/adapters/repository`
- **Public API:** `pkg/huntress` (all main client/service logic)
- **Examples:** `cmd/examples/`

### Error Handling

- Domain, application, and infrastructure errors are separated and wrapped with context.
- All API errors are mapped to Go error types (see `pkg/huntress/errors.go`).

### Authentication

- Basic Authentication (API key/secret) is required.
- No hardcoded credentials; use environment variables or config.

### Testing

- Unit and integration tests for all major services.
- Run `make test` to execute the full suite.
- Example usage: `cmd/examples/basic/main.go`.

### Security

- Follows [OSSF Security Baselines](../docs/OSSF_SECURITY_BASELINES.md).
- Input validation, secure error handling, minimal dependencies.
- No hardcoded secrets; all credentials must be provided securely.

## API Structure

- **Client Configuration:**
  ```go
  client := huntress.New(
      huntress.WithCredentials("API_KEY", "API_SECRET"),
      huntress.WithTimeout(30 * time.Second),
  )
  ```
- **Resource Operations:**
  ```go
  entity, err := client.Resource.Get(ctx, id)
  entities, pagination, err := client.Resource.List(ctx, &resource.ListParams{...})
  result, err := client.Resource.CustomOperation(ctx, params)
  ```

## Implementation Status and Priorities

- See [docs/todo.md](../docs/todo.md) for up-to-date implementation status.
- All core API resources are implemented and tested.
- Webhook CRUD is scaffolded; bulk ops, audit logs, integrations are planned.
- All new code must include tests and documentation.

## Testing & Linting

- Run `make test` for all tests.
- Run `make lint` for static analysis (uses `golangci-lint`).
- Run `gosec ./...` and `govulncheck ./...` for security scanning.

## Security & OSSF Baselines

- See [OSSF Security Baselines](../docs/OSSF_SECURITY_BASELINES.md) for requirements and checklist.
- Security policy: [SECURITY.md](../SECURITY.md)

## Additional Resources

- [Huntress API Docs](https://api.huntress.io/docs#api-overview)
- [pkg.go.dev Reference](https://pkg.go.dev/github.com/greysquirr3l/bishoujo-huntress)

---

**How to Test OSSF Security Baselines Locally (macOS):**

1. **Static Analysis:**

   - Run `make lint` (uses `golangci-lint` with security linters enabled)
   - Run `gosec ./...` for Go security checks
   - Run `govulncheck ./...` for Go vulnerability scanning

2. **Dependency Scanning:**

   - Run `go mod tidy` to ensure clean dependencies
   - Run `govulncheck ./...` to check for known vulnerabilities

3. **Secret Scanning:**

   - Install [git-secrets](https://github.com/awslabs/git-secrets):
     ```bash
     brew install git-secrets
     git secrets --scan
     ```

4. **SBOM Generation:**

   - Install [syft](https://github.com/anchore/syft):
     ```bash
     brew install syft
     syft . -o cyclonedx-json > sbom.json
     ```

5. **Test Coverage:**

   - Run `make test` and check for 100% coverage on critical code paths

6. **CI/CD:**

   - Ensure GitHub Actions workflow runs all of the above on PRs

7. **Security Policy:**
   - Confirm [SECURITY.md](../SECURITY.md) is present and up to date

**To testify adherence:**

- Save the output of the above tools (e.g., `gosec`, `govulncheck`, `git secrets`, `syft`) and include them in a security review or as CI artifacts.
- Ensure all checkboxes in [OSSF_SECURITY_BASELINES.md](../docs/OSSF_SECURITY_BASELINES.md) are addressed and documented.
