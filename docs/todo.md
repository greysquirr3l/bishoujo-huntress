
# Bishoujo-Huntress: Implementation Todo List

## Open Items

### What’s Missing or Incomplete (vs. Huntress API Docs/Swagger)

- [ ] **User Invitations**
  - [ ] No explicit invite user method yet
- [ ] **API Versioning**
  - [ ] Version field present, explicit switching/custom headers not yet implemented
- [ ] **Request Tracing**
  - [ ] No explicit support for request tracing (e.g., correlation IDs, OpenTelemetry)
- [ ] **Testing/Mocking**
  - [ ] No explicit test helpers or mock server for simulating Huntress API responses

### Features/Enhancements (vs. Huntress API Docs/Swagger)

- [ ] **Strict Model Validation**
  - [ ] Ensure all request/response models match the latest Swagger schema (field types, required/optional, enums).
  - [ ] Some fields may be missing or have different types (e.g., time.Time vs. string).
- [ ] **Enum Types**
  - [ ] Use Go enums for fields with limited values (e.g., incident status, agent platform).
- [ ] **OpenAPI/Swagger Codegen**
  - [ ] Consider using Swagger/OpenAPI codegen to cross-check your models and endpoints for completeness and accuracy.

### Next Steps

- Expand advanced filtering in all list methods (date ranges, enums, search, tags)
- Review and update error handling for context and mapping
- Complete bulk, audit log, and integration adapters
- Add user invitation helpers
- Add request tracing and response caching (optional)
- Periodically review the Huntress API Swagger/OpenAPI spec for new endpoints or changes
- Prioritize missing endpoints and features based on user needs and Huntress API changes
- Continue to follow DDD, CQRS, and Clean Architecture as described in `docs/`
- Ensure all new code is covered by unit and integration tests
- Keep documentation and examples up to date as new features are added
- Regularly review the Swagger/OpenAPI spec for new endpoints or changes

### 🚩 Required & Missing Tests (Spring 2025 Audit)

#### Domain Layer

- [ ] **Account domain**
  - [ ] Test all account entity validation logic
  - [ ] Test account domain events
- [ ] **Organization domain**
  - [ ] Test organization entity validation logic
  - [ ] Test organization domain events
- [ ] **Agent domain**
  - [ ] Test agent entity validation logic
  - [ ] Test agent platform enum logic
  - [ ] Test agent domain events
- [ ] **Incident domain**
  - [ ] Test incident entity validation logic
  - [ ] Test incident status enum logic
  - [ ] Test incident domain events
- [ ] **Report domain**
  - [ ] Test report entity validation logic
  - [ ] Test report type enum logic
  - [ ] Test report domain events
- [ ] **Billing domain**
  - [ ] Test billing entity validation logic
  - [ ] Test billing domain events
- [ ] **Webhook domain**
  - [ ] Test webhook entity validation logic
  - [ ] Test webhook event type enums

#### Application Layer

- [ ] **Command Handlers**
  - [ ] Test create/update/delete org command handlers (happy path and error cases)
  - [ ] Test agent update/status command handlers
  - [ ] Test incident update/assign command handlers
  - [ ] Test report generation command handlers
- [ ] **Query Handlers**
  - [ ] Test list/get org query handlers (with filters)
  - [ ] Test list/get agent query handlers (with filters)
  - [ ] Test list/get incident query handlers (with filters)
  - [ ] Test list/get report query handlers (with filters)
  - [ ] Test billing query handlers

#### Infrastructure Layer

- [ ] **HTTP Client**
  - [ ] Test retry logic for all error codes
  - [ ] Test rate limiting under concurrent load
  - [ ] Test authentication header injection
  - [ ] Test pagination handling for all list endpoints
  - [ ] Test error mapping for all API error types
  - [ ] Test context cancellation and timeout handling
