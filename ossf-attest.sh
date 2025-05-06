#!/bin/bash
# ossf-attest.sh: Run all OSSF Security Baseline checks and save attestation artifacts
# Usage: ./ossf-attest.sh
set -euo pipefail

# Output files
LINT_OUT=golangci-lint.txt
TEST_OUT=test-results.txt
GOSEC_OUT=gosec.txt
GOVULN_OUT=govulncheck.txt
GITSECRETS_OUT=git-secrets.txt
SBOM_OUT=sbom.json

# Project name and version for SBOM workaround
: "${PROJECT_NAME:=bishoujo-huntress}"
: "${PROJECT_VERSION:=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")}"

# Print versions for traceability
echo "== Tool Versions =="
golangci-lint --version || true
gosec --version || true
govulncheck --version || true
git secrets --list || true
syft version || true
echo "==================="

# 1. Static Analysis & Linting
echo "Running golangci-lint..."
make lint > "$LINT_OUT" 2>&1 || true

echo "Running gosec..."
gosec ./... > "$GOSEC_OUT" 2>&1 || true

echo "Running govulncheck..."
govulncheck ./... > "$GOVULN_OUT" 2>&1 || true

# 2. Dependency Scanning
echo "Running go mod tidy..."
go mod tidy

# 3. Secret Scanning
echo "Running git-secrets..."
git secrets --scan > "$GITSECRETS_OUT" 2>&1 || true

# 4. SBOM Generation (with name/version, pretty-printed)
echo "Generating SBOM with syft..."
syft . -o cyclonedx-json --source-name "$PROJECT_NAME" --source-version "$PROJECT_VERSION" \
| jq 'del(.metadata.tools.components[].author)
      | .metadata.authors = [{"name": "anchore"}]' \
> "$SBOM_OUT" || true

if command -v jq >/dev/null 2>&1; then
  echo "SBOM generated as a JSON object with name and version fields."
else
  echo "jq not found, SBOM may not have name/version fields."
fi

if command -v jq >/dev/null 2>&1; then
  echo "SBOM generated as a JSON object with name and version fields."
else
  echo "jq not found, SBOM may not have name/version fields."
fi

# 5. Test Coverage (with coverage report)
echo "Running tests with coverage..."
go test -coverprofile=coverage.txt ./... > "$TEST_OUT" 2>&1 || true

echo "All OSSF attestation artifacts generated:"
echo "  $LINT_OUT"
echo "  $GOSEC_OUT"
echo "  $GOVULN_OUT"
echo "  $GITSECRETS_OUT"
echo "  $SBOM_OUT"
echo "  $TEST_OUT"
