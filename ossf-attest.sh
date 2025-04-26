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
PROJECT_NAME="bishoujo-huntress"
PROJECT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")

# Print versions for traceability
echo "== Tool Versions =="
golangci-lint --version || true
gosec --version || true
govulncheck --version || true
git secrets --version || true
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

# 4. SBOM Generation (with workaround for name/version)
echo "Generating SBOM with syft..."
syft . -o cyclonedx-json > "$SBOM_OUT" 2>&1 || true
if command -v jq >/dev/null 2>&1; then
  echo "Injecting name and version into SBOM (workaround)..."
  jq --arg name "$PROJECT_NAME" --arg version "$PROJECT_VERSION" '(.name // $name) as $n | (.version // $version) as $v | . + {name: $n, version: $v}' "$SBOM_OUT" > "${SBOM_OUT}.tmp" && mv "${SBOM_OUT}.tmp" "$SBOM_OUT"
else
  echo "jq not found, skipping SBOM name/version injection workaround."
fi

# 5. Test Coverage
echo "Running tests..."
make test > "$TEST_OUT" 2>&1 || true

echo "All OSSF attestation artifacts generated:"
echo "  $LINT_OUT"
echo "  $GOSEC_OUT"
echo "  $GOVULN_OUT"
echo "  $GITSECRETS_OUT"
echo "  $SBOM_OUT"
echo "  $TEST_OUT"
