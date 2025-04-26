# Bishoujo-Huntress: Implementation Todo List

## API Coverage: What’s Present

Your codebase already covers the following Huntress API resources and operations:

- **Accounts**
  - Get current account details
  - Update account settings
  - List, add, update, and remove users
  - Get account statistics
- **Organizations**
  - CRUD operations (create, get, update, delete)
  - List organizations
  - Manage organization users
- **Agents**
  - Get agent details
  - List agents (with filters)
  - Update and delete agents
  - Get agent statistics
- **Incidents**
  - Get incident details
  - List incidents (with filters)
  - Update incident status
  - Assign incidents
- **Reports**
  - Generate, get, list, download, export, schedule reports
  - Get summary and detailed reports
- **Billing**
  - Get billing summary
  - List and get invoices
  - Get usage statistics
- **Infrastructure**
  - Pagination, error handling, authentication, and rate limiting are all present in the infrastructure.

---

## What’s Missing or Incomplete (vs. Huntress API Docs/Swagger)

### A. API Endpoints Not Yet Implemented

- **Webhooks**
  - The Huntress API supports webhooks for event notifications (see `/webhooks` endpoints).
  - Your codebase does not expose webhook management (list, create, update, delete webhooks).
- **Bulk Operations**
  - The API supports some bulk operations (e.g., bulk agent actions, bulk organization updates).
  - No bulk operation helpers or endpoints are present in your client.
- **Audit Logs**
  - The API provides endpoints for audit logs (e.g., `/audit-logs`).
  - No audit log service or models are present.
- **Integrations**
  - Endpoints for managing integrations (e.g., `/integrations`, `/integration-settings`) are not present.
- **User Invitations**
  - The API supports inviting users (e.g., `/users/invite`).
  - No explicit invite user method is present.
- **API Versioning**
  - The API supports versioning via headers or URL.
  - Your client has a version field, but no explicit support for switching API versions or custom headers.

### B. Features/Behaviors Not Fully Implemented

- **Advanced Filtering and Search**
  - Some endpoints support advanced filtering, searching, and sorting (e.g., query params for incidents, agents, organizations).
  - Your List methods support basic filters, but may not expose all available query params (e.g., date ranges, status enums, search terms).
- **Response Caching**
  - No caching layer for GET requests (optional, but can improve performance).
- **Request Tracing**
  - No explicit support for request tracing (e.g., correlation IDs, OpenTelemetry).
- **Webhook Event Handling**
  - No helpers for validating or parsing webhook payloads.
- **API Error Types**
  - You have a solid error model, but may want to add more granular error types for specific API error codes.
- **User-Agent Customization**
  - User-Agent is set, but not easily customizable per request.
- **Testing/Mocking**
  - No explicit test helpers or mock server for simulating Huntress API responses.

---

## Swagger/Docs-Driven Enhancements

- **Strict Model Validation**
  - Ensure all request/response models match the latest Swagger schema (field types, required/optional, enums).
  - Some fields may be missing or have different types (e.g., time.Time vs. string).
- **Enum Types**
  - Use Go enums for fields with limited values (e.g., incident status, agent platform).
- **OpenAPI/Swagger Codegen**
  - Consider using Swagger/OpenAPI codegen to cross-check your models and endpoints for completeness.

---

## Next Steps (Go Library Best Practices & Project Docs)

- Prioritize missing endpoints and features based on user needs and Huntress API changes.
- Continue to follow DDD, CQRS, and Clean Architecture as described in `docs/`.
- Ensure all new code is covered by unit and integration tests.
- Keep documentation and examples up to date as new features are added.
- Regularly review the Swagger/OpenAPI spec for new endpoints or changes.

---

# (Previous TODOs and Progress Tracking)

## Core Client Implementation

- [x] Basic client structure in pkg/huntress/client.go
- [x] Options pattern for configuration
- [x] Service interfaces for different API resources
- [x] Complete client implementation
  - [x] Basic authentication
  - [x] Service implementation initialization
  - [x] Request execution logic with context handling
  - [x] Implement proper service instance creation

## Resource Service Implementation

- [x] Complete Account service implementation
  - [x] Get account details
  - [x] Update account settings
  - [x] Manage account users
  - [x] Account statistics endpoints
- [x] Complete Organization service implementation
  - [x] Create organizations
  - [x] Get organization details
  - [x] Update organization settings
  - [x] Delete organizations
  - [x] Manage organization users
- [x] Complete Agent service implementation
  - [x] Get agent details
  - [x] List agents with filtering
  - [x] Agent status updates
  - [x] Agent statistics
- [x] Complete Incident service implementation
  - [x] List incident reports
  - [x] Get incident details
  - [x] Update incident status
  - [x] Generate incident reports
- [x] Complete Report service implementation
  - [x] Generate summary reports
  - [x] Retrieve report data
  - [x] Export reports
- [x] Complete Billing service implementation
  - [x] Get billing information
  - [x] List invoices
  - [x] Usage statistics
- [x] Complete Webhook service interface and scaffolding
  - [x] Implement Webhook API calls (List, Get, Create, Update, Delete)

## Infrastructure Layer Implementation

- [x] Basic HTTP client structure
- [x] Repository implementation scaffolding for Account
- [x] Repository implementation scaffolding for Agent
- [x] Repository implementation scaffolding for Organization
- [x] Repository implementation scaffolding for Webhook
- [x] Repository implementation scaffolding for Incident
- [x] HTTP client retry logic implementation
- [x] Complete HTTP client implementation
  - [x] Basic retry logic with backoff
  - [x] Context handling for cancellation
  - [x] Authentication handling
  - [x] Rate limiting (60 requests per minute)
  - [x] Error handling and response mapping
  - [x] Pagination handling
