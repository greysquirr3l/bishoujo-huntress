#!/bin/bash
# ossf-attest.sh: Run all OSSF Security Baseline checks and save attestation artifacts
set -euo pipefail

# Ensure required tools are installed (macOS/Homebrew or Go-based)
ensure_tool() {
  local tool="$1"
  case "$tool" in
    golangci-lint)
      if ! command -v golangci-lint >/dev/null 2>&1; then
        echo "Installing golangci-lint v2.1.6..."
        GOLANGCI_LINT_VERSION="v2.1.6"
        curl -sSfL -o golangci-lint.tar.gz "https://github.com/golangci/golangci-lint/releases/download/${GOLANGCI_LINT_VERSION}/golangci-lint-2.1.6-linux-amd64.tar.gz"
        tar -xzf golangci-lint.tar.gz
        sudo mv golangci-lint-2.1.6-linux-amd64/golangci-lint /usr/local/bin/
        rm -rf golangci-lint.tar.gz golangci-lint-2.1.6-linux-amd64
      fi
      ;;
    gosec)
      if ! command -v gosec >/dev/null 2>&1; then
        echo "Installing gosec v2.19.0..."
        go install github.com/securego/gosec/v2/cmd/gosec@v2.19.0
      fi
      ;;
    govulncheck)
      if ! command -v govulncheck >/dev/null 2>&1; then
        echo "Installing govulncheck v1.1.4..."
        go install golang.org/x/vuln/cmd/govulncheck@v1.1.4
      fi
      ;;
    syft)
      SYFT_VERSION="v1.23.1"
      if ! command -v syft >/dev/null 2>&1 || [[ "$(syft --version 2>/dev/null)" != "syft 1.23.1" ]]; then
        echo "Installing syft $SYFT_VERSION..."
        curl -sSfL -o syft.tar.gz "https://github.com/anchore/syft/releases/download/${SYFT_VERSION}/syft_${SYFT_VERSION#v}_linux_amd64.tar.gz"
        tar -xzf syft.tar.gz
        sudo mv syft /usr/local/bin/
        rm -rf syft.tar.gz syft
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
        unameOut="$(uname -s)"
        archOut="$(uname -m)"
        if [[ "$unameOut" == "Darwin" ]]; then
          if [[ "$archOut" == "arm64" ]]; then
            SEMGREP_URL="https://github.com/returntocorp/semgrep/releases/download/v${SEMGR_VERSION}/semgrep_${SEMGR_VERSION}_macos_arm64.tar.gz"
          else
            SEMGREP_URL="https://github.com/returntocorp/semgrep/releases/download/v${SEMGR_VERSION}/semgrep_${SEMGR_VERSION}_macos_x86_64.tar.gz"
          fi
        elif [[ "$unameOut" == "Linux" ]]; then
          if [[ "$archOut" == "x86_64" || "$archOut" == "amd64" ]]; then
            SEMGREP_URL="https://github.com/returntocorp/semgrep/releases/download/v${SEMGR_VERSION}/semgrep_${SEMGR_VERSION}_linux_x86_64.tar.gz"
          elif [[ "$archOut" == "aarch64" || "$archOut" == "arm64" ]]; then
            SEMGREP_URL="https://github.com/returntocorp/semgrep/releases/download/v${SEMGR_VERSION}/semgrep_${SEMGR_VERSION}_linux_arm64.tar.gz"
          else
            echo "Unsupported Linux architecture: $archOut"; exit 1
          fi
        else
          echo "Unsupported OS: $unameOut"; exit 1
        fi
        curl -sSfL -o semgrep.tar.gz "$SEMGREP_URL"
        tar -xzf semgrep.tar.gz
        sudo mv semgrep/semgrep /usr/local/bin/
        rm -rf semgrep.tar.gz semgrep
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