- [ ] **Repository Implementations**
  - [ ] Test account repository (CRUD, error, pagination)
  - [ ] Test organization repository (CRUD, error, pagination)
  - [ ] Test agent repository (CRUD, error, pagination)
  - [ ] Test incident repository (CRUD, error, pagination)
  - [ ] Test report repository (CRUD, error, pagination)
  - [ ] Test billing repository (CRUD, error, pagination)
  - [ ] Test webhook repository (CRUD, error, pagination)
  - [ ] Test integration repository (CRUD, error, pagination)
- [ ] **Adapters**
  - [ ] Test all API adapters for correct DTO mapping
  - [ ] Test all repository adapters for correct domain mapping

#### Public API (`pkg/huntress`)

- [ ] **Client**
  - [ ] Test client initialization with all options
  - [ ] Test error handling for invalid config
  - [ ] Test all service method calls (happy path and error)
- [ ] **Service Methods**
  - [ ] Test all account service methods
  - [ ] Test all organization service methods
  - [ ] Test all agent service methods
  - [ ] Test all incident service methods
  - [ ] Test all report service methods
  - [ ] Test all billing service methods
  - [ ] Test all webhook service methods
  - [ ] Test all integration service methods

#### Security & Compliance

- [ ] **Security Tests**
  - [ ] Test for secret leakage in logs/errors
  - [ ] Test input validation for all public API methods
  - [ ] Test error handling does not leak sensitive info
  - [ ] Test SBOM generation and artifact upload
  - [ ] Test static analysis (gosec, govulncheck) passes on all code
- [ ] **OSSF Scorecard**
  - [ ] Test all workflow artifacts are generated and uploaded
  - [ ] Test all actions are pinned by SHA in workflows

#### CLI & Examples

- [ ] **Examples**
  - [ ] Test all CLI examples in `cmd/examples` run without error
  - [ ] Test documentation code snippets compile and run

#### Integration & End-to-End

- [ ] **Integration Tests**
  - [ ] Test end-to-end flows for each resource (account, org, agent, incident, report, billing, webhook)
  - [ ] Test error propagation across layers
  - [ ] Test concurrent usage of the client (thread safety)

### Test Coverage Goals

- [ ] **Domain Layer:** 100% of business logic and validation
- [ ] **Application Layer:** 100% of command/query handlers
- [ ] **Infrastructure Layer:** 90%+ of HTTP/repo logic, 100% of error mapping
- [ ] **Public API:** 100% of exported methods
- [ ] **Security:** 100% of input validation and error handling
- [ ] **End-to-End:** At least one test for each major workflow

**Next Steps:**
- [ ] Add missing unit and integration tests as listed above.
- [ ] Review test coverage reports after each CI run.
- [ ] Ensure all new features and bugfixes include corresponding tests.
- [ ] Periodically review and update this checklist as the codebase evolves.

---

## Completed Items

### API Coverage: What’s Present

- [x] **Accounts**
  - [x] Get current account details
  - [x] Update account settings
  - [x] List, add, update, and remove users
  - [x] Get account statistics
- [x] **Organizations**
  - [x] CRUD operations (create, get, update, delete)
  - [x] List organizations
  - [x] Manage organization users
- [x] **Agents**
  - [x] Get agent details
  - [x] List agents (with filters)
  - [x] Update and delete agents
  - [x] Get agent statistics
- [x] **Incidents**
  - [x] Get incident details
  - [x] List incidents (with filters)
  - [x] Update incident status
  - [x] Assign incidents
- [x] **Reports**
  - [x] Generate, get, list, download, export, schedule reports
  - [x] Get summary and detailed reports
- [x] **Billing**
  - [x] Get billing summary
  - [x] List and get invoices
  - [x] Get usage statistics
- [x] **Webhooks**
  - [x] CRUD operations (create, get, update, delete)
  - [x] List webhooks
  - [x] Webhook service and models implemented
- [x] **Infrastructure**
  - [x] Pagination, error handling, authentication, and rate limiting are all present in the infrastructure.

### Features/Enhancements (vs. Huntress API Docs/Swagger)

- [x] **Bulk Operations**
- [x] **Integrations**
- [x] **Advanced Filtering and Search**
- [x] **Response Caching**
- [x] **Webhook Event Handling**
- [x] **API Error Types**
- [x] **User-Agent Customization**

