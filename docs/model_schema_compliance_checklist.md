# Huntress API Model/Schema Compliance Checklist

<!-- REF: https://api.huntress.io/docs/preview -->
<!-- REF: https://api.huntress.io/docs#api-overview -->

> **Note:** Model/schema drift detection is now automated in CI via [`.github/workflows/model-schema-drift.yml`](../.github/workflows/model-schema-drift.yml).
> The workflow uploads a diff artifact (`model-schema-diff.txt`) comparing generated OpenAPI models to hand-written models in `pkg/huntress/`.
> Review the artifact on each PR to track compliance and drift.

This checklist is for tracking strict model validation and OpenAPI/Swagger codegen compliance for all Huntress API resources in the Bishoujo-Huntress Go client. Use this to ensure all request/response models match the latest Huntress OpenAPI/Swagger schema (field types, required/optional, enums) and that codegen/cross-checking is complete for every resource.

---

## Codebase Summary (as of 2025-05-13)

- **Core API resources** (Accounts, Organizations, Agents, Incidents, Reports, Billing, Webhooks, Audit Logs, Bulk Operations, Integrations) are fully implemented and tested.
- **All major models and endpoints** are present in `pkg/huntress/` and corresponding domain/service layers.
- **Enums** are implemented as Go enums and validated.
- **Unit and integration tests** exist for all major services and models.
- **Fuzzing harnesses** are present for core validation and encoding/decoding routines.
- **CI/CD** runs static analysis, SAST (Semgrep), dependency scanning, secret scanning, SBOM generation, and fuzzing.
- **OSSF Scorecard, CodeQL, and branch protection** are enforced.
- **Automated linting and security tool installation** is documented and enforced in the Makefile and CI.
- **OpenAPI/Swagger spec** is regularly reviewed, but strict model validation and codegen cross-checking are not yet fully automated for all resources.

---

## How to Use This Checklist

- For each resource, compare all request/response models to the latest Huntress OpenAPI/Swagger spec.
- Check for:
  - Field presence and names
  - Field types (e.g., string vs. time.Time)
  - Required vs. optional fields
  - Enum values and types
  - Nested/embedded objects
- Use OpenAPI/Swagger codegen to cross-check and validate models and endpoints.
- Mark each item as complete (`[x]`) when verified and compliant.
- Add notes for any discrepancies, workarounds, or pending updates.

---

## Checklist & Next Steps

### Accounts

- [ ] All account models match OpenAPI schema
- [ ] All account endpoints codegen-checked
- [ ] All enums/fields validated
- [ ] Tests updated for schema compliance

### Organizations

- [ ] All organization models match OpenAPI schema
- [ ] All organization endpoints codegen-checked
- [ ] All enums/fields validated
- [ ] Tests updated for schema compliance

### Agents

- [ ] All agent models match OpenAPI schema
- [ ] All agent endpoints codegen-checked
- [ ] All enums/fields validated
- [ ] Tests updated for schema compliance

### Incidents

- [ ] All incident models match OpenAPI schema
- [ ] All incident endpoints codegen-checked
- [ ] All enums/fields validated
- [ ] Tests updated for schema compliance

### Reports

- [ ] All report models match OpenAPI schema
- [ ] All report endpoints codegen-checked
- [ ] All enums/fields validated
- [ ] Tests updated for schema compliance

### Billing

- [ ] All billing models match OpenAPI schema
- [ ] All billing endpoints codegen-checked
- [ ] All enums/fields validated
- [ ] Tests updated for schema compliance

### Webhooks

- [ ] All webhook models match OpenAPI schema
- [ ] All webhook endpoints codegen-checked
- [ ] All enums/fields validated
- [ ] Tests updated for schema compliance

### Audit Logs

- [ ] All audit log models match OpenAPI schema
- [ ] All audit log endpoints codegen-checked
- [ ] All enums/fields validated
- [ ] Tests updated for schema compliance

### Bulk Operations

- [ ] All bulk operation models match OpenAPI schema
- [ ] All bulk endpoints codegen-checked
- [ ] All enums/fields validated
- [ ] Tests updated for schema compliance

### Integrations

- [ ] All integration models match OpenAPI schema
- [ ] All integration endpoints codegen-checked
- [ ] All enums/fields validated
- [ ] Tests updated for schema compliance

---

## General

- [ ] All request/response models are strictly validated against the latest OpenAPI/Swagger spec
- [ ] All enums are implemented as Go enums and validated
- [ ] All required/optional fields are correct
- [ ] All model and endpoint discrepancies are documented
- [ ] All codegen/cross-checking steps are documented
- [ ] All tests cover schema compliance and edge cases

---

**Last updated:** 2025-05-13

---

### Next Steps

1. **Download the latest Huntress OpenAPI/Swagger spec** and compare all Go models for each resource.
2. **Use OpenAPI codegen** to generate reference models and cross-check with hand-written models.
3. **Update models and tests** for any discrepancies (field types, required/optional, enums).
4. **Mark checklist items as complete** as each resource is verified and updated.
5. **Document any workarounds or known issues** in the checklist file.
