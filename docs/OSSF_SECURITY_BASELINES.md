# OSSF Security Baselines

This document provides guidance on implementing security practices aligned with the Open Source Security Foundation (OSSF) security baselines for open source projects.

<!-- REF: https://github.com/ossf/security-baselines -->
<!-- REF: https://bestpractices.coreinfrastructure.org/en -->

## 🔐 Core Security Principles

The OSSF security baselines focus on these key areas:

1. **Code Security**
2. **Dependency Management**
3. **Build & Release Security**
4. **Vulnerability Disclosure**
5. **Security Testing**
6. **Documentation**

## 📋 Security Baseline Checklist

### Code Security

- [ ] **Secure Coding Practices**
  - Follow language-specific secure coding guidelines
  - Enforce code quality standards through linters
  - Implement proper error handling
  - Validate all inputs, especially user inputs

- [ ] **Source Control Protection**
  - Protect default branches with required reviews
  - Enforce signed commits
  - Implement branch protection rules
  - Use tools to detect secrets and credentials in code

- [ ] **Authentication & Authorization**
  - Use strong authentication mechanisms
  - Implement proper authorization checks
  - Avoid hardcoded credentials
  - Practice least privilege principles

<!-- REF: https://github.com/ossf/secure-code-fundamentals -->

### Dependency Management

- [ ] **Dependency Verification**
  - Use a dependency scanning tool
  - Verify dependency integrity (checksums)
  - Implement Software Bill of Materials (SBOM)

- [ ] **Dependency Updates**
  - Regularly update dependencies
  - Automate dependency updates when possible
  - Monitor for vulnerabilities in dependencies

- [ ] **Dependency Minimization**
  - Minimize unnecessary dependencies
  - Document why each dependency is needed
  - Prefer well-maintained dependencies

<!-- REF: https://github.com/ossf/package-manager-best-practices -->

### Build & Release Security

- [ ] **Build Reproducibility**
  - Ensure builds are reproducible
  - Document build process thoroughly
  - Use automated builds

- [ ] **Artifact Signing**
  - Sign release artifacts
  - Verify signatures during installation/deployment
  - Document signature verification process

- [ ] **Supply Chain Protection**
  - Use trustworthy build environments
  - Implement CI/CD security controls
  - Consider SLSA (Supply-chain Levels for Software Artifacts) framework

<!-- REF: https://slsa.dev/ -->
<!-- REF: https://github.com/ossf/package-analysis -->

### Vulnerability Disclosure

- [ ] **Security Policy**
  - Maintain a clear SECURITY.md file
  - Define the vulnerability reporting process
  - Document supported versions

- [ ] **Vulnerability Management**
  - Track security issues appropriately
  - Provide timely fixes for security issues
  - Follow coordinated vulnerability disclosure practices

- [ ] **Security Advisories**
  - Publish security advisories for vulnerabilities
  - Use standard formats (e.g., CVE)
  - Communicate impact and mitigation clearly

<!-- REF: https://github.com/ossf/oss-vulnerability-guide -->

### Security Testing

- [ ] **Automated Testing**
  - Implement security-focused test cases
  - Use SAST (Static Application Security Testing) tools
  - Consider DAST (Dynamic Application Security Testing) if applicable

- [ ] **Fuzz Testing**
  - Implement fuzzing for parsing or complex logic
  - Integrate fuzzing into CI pipeline
  - Have a process to triage fuzzing results

- [ ] **Penetration Testing**
  - Consider regular security reviews
  - Document security testing approach
  - Fix identified security issues promptly

<!-- REF: https://github.com/ossf/fuzz-introspector -->
<!-- REF: https://owasp.org/www-project-web-security-testing-guide/ -->

### Security Documentation

- [ ] **User Documentation**
  - Document security features
  - Provide secure configuration guidance
  - Include threat model where appropriate

- [ ] **Developer Documentation**
  - Document security expectations for contributors
  - Provide security testing information
  - Include architecture security considerations

- [ ] **Security Risk Assessment**
  - Identify key security risks
  - Document trust boundaries
  - Maintain security assumptions

## 🛠️ Implementation Guidelines

### Starting Small

1. Begin with basic security hygiene:
   - Enable branch protection
   - Add a SECURITY.md file
   - Set up automated dependency scanning

2. Progress to intermediate measures:
   - Implement automated security testing
   - Sign releases
   - Create a vulnerability management process

3. Advanced security measures:
   - Generate and publish SBOMs
   - Implement fuzzing
   - Conduct regular security audits

## Tool Recommendations

### General Purpose Tools

