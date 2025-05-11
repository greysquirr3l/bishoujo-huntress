#!/bin/bash
# ossf-attest.sh: Run all OSSF Security Baseline checks and save attestation artifacts
set -euo pipefail

# Check for required tools
for tool in golangci-lint gosec govulncheck git syft jq; do
  if ! command -v "$tool" >/dev/null 2>&1; then
    echo "ERROR: Required tool '$tool' is not installed or not in PATH."
    exit 1
  fi
done

# Output files
LINT_OUT=golangci-lint.txt
TEST_OUT=test-results.txt
GOSEC_OUT=gosec.txt
GOVULN_OUT=govulncheck.txt
GITSECRETS_OUT=git-secrets.txt
SBOM_OUT=sbom.json
COVERAGE_OUT=coverage.txt

: "${PROJECT_NAME:=bishoujo-huntress}"
: "${PROJECT_VERSION:=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")}"

echo "== Tool Versions =="
go version || true
go env GOMOD || true
golangci-lint --version || true
gosec --version || true
govulncheck --version || true
git secrets --list || true
syft version || true
jq --version || true
echo "==================="

echo "== Running OSSF Security Baseline Checks =="
echo "Running golangci-lint..."
make lint > "$LINT_OUT" 2>&1 || true &

echo "Running gosec..."
gosec ./... > "$GOSEC_OUT" 2>&1 || true &

echo "Running govulncheck..."
govulncheck ./... > "$GOVULN_OUT" 2>&1 || true &

echo "Running git-secrets..."
git secrets --scan > "$GITSECRETS_OUT" 2>&1 || true &

wait

echo "Running go mod tidy..."
go mod tidy

echo "Generating SBOM with syft..."
if command -v jq >/dev/null 2>&1; then
  syft . -o cyclonedx-json --source-name "$PROJECT_NAME" --source-version "$PROJECT_VERSION" \
    | jq 'del(.metadata.tools.components[].author) | .metadata.authors = [{"name": "anchore"}]' \
    > "$SBOM_OUT" || true
  echo "SBOM generated as a JSON object with name and version fields."
else
  syft . -o cyclonedx-json --source-name "$PROJECT_NAME" --source-version "$PROJECT_VERSION" > "$SBOM_OUT" || true
  echo "jq not found, SBOM may not have name/version fields."
fi

echo "Running tests with coverage..."
go test -v -coverprofile="$COVERAGE_OUT" ./... > "$TEST_OUT" 2>&1 || true

echo
echo "== OSSF Attestation Artifact Summary =="
for f in "$LINT_OUT" "$GOSEC_OUT" "$GOVULN_OUT" "$GITSECRETS_OUT" "$SBOM_OUT" "$TEST_OUT" "$COVERAGE_OUT"; do
  if [[ -f "$f" ]]; then
    echo "  $(ls -lh "$f" | awk '{print $9, $5}')"
    head -5 "$f" | sed 's/^/    /'
    echo "    ..."
  else
    echo "  $f (not generated)"
  fi
done

if [[ -f "$COVERAGE_OUT" ]]; then
  echo
  echo "== Go Test Coverage Summary =="
  go tool cover -func="$COVERAGE_OUT" | grep total: || true
fi

echo "========================================"
echo "Artifacts saved in: $(pwd)"