- [x] Implement HTTP client utilities
  - [x] Retry configuration
  - [x] Backoff calculation
  - [x] Request/Response logging
  - [x] Metrics collection
  - [x] Debug mode
- [x] Complete repository implementations
  - [x] Account repository
  - [x] Organization repository
  - [x] Agent repository
  - [x] Incident repository
  - [x] Report repository
  - [x] Billing repository

## Domain Model Implementation

- [x] Complete Account domain model
  - [x] Basic entity structure
  - [x] Complete validation logic
  - [x] Domain events
- [x] Complete Organization domain model
  - [x] Basic entity structure
  - [x] Complete validation logic
  - [x] Domain events
- [x] Complete Agent domain model
  - [x] Basic entity structure
  - [x] Complete validation logic
  - [x] Domain events
- [x] Complete Incident domain model
  - [x] Basic entity structure
  - [x] Complete validation logic
  - [x] Domain events
- [x] Complete Report domain model
  - [x] Define entity structure
  - [x] Implement validation logic
  - [x] Domain events
- [x] Complete Billing domain model
  - [x] Define entity structure
  - [x] Implement validation logic
  - [x] Domain events

## API-Specific Enhancements

### Rate Limiting

- [x] Initial retry mechanism for 429 responses
- [x] Implement sliding window rate limiter (60 requests/minute)
- [x] Add proactive rate limiting
- [x] Create request queue for high-volume operations
- [x] Implement rate limit tracking across concurrent requests

### Error Handling

- [x] Basic retry configuration
- [x] Retryable status code handling
- [x] API-specific error types
- [x] Error response mapping
- [x] Detailed error context
- [x] Error recovery strategies
- [x] Custom error types per domain

### Authentication

- [x] Basic auth implementation
- [x] Secure credential storage
- [x] Auth header handling
- [x] Token refresh mechanism (if needed)
- [x] Support for environment variable configuration

### API Client Features

- [x] Add API version support
- [x] User agent configuration
- [x] Request tracing
- [x] Response caching
- [x] Automatic retries for idempotent operations
- [x] Webhook handling for event notifications
- [x] Bulk operations support

## Application Layer Implementation

- [x] Command Handlers
  - [x] Create organization command
  - [x] Update organization command
  - [x] Delete organization command
  - [x] Update agent settings command
  - [x] Update incident status command
  - [x] Generate report command
- [x] Query Handlers
  - [x] List organizations query
  - [x] Get organization details query
  - [x] List agents query
  - [x] Get agent details query
  - [x] List incidents query
  - [x] Get incident details query
  - [x] List reports query
  - [x] Get report query
  - [x] Download report query
  - [x] Get summary report query

## Testing Improvements

### HTTP Client Tests

- [x] Retry logic unit tests
- [x] Rate limiting tests
- [x] Auth handling tests
- [x] Mock server implementation
- [x] Integration test suite

### Repository Tests

- [x] Mock HTTP client for testing
- [x] Repository unit tests
- [x] Integration tests with API
- [x] Error handling tests
- [x] Rate limit handling tests

### Domain Model Tests

- [x] Entity validation tests
- [x] Value object tests
- [x] Domain service tests
- [x] Entity factory tests

### Application Layer Tests

- [x] Command handler tests
- [x] Query handler tests
- [x] Integration tests between layers

## Documentation Updates

### API Documentation

- [x] Complete API method documentation
- [x] Authentication examples
- [x] Rate limiting examples
- [x] Error handling examples
- [x] Pagination examples

### Developer Documentation

- [x] Contributing guidelines
- [x] Architecture overview
- [x] Testing guide
- [x] Release process
- [x] Security considerations

### Example Applications

- [x] Basic client usage examples
- [x] Organization management examples
- [x] Agent monitoring examples
- [x] Incident response examples
- [x] Reporting examples
- [x] CLI tool for common operations

## CI/CD and Project Infrastructure

- [x] GitHub Actions workflow
- [x] Automated testing
- [x] Code coverage reporting
- [x] Static analysis
- [x] Documentation generation
- [x] Version tagging and release automation
- [x] Go module versioning strategy

---

## Next Steps (Spring 2025)

- [x] Expand advanced filtering in all list methods (date ranges, enums, search, tags)
- [x] Review and update error handling for context and mapping

### Progress Notes (Spring 2025)

- Strict model validation and enum enforcement are now implemented for all create/update operations and public API params.
- Unit tests for validation, service stubs, and error cases are present in `pkg/huntress/*_test.go`.
- Enum usage and validation are documented in README and API docs.
- Advanced filtering and error context/mapping review remain as next priorities.

## Swagger/OpenAPI Considerations

- [ ] Periodically review the Huntress API Swagger/OpenAPI spec to ensure all models and endpoints are up to date.
- [ ] Consider using Swagger/OpenAPI codegen to cross-check Go models and endpoints for completeness and accuracy.
- [ ] Document any intentional deviations from the upstream API spec.
- [ ] Use Swagger/OpenAPI for automated documentation or SDK generation if needed in the future.

> Note: Swagger/OpenAPI integration is not required for OSSF/security compliance, but is recommended for long-term maintainability and completeness.
