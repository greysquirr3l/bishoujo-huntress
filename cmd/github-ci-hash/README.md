# GitHub CI Hash Updater

A comprehensive Go tool for managing GitHub Actions in your CI/CD workflows. This tool automatically checks for updates, resolves the latest stable releases, fetches commit SHAs, and updates your workflow files with proper SHA pinning for enhanced security.

## Features

- 🔍 **Check for Updates**: Scan all workflow files and identify actions with available updates
- 🔄 **Update with Confirmation**: Update actions to latest versions with user confirmation
- 🔒 **SHA Verification**: Verify all actions are properly pinned to commit SHAs
- 🎯 **Selective Updates**: Update all workflows or target specific workflow files
- 📊 **Detailed Reports**: Comprehensive summaries of action status and available updates
- 🛡️ **Security Focused**: Follows OSSF security best practices with SHA pinning

## Installation

The tool is built automatically when using the Makefile targets, or you can build it manually:

```bash
go build -o ./build/github-ci-hash ./cmd/github-ci-hash/
```

## Usage

### Using Makefile (Recommended)

```bash
# Check for action updates
make ci-hash-check

# Update actions with confirmation prompts
make ci-hash-update

# Verify all actions are pinned to SHAs
make ci-hash-verify
```

### Direct Usage

```bash
# Check for updates without applying
./build/github-ci-hash check

# Update all workflows (with confirmation)
./build/github-ci-hash update

# Update specific workflow file
./build/github-ci-hash update ci.yml

# Verify all actions are pinned to SHAs
./build/github-ci-hash verify
```

## Authentication

For higher API rate limits, set a GitHub token:

```bash
export GITHUB_TOKEN="your_github_token"
# or
export GH_TOKEN="your_github_token"
```

Without authentication, you'll hit GitHub's rate limits faster but the tool will still work for smaller projects.

## Dependency Graph Integration

The tool can leverage GitHub's dependency graph APIs when authenticated:

1. **Automated Detection**: Automatically discovers all GitHub Actions in your workflows
2. **Latest Release Fetching**: Uses GitHub API to get the latest stable releases
3. **SHA Resolution**: Resolves tags and branches to commit SHAs
4. **Special Handling**: Proper handling for complex actions like CodeQL bundles

## Security Benefits

- **SHA Pinning**: Ensures all actions are pinned to specific commit SHAs
- **Supply Chain Security**: Prevents attacks via compromised action tags
- **OSSF Compliance**: Follows OpenSSF Scorecard security recommendations
- **Verification**: Built-in verification to ensure proper pinning

## Example Output

```bash
🔍 Scanning workflow files...
Checking for action updates...

📁 .github/workflows/ci.yml:
  🔍 Checking actions/checkout... ✅ Up to date (v4.2.2)
  🔍 Checking step-security/harden-runner... 🔄 Update available: v2.12.0 → v2.12.1
  🔍 Checking actions/setup-go... ✅ Up to date (v5.5.0)

📊 Summary:
📈 Total: 15 actions
✅ Up to date: 13
🔄 Need updates: 2
```

## Integration Options

### Pre-commit Hook

Add to `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: local
    hooks:
      - id: github-ci-hash-verify
        name: Verify GitHub Actions are pinned to SHAs
        entry: make ci-hash-verify
        language: system
        files: '^\.github/workflows/.*\.ya?ml$'
        pass_filenames: false

      - id: ossf-security-check
        name: OSSF Security Baseline Check
        entry: make ossf-attest
        language: system
        pass_filenames: false
```

### GitHub Workflow

Create `.github/workflows/security-compliance.yml`:

```yaml
name: Security Compliance Check
on:
  schedule:
    - cron: '0 0 * * 1' # Weekly on Monday
  workflow_dispatch:
  pull_request:
    paths:
      - '.github/workflows/**'

jobs:
  security-compliance:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      security-events: write
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      - name: Check GitHub Actions are pinned
        run: make ci-hash-verify
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      - name: Run OSSF Security Baseline
        run: make ossf-attest
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      - name: Check for action updates
        run: make ci-hash-check
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

## Related Tools

This tool is part of a comprehensive security toolkit in the Bishoujo-Huntress project:

### OSSF Security Baseline Attestation Tool

The project also includes an **OSSF Security Baseline Attestation Tool** (`ossf-attest`) that provides comprehensive security scanning and compliance checking:

```bash
# Run OSSF security baseline checks
make ossf-attest

# Run with verbose output
./build/ossf-attest -v

# Run tools sequentially instead of parallel
./build/ossf-attest -s

# Save reports to custom directory
./build/ossf-attest -o security-reports/
```

The OSSF tool performs:
- Static analysis with `gosec` and `golangci-lint`
- Vulnerability scanning with `govulncheck`
- Dependency analysis and SBOM generation
- Secret scanning with `git-secrets`
- License compliance checking
- Security policy validation

### Integrated Workflow

Both tools work together to ensure comprehensive supply chain security:

1. **GitHub CI Hash Updater**: Ensures GitHub Actions are pinned to commit SHAs
2. **OSSF Attestation Tool**: Validates overall security posture and compliance

```bash
# Complete security workflow
make ci-hash-verify    # Verify action SHA pinning
make ossf-attest       # Run OSSF security baseline checks
make security-check    # Additional security scans
```

## Advanced Features

### Special Action Handling

- **CodeQL Actions**: Automatically handles CodeQL bundle versioning
- **Sub-actions**: Properly resolves SHAs for sub-actions like `github/codeql-action/upload-sarif`
- **Version Normalization**: Handles different version formats consistently
- **Idempotent Updates**: Checks existing file state and only modifies files when actual changes are needed
- **Atomic Backup/Restore**: Creates backups before updates and can restore on failure

### Error Handling

- **Retry Logic**: Built-in retry with exponential backoff for API calls
- **Graceful Degradation**: Works with or without authentication
- **Detailed Error Messages**: Clear error reporting for troubleshooting
- **Rollback Support**: Automatic restoration from backup files on update failures

## Contributing

This tool is part of the Bishoujo-Huntress project. See the main project documentation for contribution guidelines.

## License

Same as the main Bishoujo-Huntress project.
