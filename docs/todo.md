
# Bishoujo-Huntress: Implementation Todo List

## Open Items

### What’s Missing or Incomplete (vs. Huntress API Docs/Swagger)

- [ ] **User Invitations**
  - [ ] No explicit invite user method yet
[x] **API Versioning**
  - [x] Version field present; explicit switching/custom headers implemented and tested
[x] **Request Tracing**
  - [x] Request tracing (correlation IDs, OpenTelemetry) implemented in HTTP client and infrastructure
[x] **Testing/Mocking**
  - [x] Test helpers and mock server for simulating Huntress API responses implemented and in use

### Features/Enhancements (vs. Huntress API Docs/Swagger)

- [ ] **Strict Model Validation**
  - [ ] Ensure all request/response models match the latest Swagger schema (field types, required/optional, enums).
  - [ ] Some fields may be missing or have different types (e.g., time.Time vs. string).
- [x] **Enum Types**
  - [x] Go enums are now used for all major fields with limited values (incident status, agent platform, report type, webhook event type, etc.)
- [ ] **OpenAPI/Swagger Codegen**
  - [ ] Consider using Swagger/OpenAPI codegen to cross-check your models and endpoints for completeness and accuracy.

### Mock Generation & Testing

- [x] **Manual/Testify Mock Implementations**
  - [x] Interfaces are mocked using hand-written or testify-based mocks in test files.
  - [x] No use of Mockery or auto-generated mocks; mocks are maintained alongside tests.
  - [x] All critical interfaces have corresponding test mocks for unit and integration testing.

- [x] Advanced filtering in all list methods (date ranges, enums, search, tags) implemented and tested
- [x] Error handling reviewed and updated for context and mapping
- [x] Complete bulk, audit log, and integration adapters
- [x] Add user invitation helpers (helpers scaffolded, API endpoint pending)
- [x] Request tracing and response caching implemented (optional, can be toggled)
- [x] Huntress API Swagger/OpenAPI spec reviewed for new endpoints or changes (as of May 2025)
- [x] Prioritize missing endpoints and features based on user needs and Huntress API changes
- [x] DDD, CQRS, and Clean Architecture followed as described in `docs/`
- [x] All new code is covered by unit and integration tests
- [x] Documentation and examples up to date as new features are added
- [x] Swagger/OpenAPI spec regularly reviewed for new endpoints or changes

### Linting Automation

- [x] `golangci-lint` installation is now automated in the Makefile and documented in the README. The Makefile will install the correct version if not present in your `$PATH`.

### 🚩 Required & Missing Tests (Spring 2025 Audit)

#### Domain Layer

- [x] **Account domain**
  - [x] Test all account entity validation logic
  - [x] Test account domain events
- [x] **Organization domain**
  - [x] Test organization entity validation logic
  - [x] Test organization domain events
- [x] **Agent domain**
  - [x] Test agent entity validation logic
  - [x] Test agent platform enum logic
  - [x] Test agent domain events
- [x] **Incident domain**
  - [x] Test incident entity validation logic
  - [x] Test incident status enum logic
  - [x] Test incident domain events
- [x] **Report domain**
  - [x] Test report entity validation logic
  - [x] Test report type enum logic
  - [x] Test report domain events
- [x] **Billing domain**
  - [x] Test billing entity validation logic
  - [x] Test billing domain events
- [x] **Webhook domain**
  - [x] Test webhook entity validation logic
  - [x] Test webhook event type enums

#### Application Layer

- [x] **Command Handlers**
  - [x] Test create/update/delete org command handlers (happy path and error cases)
  - [x] Test agent update/status command handlers
  - [x] Test incident update/assign command handlers
  - [x] Test report generation command handlers
- [x] **Query Handlers**
  - [x] Test list/get org query handlers (with filters)
  - [x] Test list/get agent query handlers (with filters)
  - [x] Test list/get incident query handlers (with filters)
  - [x] Test list/get report query handlers (with filters)
  - [x] Test billing query handlers

#### Infrastructure Layer

- [x] **HTTP Client**
  - [x] Test retry logic for all error codes
  - [x] Test rate limiting under concurrent load
  - [x] Test authentication header injection
  - [x] Test pagination handling for all list endpoints
  - [x] Test error mapping for all API error types
  - [x] Test context cancellation and timeout handling
- [x] **Repository Implementations**
  - [x] Test account repository (CRUD, error, pagination)
  - [x] Test organization repository (CRUD, error, pagination)
  - [x] Test agent repository (CRUD, error, pagination)
  - [x] Test incident repository (CRUD, error, pagination)
  - [x] Test report repository (CRUD, error, pagination)
  - [x] Test billing repository (CRUD, error, pagination)
  - [x] Test webhook repository (CRUD, error, pagination)
  - [x] Test integration repository (CRUD, error, pagination)
- [x] **Adapters**
  - [x] Test all API adapters for correct DTO mapping
  - [x] Test all repository adapters for correct domain mapping

#### Public API (`pkg/huntress`)

- [x] **Client**
  - [x] Test client initialization with all options
  - [x] Test error handling for invalid config
  - [x] Test all service method calls (happy path and error)
- [x] **Service Methods**
  - [x] Test all account service methods
  - [x] Test all organization service methods
  - [x] Test all agent service methods
  - [x] Test all incident service methods
  - [x] Test all report service methods
  - [x] Test all billing service methods
  - [x] Test all webhook service methods
  - [x] Test all integration service methods

#### Security & Compliance

- [x] **Security Tests**
  - [x] Test for secret leakage in logs/errors  # see account_service_test.go, organization_service_test.go, logs/errors tests, git-secrets.txt
  - [x] Test input validation for all public API methods  # see account_service_test.go, organization_service_test.go, agent_service_test.go, etc.
  - [x] Test error handling does not leak sensitive info  # see account_service_test.go, error message assertions
  - [x] Test SBOM generation and artifact upload  # see sbom.json, ci.yml
  - [x] Test static analysis (gosec, govulncheck) passes on all code  # see gosec.txt, govulncheck.txt, ci.yml
- [x] **OSSF Scorecard**
  - [x] Test all workflow artifacts are generated and uploaded  # see ci.yml, scorecard.yml, artifact upload steps
  - [x] Test all actions are pinned by SHA in workflows  # see ci.yml, scorecard.yml, github-ci-hash.sh

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
