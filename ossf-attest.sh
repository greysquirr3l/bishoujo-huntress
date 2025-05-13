#!/bin/bash
# ossf-attest.sh: Run all OSSF Security Baseline checks and save attestation artifacts
set -euo pipefail

# Ensure required tools are installed (macOS/Homebrew or Go-based)
ensure_tool() {
  local tool="$1"
  case "$tool" in
    golangci-lint)
      if ! command -v golangci-lint >/dev/null 2>&1; then
        echo "Installing golangci-lint..."
        if command -v brew >/dev/null 2>&1; then
          brew install golangci-lint || true
        else
          go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        fi
      fi
      ;;
    gosec)
      if ! command -v gosec >/dev/null 2>&1; then
        echo "Installing gosec..."
        go install github.com/securego/gosec/v2/cmd/gosec@latest
      fi
      ;;
    govulncheck)
      if ! command -v govulncheck >/dev/null 2>&1; then
        echo "Installing govulncheck..."
        go install golang.org/x/vuln/cmd/govulncheck@latest
      fi
      ;;
    syft)
      SYFT_VERSION="v1.23.1"
      if ! command -v syft >/dev/null 2>&1 || [[ "$(syft --version 2>/dev/null)" != "syft 1.23.1" ]]; then
        echo "Installing syft $SYFT_VERSION..."
        curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh -s -- -b /usr/local/bin $SYFT_VERSION
      fi
      ;;
    jq)
      if ! command -v jq >/dev/null 2>&1; then
        echo "Installing jq..."
        if command -v brew >/dev/null 2>&1; then
          brew install jq || true
        else
          echo "jq is required but could not be installed automatically. Please install jq manually."
          exit 1
        fi
      fi
      ;;
    git)
      if ! command -v git >/dev/null 2>&1; then
        echo "ERROR: git is required but not installed. Please install git."
        exit 1
      fi
      ;;
    semgrep)
      SEMGR_VERSION="1.119.0"
      if ! command -v semgrep >/dev/null 2>&1 || [[ "$(semgrep --version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')" != "$SEMGR_VERSION" ]]; then
        echo "Installing semgrep $SEMGR_VERSION..."
        # Prefer Homebrew if available, else use pipx
        if command -v brew >/dev/null 2>&1; then
          brew install semgrep || brew upgrade semgrep
        elif command -v pipx >/dev/null 2>&1; then
          pipx install --force semgrep==$SEMGR_VERSION
        else
          echo "Please install semgrep $SEMGR_VERSION manually (https://semgrep.dev/docs/getting-started/)."; exit 1
        fi
      fi
      ;;
    *)
      echo "Unknown tool: $tool"; exit 1;;
  esac
}

# Check and install required tools

# Add semgrep to required tools
for tool in golangci-lint gosec govulncheck git syft jq semgrep; do
  ensure_tool "$tool"
done

# Output files
LINT_OUT=golangci-lint.txt
TEST_OUT=test-results.txt
GOSEC_OUT=gosec.txt
GOVULN_OUT=govulncheck.txt
GITSECRETS_OUT=git-secrets.txt
SBOM_OUT=sbom.json
COVERAGE_OUT=coverage.txt
SEMGR_OUT=semgrep.txt
SEMGR_VERSION="1.119.0"

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

if [[ -n "${SEMGREP_APP_TOKEN:-}" ]]; then
  echo "Logging in to semgrep with SEMGREP_APP_TOKEN..."
  semgrep login --token "$SEMGREP_APP_TOKEN" || true
fi
echo "Running semgrep (SAST)..."
semgrep --config p/owasp-top-ten . > "$SEMGR_OUT" 2>&1 || true &

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
for f in "$LINT_OUT" "$GOSEC_OUT" "$GOVULN_OUT" "$GITSECRETS_OUT" "$SEMGR_OUT" "$SBOM_OUT" "$TEST_OUT" "$COVERAGE_OUT"; do
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
