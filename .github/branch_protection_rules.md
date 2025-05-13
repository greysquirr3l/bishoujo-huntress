# Branch Protection Rules for OSSF Compliance

This document describes the recommended GitHub branch protection rules for the Bishoujo-Huntress project to comply with the Open Source Security Foundation (OSSF) Scorecard and security best practices.

## Recommended Branch Protection Settings

Apply these rules to your main development branch (e.g., `main` or `master`).

### 1. Require Pull Request Reviews

- Require at least **1 approving review** before merging.
- Dismiss stale pull request approvals when new commits are pushed.
- Require review from Code Owners (if applicable).

### 2. Require Status Checks to Pass Before Merging

- Enable **Require status checks to pass before merging**.
- Select the following required status checks:
  - `ci` (main CI workflow)
  - `scorecard` (OSSF Scorecard workflow)
- Optionally, require additional checks (e.g., `lint`, `test`, `ossf-artifacts`) as needed.
- Do **not** allow bypassing required status checks.

### 3. Require Linear History

- Enable **Require linear history** to prevent merge commits.

### 4. Require Signed Commits

- Enable **Require signed commits** (optional but recommended for OSSF compliance).

### 5. Restrict Who Can Push to Matching Branches

- Restrict direct pushes to the protected branch.
- Allow only maintainers or GitHub Actions bots to push if necessary.

### 6. Require Conversation Resolution

- Enable **Require all conversations to be resolved before merging**.

### 7. Include Administrators (Recommended)

- Apply these rules to administrators for maximum protection.

## How to Configure

1. Go to your repository on GitHub.
2. Click **Settings** > **Branches** > **Branch protection rules**.
3. Click **Add rule** and set the pattern to your main branch (e.g., `main`).
4. Enable the options above as described.
5. Save changes.

## OSSF Scorecard Requirements

- OSSF Scorecard checks for branch protection on the default branch.
- Required status checks must include your main CI and Scorecard workflows.
- See [OSSF Scorecard documentation](https://github.com/ossf/scorecard/blob/main/docs/checks.md#branch-protection) for more details.

## Example

<!--![Branch Protection Example](../docs/img/github-branch-protection-example.png)-->

---

For more information, see:
- [GitHub Docs: About protected branches](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-branch-protection-rules/about-protected-branches)
- [OSSF Scorecard: branch-protection check](https://github.com/ossf/scorecard/blob/main/docs/checks.md#branch-protection)
