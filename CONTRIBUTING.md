# Contributing to Bishoujo-Huntress

Thank you for your interest in contributing to Bishoujo-Huntress! We welcome contributions of all kinds—code, documentation, tests, and feedback.

## Getting Started

1. **Fork the repository** and clone your fork locally.
2. **Create a new branch** for your feature or bugfix:

   ```bash
   git checkout -b my-feature
   ```

3. **Install dependencies:**

   ```bash
   go mod download
   ```

4. **Run tests, linters, and generate mocks:**

   ```bash
   make test
   make lint
   make generate-mocks   # Regenerate mocks if you change interfaces
   ./ossf-attest.sh
   ```

## Mocking & Test Doubles

- We use [Mockery](https://github.com/vektra/mockery) for generating mocks for interfaces.
- To regenerate mocks after changing interfaces, run:

  ```bash
  make generate-mocks
  ```

- Mocks are generated in `internal/mocks/` (for internal interfaces) and `pkg/huntress/mocks/` (for public API interfaces).
- See the `//go:generate` directives in interface files for details.

5. **Make your changes.**
6. **Add or update tests** to ensure coverage for your changes.
7. **Commit and push** your branch:

   ```bash
   git add .
   git commit -m "Describe your change"
   git push origin my-feature
   ```

8. **Open a Pull Request** on GitHub. Please describe your change and reference any related issues.

## Code Style & Quality

- Follow Go best practices and idioms.
- All code must pass `make lint` (uses golangci-lint v2).
- All code must pass `make test` and have good test coverage.
- Run `./ossf-attest.sh` before pushing to generate security attestation artifacts.
- Use enums and strict model validation for all API params and models.
- Document all exported symbols with GoDoc comments.

## Security & Compliance

- No hardcoded credentials or secrets.
- All API keys/secrets must be provided via environment variables or config.
- Follow [OSSF Security Baselines](docs/OSSF_SECURITY_BASELINES.md).
- Run secret scanning (`git secrets --scan`) before submitting.
- Generate and include SBOM (`syft . -o cyclonedx-json > sbom.json`).

## Pull Request Checklist

- [ ] All tests pass (`make test`)
- [ ] Lint passes (`make lint`)
- [ ] Security checks pass (`./ossf-attest.sh`)
- [ ] SBOM and attestation artifacts are up to date
- [ ] Documentation is updated (if needed)
- [ ] PR description explains the change and references issues

## Reporting Issues & Security

- For bugs or feature requests, open a GitHub issue.
- For security issues, see [SECURITY.md](SECURITY.md) and do **not** file a public issue.

## Resources

- [Architecture](docs/architecture.md)
- [OSSF Security Baselines](docs/OSSF_SECURITY_BASELINES.md)
- [Project Guidelines](docs/project_guidelines.md)
- [Project Structure](docs/project_structure.md)
- [Implementation Todo List](docs/todo.md)

Thank you for helping make Bishoujo-Huntress better!
