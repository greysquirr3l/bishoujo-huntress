# OSSF Security Baseline Attestation Tool

A Go-based implementation of the OSSF (Open Source Security Foundation) Security Baseline attestation process. This tool replaces the bash-based `ossf-attest.sh` script with a more robust, cross-platform solution.

## Features

- **Cross-platform**: Works on Windows, macOS, and Linux
- **Automated tool installation**: Handles tool version management and installation
- **Parallel execution**: Runs security tools concurrently for faster results
- **Structured reporting**: Generates both JSON and human-readable reports
- **Error handling**: Comprehensive error recovery and context
- **Reproducible**: Version-pinned tools with deterministic output

## Security Tools Included

- **golangci-lint** (v2.1.6): Go static analysis and linting
- **gosec** (v2.22.4): Go security scanning
- **govulncheck** (v1.1.4): Go vulnerability detection
- **semgrep** (1.119.0): Static analysis security testing
- **syft** (v1.23.1): Software Bill of Materials (SBOM) generation
- **go test**: Test execution with coverage

## Usage

### Basic Usage

```bash
# Build the tool
make ossf-build

# Run attestation with all checks
./build/ossf-attest

# Run with verbose output
./build/ossf-attest --verbose

# Run tools sequentially instead of parallel
./build/ossf-attest --sequential

# Custom output directory
./build/ossf-attest --output ./attestation-results
```

### Environment Variables

- `PROJECT_NAME`: Override project name (default: "bishoujo-huntress")
- `PROJECT_VERSION`: Override project version (default: from git tags)
- `SEMGREP_APP_TOKEN`: Semgrep authentication token (optional)

### Make Integration

The tool is integrated into the project Makefile:

```bash
# Run OSSF attestation
make ossf-attest

# Build the OSSF tool
make ossf-build

# Clean OSSF artifacts
make ossf-clean
```

## Output Files

The tool generates several output files:

- `ossf-attestation-report.json`: Complete attestation results in JSON format
- `ossf-attestation-summary.txt`: Human-readable summary
- `golangci-lint.txt`: Linter output
- `gosec.txt`: Security scanner output
- `govulncheck.txt`: Vulnerability check output
- `semgrep.txt`: Static analysis output
- `test-results.txt`: Test execution output
- `coverage.txt`: Test coverage data
- `sbom.json`: Software Bill of Materials

## Tool Installation

The tool automatically handles installation of required security tools:

1. **Local installation**: Tools are installed to `./.local_tools/bin/`
2. **PATH management**: Automatically adds local tools to PATH
3. **Version checking**: Verifies correct tool versions are installed
4. **Fallback strategies**: Multiple installation methods for each tool

### Manual Tool Setup

If you prefer to install tools manually:

```bash
# golangci-lint
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.6

# gosec
go install github.com/securego/gosec/v2/cmd/gosec@v2.22.4

# govulncheck
go install golang.org/x/vuln/cmd/govulncheck@v1.1.4

# syft
curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.23.1

# semgrep
pipx install semgrep==1.119.0
# or
pip3 install --user semgrep==1.119.0
```

## Report Format

### JSON Report

The JSON report includes:

```json
{
  "project_name": "bishoujo-huntress",
  "project_version": "v1.0.0",
  "timestamp": "2024-01-15T10:30:00Z",
  "go_version": "go1.24.4",
  "os": "darwin",
  "arch": "arm64",
  "results": {
    "golangci-lint": {
      "tool": "golangci-lint",
      "version": "1.61.0",
      "success": true,
      "output": "...",
      "duration": "15.2s",
      "output_file": "golangci-lint.txt"
    }
  },
  "summary": {
    "total_tools": 6,
    "success_count": 6,
    "failure_count": 0,
    "warning_count": 0,
    "total_duration": "45.8s"
  }
}
```

### Text Summary

Human-readable summary with pass/fail status and timing information.

## Integration with CI/CD

### GitHub Actions

```yaml
- name: Run OSSF Security Baseline Attestation
  run: |
    make ossf-attest

- name: Upload Attestation Reports
  uses: actions/upload-artifact@v4
  with:
    name: ossf-attestation-reports
    path: |
      ossf-attestation-*.json
      ossf-attestation-*.txt
      *.txt
      sbom.json
```

### Pre-commit Hooks

Add to `.pre-commit-config.yaml`:

```yaml
- repo: local
  hooks:
    - id: ossf-attest
      name: OSSF Security Baseline Attestation
      entry: make ossf-attest
      language: system
      pass_filenames: false
```

## Troubleshooting

### Tool Installation Issues

1. **golangci-lint**: Ensure curl is available and GitHub is accessible
2. **gosec/govulncheck**: Requires Go 1.21+ and GitHub proxy access
3. **syft**: Requires curl and GitHub releases access
4. **semgrep**: Requires Python 3.8+ and pip/pipx

### Permission Issues

```bash
# Ensure local tools directory is writable
chmod 755 ./.local_tools/bin/

# Fix PATH issues
export PATH="$PWD/.local_tools/bin:$PATH"
```

### Network Issues

Some tools require internet access for:
- Downloading tool binaries
- Fetching vulnerability databases
- Accessing rule repositories

Configure proxy settings if needed:

```bash
export HTTP_PROXY=http://proxy.example.com:8080
export HTTPS_PROXY=http://proxy.example.com:8080
```

## Comparison with Bash Script

| Feature | Bash Script | Go Implementation |
|---------|-------------|-------------------|
| Cross-platform | ❌ Unix only | ✅ Windows/macOS/Linux |
| Error handling | ⚠️ Basic | ✅ Comprehensive |
| Parallel execution | ⚠️ Limited | ✅ Full control |
| Tool management | ⚠️ Manual | ✅ Automated |
| Structured output | ❌ Text only | ✅ JSON + Text |
| Testing | ❌ Difficult | ✅ Unit testable |
| Dependencies | ⚠️ Many shell tools | ✅ Self-contained |
| Maintenance | ⚠️ Script complexity | ✅ Type-safe Go |

## Security Considerations

- Tools are installed to local directory (not system-wide)
- Version pinning prevents supply chain attacks
- No hardcoded credentials or secrets
- Secure file permissions (0644 for outputs, 0755 for executables)
- Input validation for all file paths and commands
- Error messages don't leak sensitive information

## Development

### Building

```bash
go build -o build/ossf-attest ./cmd/ossf-attest
```

### Testing

```bash
go test ./cmd/ossf-attest/...
```

### Linting

```bash
golangci-lint run ./cmd/ossf-attest/...
```

## Migration from Bash Script

The Go implementation provides the same functionality as the original `ossf-attest.sh` script:

1. **Same tools**: All security tools from the bash script
2. **Same outputs**: Compatible file formats and names
3. **Enhanced features**: Better error handling, parallel execution, structured reporting
4. **Easier maintenance**: Type-safe Go code instead of shell scripting

To migrate:

1. Update Makefile target to use Go tool
2. Update CI/CD workflows to use new binary
3. Remove old bash script dependencies
4. Update documentation references