### API Endpoints Implemented

- [x] **Webhooks**
  - [x] Webhook management (list, create, update, delete) fully implemented, tested, and documented.
- [x] **Bulk Operations**
  - [x] Bulk agent/org actions: helpers and endpoints implemented
- [x] **Audit Logs**
  - [x] Full audit log service and domain models implemented.
  - [x] Application layer (command/query) handlers for audit logs implemented.
  - [x] Public API and documentation for audit log operations present.
  - [x] Tests for audit log functionality present.
  - [x] Implementation and documentation of audit log service/models complete.
- [x] **Integrations**
  - [x] Endpoints for managing integrations (e.g., `/integrations`, `/integration-settings`) are implemented and documented.

### Progress Notes (Spring 2025)

- All core API resources (Accounts, Organizations, Agents, Incidents, Reports, Billing) are fully implemented and tested.
- Webhook CRUD is now fully implemented and tested.
- Bulk, audit log, and integration adapters are scaffolded.
- Strict model validation and enum enforcement are implemented for all create/update operations and public API params.
- Unit and integration tests for validation, service stubs, and error cases are present.
- Enum usage and validation are documented in README and API docs.
- Advanced filtering and error context/mapping review remain as next priorities.

### Core Client Implementation

- [x] Basic client structure in pkg/huntress/client.go
- [x] Options pattern for configuration
- [x] Service interfaces for different API resources
- [x] Complete client implementation
  - [x] Basic authentication
  - [x] Service implementation initialization
  - [x] Request execution logic with context handling
  - [x] Implement proper service instance creation

### Resource Service Implementation

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

### Infrastructure Layer Implementation

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

### Domain Model Implementation

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

### API-Specific Enhancements

- [x] Initial retry mechanism for 429 responses
- [x] Implement sliding window rate limiter (60 requests/minute)
- [x] Add proactive rate limiting
- [x] Create request queue for high-volume operations
- [x] Implement rate limit tracking across concurrent requests
- [x] Basic retry configuration
- [x] Retryable status code handling
- [x] API-specific error types
- [x] Error response mapping
- [x] Detailed error context
- [x] Error recovery strategies
- [x] Custom error types per domain
- [x] Basic auth implementation
- [x] Secure credential storage
- [x] Auth header handling
- [x] Token refresh mechanism (if needed)
- [x] Support for environment variable configuration
- [x] Add API version support
- [x] User agent configuration
- [x] Request tracing
- [x] Response caching
- [x] Automatic retries for idempotent operations
- [x] Webhook handling for event notifications
- [x] Bulk operations support

### Application Layer Implementation

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

### Testing Improvements

- [x] Retry logic unit tests
- [x] Rate limiting tests
- [x] Auth handling tests
- [x] Mock server implementation
- [x] Integration test suite
- [x] Mock HTTP client for testing
- [x] Repository unit tests
- [x] Integration tests with API
- [x] Error handling tests
- [x] Rate limit handling tests
- [x] Entity validation tests
- [x] Value object tests
- [x] Domain service tests
- [x] Entity factory tests
- [x] Command handler tests
- [x] Query handler tests
- [x] Integration tests between layers

### Documentation Updates

- [x] Complete API method documentation
- [x] Authentication examples
- [x] Rate limiting examples
- [x] Error handling examples
- [x] Pagination examples
- [x] Contributing guidelines
- [x] Architecture overview
- [x] Testing guide
- [x] Release process
- [x] Security considerations
- [x] Basic client usage examples
- [x] Organization management examples
- [x] Agent monitoring examples
- [x] Incident response examples
- [x] Reporting examples
- [x] CLI tool for common operations

### CI/CD and Project Infrastructure

- [x] GitHub Actions workflow
- [x] Automated testing
- [x] Code coverage reporting
- [x] Static analysis
- [x] Documentation generation
- [x] Version tagging and release automation
- [x] Go module versioning strategy

### Next Steps (Spring 2025)

- [x] Expand advanced filtering in all list methods (date ranges, enums, search, tags)
- [x] Review and update error handling for context and mapping

---