```bash
# Dependency scanning
$ dependency-check --project "Project Name" --scan /path/to/code

# Secret scanning
$ git-secrets --scan

#### Semgrep (SAST)

This project uses [Semgrep](https://semgrep.dev/) for static application security testing (SAST). The recommended version is **v1.119.0** (pinned in CI and `ossf-attest.sh`).

To run Semgrep locally:

```bash
brew install semgrep  # or pipx install semgrep==1.119.0
semgrep --config p/owasp-top-ten .
```

Semgrep is also run automatically in CI and as part of the OSSF attestation script (`ossf-attest.sh`).

#### SBOM generation

$ syft /path/to/project -o cyclonedx-json > sbom.json

### Language-Specific Tools

Different tools are recommended based on programming language:

- **Go**: gosec, govulncheck
- **JavaScript/Node.js**: npm audit, eslint-plugin-security
- **Python**: bandit, safety
- **Java**: SpotBugs, OWASP Dependency Check
- **Ruby**: Brakeman, bundler-audit
- **Rust**: cargo-audit, cargo-deny

## 📊 Assessment & Improvement

### Measuring Security Maturity

The OSSF provides tools to assess your project's security posture:

- [OSSF Scorecard](https://securityscorecards.dev/): Automated checks for security best practices
- [OpenSSF Best Practices Badge Program](https://bestpractices.coreinfrastructure.org/): A way to show your project follows best practices

### Continuous Improvement

- Regularly review security posture
- Subscribe to security advisories for your dependencies
- Participate in security-focused communities
- Consider having periodic external security reviews

<!-- REF: https://github.com/ossf/scorecard -->
<!-- REF: https://github.com/ossf/allstar -->

## 📚 Additional Resources

- [OSSF Security Insights](https://github.com/ossf/security-insights-spec): Standard format for security information
- [OSSF Security Tooling](https://github.com/ossf/wg-security-tooling): Working group on security tools
- [OSSF Best Practices](https://github.com/ossf/wg-best-practices-os-developers): Best practices for open source developers
- [OWASP Top 10](https://owasp.org/www-project-top-ten/): Common web application security risks
- [CII Best Practices](https://bestpractices.coreinfrastructure.org/): Security best practices for open source

## References

1. [OSSF Security Baselines] (<https://github.com/ossf/security-baselines>)
2. [Open Source Security Foundation] (<https://openssf.org/>)
3. [SLSA Framework] (<https://slsa.dev/>)
4. [OSSF Scorecard] (<https://securityscorecards.dev/>)
5. [OWASP Top 10] (<https://owasp.org/www-project-top-ten/>)
6. [CII Best Practices] (<https://bestpractices.coreinfrastructure.org/>)

## 🧪 How to Test OSSF Security Baselines Locally (macOS)

To verify your adherence to the OSSF Security Baselines, run the following checks locally:

### 1. Static Analysis & Linting

- Run `make lint` to check for code quality, security, and style issues (uses `golangci-lint` with security linters enabled).
- If you do not have `golangci-lint` installed, the Makefile will install it automatically (see README for details).
- Run `gosec ./...` for Go security static analysis.
- Run `govulncheck ./...` for Go vulnerability scanning.

### 2. Dependency Scanning

- Run `go mod tidy` to ensure dependencies are clean and minimal.
- Run `govulncheck ./...` to check for known vulnerabilities in dependencies.

### 3. Secret Scanning

- Install [git-secrets](https://github.com/awslabs/git-secrets) and scan for secrets:

```bash
brew install git-secrets
git secrets --scan
```

### 4. SBOM Generation

This project uses [syft](https://github.com/anchore/syft) **v1.23.1** for SBOM generation. Please use this version for reproducibility and OSSF Scorecard compliance.

Install syft v1.23.1 (recommended):

```bash
curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh -s -- -b /usr/local/bin v1.23.1
syft --version  # should print syft 1.23.1
syft . -o cyclonedx-json > sbom.json
```

### 5. Test Coverage

- Run `make test` to execute all tests and check for high coverage, especially on critical code paths.

### 6. CI/CD Verification

- Ensure your GitHub Actions workflow runs all of the above on every PR.

### 7. Security Policy

- Confirm that [SECURITY.md](../SECURITY.md) is present and up to date.

### 8. Testify Adherence

- Save the output of the above tools (e.g., `gosec`, `govulncheck`, `git secrets`, `syft`) as text files:
  - `gosec.txt`, `govulncheck.txt`, `git-secrets.txt`, `sbom.json`
- Attach these files to a security review or upload as CI artifacts.
- For each checkbox in this document, add a comment or evidence (e.g., filename, summary, or link to artifact) showing how it is satisfied.
- Example:
  - `[x] Secure Coding Practices  # see gosec.txt, govulncheck.txt`
  - `[x] Dependency Verification  # see sbom.json, govulncheck.txt`
  - `[x] Secret Scanning  # see git-secrets.txt`
- Optionally, add a section at the end of this file summarizing the evidence and linking to the artifacts.

---

#### Example: Saving Evidence

```bash
# Run and save outputs
make lint > golangci-lint.txt
make test > test-results.txt
gosec ./... > gosec.txt
govulncheck ./... > govulncheck.txt
git secrets --scan > git-secrets.txt
syft . -o cyclonedx-json > sbom.json
```

---
