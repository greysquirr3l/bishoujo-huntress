# Security Policy

## Supported Versions

The following versions of Bishoujo-Huntress are currently receiving security updates:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

## Reporting a Vulnerability

The Bishoujo-Huntress team takes security issues seriously. We appreciate your efforts to responsibly disclose your findings and will make every effort to acknowledge your contributions.

To report a security vulnerability, please follow these steps:

1. **DO NOT** create a public GitHub issue for the vulnerability.
2. Email your findings to [s0ma@proton.me](mailto:s0ma@proton.me).
3. Include as much information as possible:
   - A detailed description of the vulnerability
   - Steps to reproduce the issue
   - Potential impact of the vulnerability
   - Suggestions for mitigation or resolution (if any)
4. Allow time for the team to review and respond to your report (typically within 48 hours).

## Security Response Process

After receiving a vulnerability report, the following process will be followed:

1. The maintainer will acknowledge receipt of the report within 48 hours.
2. The maintainer will work to confirm the vulnerability and determine its impact.
3. If confirmed, the maintainer will develop and test a fix.
4. Once a fix is ready, the maintainer will:
   - Release a new version addressing the vulnerability
   - Credit the reporter (if desired)
   - Publish details about the vulnerability including CVE if applicable

## Security Practices

This project follows the [OSSF Security Baselines](https://github.com/ossf/security-baselines) and implements the following security practices:

- Enforcing code reviews for all changes
- Using static analysis tools to detect security issues
- Regular dependency scanning for known vulnerabilities
- Supporting only TLS 1.2+ for all API communications
- Following secure coding practices for Go development

## Security-related Documentation

For more information about our security practices, please refer to:

- [OSSF Security Baselines](docs/OSSF_SECURITY_BASELINES.md)
- [Architecture Security Considerations](docs/architecture.md#security-considerations)

## Acknowledgments

We would like to thank the following individuals for reporting security vulnerabilities:

- (This section will be updated as applicable)
