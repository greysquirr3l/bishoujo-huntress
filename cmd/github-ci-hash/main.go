// Package main implements a GitHub CI Hash Updater tool for managing GitHub Actions
// in CI/CD workflows with proper SHA pinning for enhanced security.
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"
)

const (
	// codeQLAction is the GitHub CodeQL action repository name
	codeQLAction = "codeql-action"
)

var (
	// shaRegex is a compiled regex for matching 40-character SHA hashes
	shaRegex = regexp.MustCompile(`^[a-f0-9]{40}$`)

	// Version information (set by build flags)
	// Version is the current version of the application
	Version = "dev"
	// GitCommit is the git commit hash the binary was built from
	GitCommit = "unknown"
	// BuildTime is the time when the binary was built
	BuildTime = "unknown"
)

// ActionInfo represents information about a GitHub Action
type ActionInfo struct {
	Repo         string `json:"repo"`
	CurrentRef   string `json:"current_ref"`
	CurrentSHA   string `json:"current_sha"`
	LatestTag    string `json:"latest_tag"`
	LatestSHA    string `json:"latest_sha"`
	NeedsUpdate  bool   `json:"needs_update"`
	Line         int    `json:"line"`
	OriginalLine string `json:"original_line"`
	WorkflowFile string `json:"workflow_file"`
}

// WorkflowActions represents all actions found in workflows
type WorkflowActions map[string][]ActionInfo

// GitHubClient wraps the GitHub API client with additional functionality
type GitHubClient struct {
	client *github.Client
	ctx    context.Context
}

// NewGitHubClient creates a new GitHub client with optional authentication
func NewGitHubClient() *GitHubClient {
	ctx := context.Background()
	var client *github.Client

	// Try to use GitHub token from environment
	if token, source := getGitHubToken(); token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)

		// Show green status indicator for authenticated access
		fmt.Printf("🟢 GitHub API: \033[32mAuthenticated\033[0m via %s (higher rate limits available)\n", source)
	} else {
		client = github.NewClient(nil)
		fmt.Printf("🟡 GitHub API: \033[33mUnauthenticated\033[0m (lower rate limits)\n")
		fmt.Println("   Set GITHUB_TOKEN or GH_TOKEN environment variable, or authenticate with 'gh auth login'.")
	}

	return &GitHubClient{
		client: client,
		ctx:    ctx,
	}
}

// getGitHubToken retrieves GitHub token from environment variables or gh CLI
func getGitHubToken() (string, string) {
	// Try environment variables first
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token, "GITHUB_TOKEN"
	}
	if token := os.Getenv("GH_TOKEN"); token != "" {
		return token, "GH_TOKEN"
	}

	// Try to get token from gh CLI if available
	if token := getTokenFromGHCLI(); token != "" {
		return token, "gh CLI"
	}

	return "", ""
}

// getTokenFromGHCLI attempts to get the GitHub token from gh CLI
func getTokenFromGHCLI() string {
	cmd := exec.Command("gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		// gh CLI not available or not authenticated
		return ""
	}

	token := strings.TrimSpace(string(output))
	if token != "" {
		return token
	}

	return ""
}

// GetLatestRelease fetches the latest release for a repository
func (gc *GitHubClient) GetLatestRelease(owner, repo string) (*github.RepositoryRelease, error) {
	release, _, err := gc.client.Repositories.GetLatestRelease(gc.ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest release for %s/%s: %w", owner, repo, err)
	}
	return release, nil
}

// ResolveSHA resolves a tag or branch to its commit SHA
func (gc *GitHubClient) ResolveSHA(owner, repo, ref string) (string, error) {
	// Special handling for CodeQL action bundle tags
	if owner == "github" && repo == codeQLAction && strings.HasPrefix(ref, "v") {
		ref = "codeql-bundle-" + ref
	}

	// Try to get tag first
	gitRef, _, err := gc.client.Git.GetRef(gc.ctx, owner, repo, "tags/"+ref)
	if err == nil && gitRef.Object != nil {
		if gitRef.Object.GetType() == "tag" {
			// Dereference annotated tag
			tag, _, tagErr := gc.client.Git.GetTag(gc.ctx, owner, repo, gitRef.Object.GetSHA())
			if tagErr == nil && tag.Object != nil {
				return tag.Object.GetSHA(), nil
			}
		}
		return gitRef.Object.GetSHA(), nil
	}

	// Try branch if tag fails
	gitRef, _, err = gc.client.Git.GetRef(gc.ctx, owner, repo, "heads/"+ref)
	if err == nil && gitRef.Object != nil {
		return gitRef.Object.GetSHA(), nil
	}

	return "", fmt.Errorf("could not resolve ref %s for %s/%s", ref, owner, repo)
}

// parseWorkflowFile parses a workflow file and extracts GitHub Actions
func parseWorkflowFile(filename string) ([]ActionInfo, error) {
	content, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		return nil, fmt.Errorf("failed to read workflow file %s: %w", filename, err)
	}

	var actions []ActionInfo
	lines := strings.Split(string(content), "\n")

	// Regex to match uses: statements
	usesRegex := regexp.MustCompile(`^\s*uses:\s+([^@]+)@([a-f0-9]{40}|[^#\s]+)(?:\s*#\s*([^\s]+))?`)

	for i, line := range lines {
		matches := usesRegex.FindStringSubmatch(line)
		if matches != nil {
			repo := matches[1]
			currentRef := matches[2]
			// comment := "" // Available for future use
			// if len(matches) > 3 {
			// 	comment = matches[3]
			// }

			// Determine current SHA (if ref is already a SHA)
			currentSHA := ""
			if shaRegex.MatchString(currentRef) {
				currentSHA = currentRef
			}

			actions = append(actions, ActionInfo{
				Repo:         repo,
				CurrentRef:   currentRef,
				CurrentSHA:   currentSHA,
				Line:         i + 1,
				OriginalLine: line,
				WorkflowFile: filename,
			})
		}
	}

	return actions, nil
}

// scanWorkflows scans all workflow files and extracts GitHub Actions
func scanWorkflows() (WorkflowActions, error) {
	workflowActions := make(WorkflowActions)

	workflowDir := ".github/workflows"
	entries, err := os.ReadDir(workflowDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read workflow directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !strings.HasSuffix(filename, ".yml") && !strings.HasSuffix(filename, ".yaml") {
			continue
		}

		fullPath := filepath.Join(workflowDir, filename)
		actions, err := parseWorkflowFile(fullPath)
		if err != nil {
			fmt.Printf("Warning: Failed to parse %s: %v\n", fullPath, err)
			continue
		}

		if len(actions) > 0 {
			workflowActions[fullPath] = actions
		}
	}

	return workflowActions, nil
}

// checkForUpdates checks if actions have newer versions available
func checkForUpdates(gc *GitHubClient, actions WorkflowActions) {
	fmt.Println("Checking for action updates...")

	for workflow, actionList := range actions {
		fmt.Printf("\n📁 %s:\n", workflow)

		for i := range actionList {
			action := &actionList[i]

			// Parse owner/repo from action repo
			parts := strings.Split(action.Repo, "/")
			if len(parts) < 2 {
				fmt.Printf("  ⚠️  Invalid repo format: %s\n", action.Repo)
				continue
			}

			owner := parts[0]
			repo := parts[1]

			// For sub-actions (like github/codeql-action/upload-sarif), use the main repo
			if len(parts) > 2 && owner == "github" && repo == codeQLAction {
				// Keep the original repo path but fetch from main repo
				repo = codeQLAction
			}

			fmt.Printf("  🔍 Checking %s...", action.Repo)

			// Get latest release
			release, err := gc.GetLatestRelease(owner, repo)
			if err != nil {
				fmt.Printf(" ❌ Error: %v\n", err)
				continue
			}

			action.LatestTag = release.GetTagName()

			// Resolve SHA for latest tag
			sha, err := gc.ResolveSHA(owner, repo, action.LatestTag)
			if err != nil {
				fmt.Printf(" ❌ Error resolving SHA: %v\n", err)
				continue
			}

			action.LatestSHA = sha

			// Check if update is needed
			if action.CurrentSHA == "" {
				// Current ref is not a SHA, resolve it
				currentSHA, err := gc.ResolveSHA(owner, repo, action.CurrentRef)
				if err != nil {
					fmt.Printf(" ❌ Error resolving current SHA: %v\n", err)
					continue
				}
				action.CurrentSHA = currentSHA
			}

			if action.CurrentSHA != action.LatestSHA {
				action.NeedsUpdate = true
				fmt.Printf(" 🔄 Update available: %s → %s\n", action.CurrentRef, action.LatestTag)
			} else {
				fmt.Printf(" ✅ Up to date (%s)\n", action.LatestTag)
			}
		}

		// Update the slice in the map
		actions[workflow] = actionList
	}
}

// promptForConfirmation asks user for confirmation
func promptForConfirmation(message string) bool {
	fmt.Printf("%s (y/N): ", message)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

// updateWorkflowFile updates a workflow file with new action versions
// This function is idempotent - it can be called multiple times safely
// and will only make changes when actually needed
func updateWorkflowFile(filename string, actions []ActionInfo) error {
	content, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	// Check if any updates are actually needed (idempotent check)
	hasActualUpdates := false
	for _, action := range actions {
		if !action.NeedsUpdate {
			continue
		}

		lineIndex := action.Line - 1
		if lineIndex >= len(lines) {
			continue
		}

		// Check if the line already has the target SHA
		currentLine := lines[lineIndex]
		expectedLine := regexp.MustCompile(`@[a-f0-9]{40}|@[^#\s]+`).ReplaceAllString(currentLine, fmt.Sprintf("@%s # %s", action.LatestSHA, action.LatestTag))
		if currentLine != expectedLine {
			hasActualUpdates = true
			break
		}
	}

	// If no actual updates needed, return early (idempotent behavior)
	if !hasActualUpdates {
		fmt.Printf("  ✅ %s: Already up to date, no changes needed\n", filename)
		return nil
	}

	// Sort actions by line number in reverse order to avoid line number shifting
	sort.Slice(actions, func(i, j int) bool {
		return actions[i].Line > actions[j].Line
	})

	for _, action := range actions {
		if !action.NeedsUpdate {
			continue
		}

		lineIndex := action.Line - 1
		if lineIndex >= len(lines) {
			continue
		}

		// Replace the line with updated SHA and tag comment
		oldLine := lines[lineIndex]
		newLine := regexp.MustCompile(`@[a-f0-9]{40}|@[^#\s]+`).ReplaceAllString(oldLine, fmt.Sprintf("@%s # %s", action.LatestSHA, action.LatestTag))

		// Only update if actually different (additional idempotent check)
		if oldLine != newLine {
			lines[lineIndex] = newLine
			fmt.Printf("  📝 Updated line %d: %s → %s\n", action.Line, action.CurrentRef, action.LatestTag)
		}
	}

	// Write back to file
	newContent := strings.Join(lines, "\n")
	return os.WriteFile(filename, []byte(newContent), 0600)
}

// updateActions updates the workflow files with new action versions
// This function implements atomic update semantics:
// - Creates backups before any modifications
// - Rolls back changes if any operation fails
// - Is idempotent and safe to retry
func updateActions(actions WorkflowActions, targetWorkflow string) error {
	fmt.Println("\n🚀 Updating workflow files...")

	// Collect files that need updates for atomic-like behavior
	var filesToUpdate []string
	for workflow, actionList := range actions {
		// If specific workflow is targeted, skip others
		if targetWorkflow != "" && workflow != targetWorkflow {
			continue
		}

		// Check if any actions need updates
		hasUpdates := false
		for _, action := range actionList {
			if action.NeedsUpdate {
				hasUpdates = true
				break
			}
		}

		if hasUpdates {
			filesToUpdate = append(filesToUpdate, workflow)
		}
	}

	if len(filesToUpdate) == 0 {
		fmt.Println("  ✅ No updates needed for any workflow files")
		return nil
	}

	// Create all backups first (atomic preparation)
	backupFiles := make(map[string]string)
	for _, workflow := range filesToUpdate {
		// Create backup with deterministic name
		backupFile := workflow + ".bak"
		if err := copyFile(workflow, backupFile); err != nil {
			// Clean up any backups we've already created
			for _, existingBackup := range backupFiles {
				if removeErr := os.Remove(existingBackup); removeErr != nil {
					fmt.Printf("Warning: failed to clean up backup %s: %v\n", existingBackup, removeErr)
				}
			}
			return fmt.Errorf("failed to create backup for %s: %w", workflow, err)
		}
		backupFiles[workflow] = backupFile
		fmt.Printf("  💾 Created backup: %s\n", backupFile)
	}

	// Now process each workflow with atomic rollback capability
	for workflow, actionList := range actions {
		// If specific workflow is targeted, skip others
		if targetWorkflow != "" && workflow != targetWorkflow {
			continue
		}

		// Check if any actions need updates
		hasUpdates := false
		for _, action := range actionList {
			if action.NeedsUpdate {
				hasUpdates = true
				break
			}
		}

		if !hasUpdates {
			fmt.Printf("  ✅ %s: No updates needed\n", workflow)
			continue
		}

		fmt.Printf("\n📁 %s:\n", workflow)

		// Show what will be updated
		for _, action := range actionList {
			if action.NeedsUpdate {
				fmt.Printf("  🔄 %s: %s → %s (%s)\n", action.Repo, action.CurrentRef, action.LatestTag, action.LatestSHA[:8])
			}
		}

		// Ask for confirmation
		if !promptForConfirmation(fmt.Sprintf("Update %s?", workflow)) {
			fmt.Printf("  ⏭️  Skipped %s\n", workflow)
			continue
		}

		// Update the file (now with idempotent checks)
		if err := updateWorkflowFile(workflow, actionList); err != nil {
			fmt.Printf("  ❌ Failed to update: %v\n", err)

			// Restore from backup on failure
			if backupFile, exists := backupFiles[workflow]; exists {
				if restoreErr := copyFile(backupFile, workflow); restoreErr != nil {
					fmt.Printf("  ❌ Failed to restore backup: %v\n", restoreErr)
				} else {
					fmt.Printf("  🔄 Restored from backup due to update failure\n")
				}
			}
			continue
		}

		fmt.Printf("  ✅ Updated %s\n", workflow)
	}

	return nil
}

// copyFile copies a file
func copyFile(src, dst string) error {
	source, err := os.Open(filepath.Clean(src))
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := source.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close source file: %v\n", closeErr)
		}
	}()

	destination, err := os.Create(filepath.Clean(dst))
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := destination.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close destination file: %v\n", closeErr)
		}
	}()

	_, err = io.Copy(destination, source)
	return err
}

// printSummary prints a summary of actions and their status
func printSummary(actions WorkflowActions) {
	fmt.Println("\n📊 Summary:")

	totalActions := 0
	upToDate := 0
	needsUpdate := 0

	for workflow, actionList := range actions {
		fmt.Printf("\n📁 %s:\n", workflow)

		for _, action := range actionList {
			totalActions++
			status := "✅ Up to date"
			if action.NeedsUpdate {
				needsUpdate++
				status = "🔄 Update available"
			} else {
				upToDate++
			}

			fmt.Printf("  %s: %s (%s)\n", action.Repo, status, action.LatestTag)
		}
	}

	fmt.Printf("\n📈 Total: %d actions\n", totalActions)
	fmt.Printf("✅ Up to date: %d\n", upToDate)
	fmt.Printf("🔄 Need updates: %d\n", needsUpdate)
}

// verifyPinnedSHAs verifies that all actions are pinned to SHAs
func verifyPinnedSHAs() error {
	fmt.Println("\n🔒 Verifying all actions are pinned to SHAs...")

	actions, err := scanWorkflows()
	if err != nil {
		return err
	}

	unpinned := []string{}

	for workflow, actionList := range actions {
		for _, action := range actionList {
			if !shaRegex.MatchString(action.CurrentRef) {
				unpinned = append(unpinned, fmt.Sprintf("%s:%d %s@%s", workflow, action.Line, action.Repo, action.CurrentRef))
			}
		}
	}

	if len(unpinned) > 0 {
		fmt.Println("❌ The following actions are not pinned to SHAs:")
		for _, item := range unpinned {
			fmt.Printf("  %s\n", item)
		}
		return fmt.Errorf("found %d unpinned actions", len(unpinned))
	}

	fmt.Println("✅ All actions are properly pinned to SHAs")
	return nil
}

// installPreCommitHooks installs pre-commit hooks for the repository
func installPreCommitHooks() error {
	fmt.Println("🔧 Installing pre-commit hooks...")

	// Check if we're in a git repository
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return fmt.Errorf("not in a git repository (no .git directory found)")
	}

	// Create hooks directory if it doesn't exist
	hooksDir := ".git/hooks"
	if err := os.MkdirAll(hooksDir, 0750); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	// Pre-commit hook script
	preCommitHook := `#!/bin/sh
# Pre-commit hook for github-ci-hash project
set -e

echo "🔍 Running pre-commit checks..."

# Check if golangci-lint is available
if ! command -v golangci-lint >/dev/null 2>&1; then
    echo "❌ golangci-lint is not installed"
    echo "   Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    exit 1
fi

# Run linting
echo "🔍 Running golangci-lint..."
if ! golangci-lint run; then
    echo "❌ Linting failed"
    exit 1
fi

# Run tests
echo "🧪 Running tests..."
if ! go test ./...; then
    echo "❌ Tests failed"
    exit 1
fi

# Verify all GitHub Actions are pinned to SHAs
echo "🔒 Verifying GitHub Actions are pinned to SHAs..."
if ! go run . verify >/dev/null 2>&1; then
    echo "❌ Some GitHub Actions are not pinned to SHAs"
    echo "   Run 'go run . verify' to see details"
    exit 1
fi

echo "✅ All pre-commit checks passed!"
`

	// Write pre-commit hook
	preCommitPath := filepath.Join(hooksDir, "pre-commit")
	// #nosec G306 - Git hooks must be executable (0755) to function properly
	if err := os.WriteFile(preCommitPath, []byte(preCommitHook), 0755); err != nil {
		return fmt.Errorf("failed to write pre-commit hook: %w", err)
	}

	fmt.Printf("✅ Pre-commit hook installed at %s\n", preCommitPath)

	// Pre-push hook script
	prePushHook := `#!/bin/sh
# Pre-push hook for github-ci-hash project
set -e

echo "🚀 Running pre-push checks..."

# Check for GitHub Actions updates
echo "🔍 Checking for GitHub Action updates..."
if ! go run . check >/dev/null 2>&1; then
    echo "⚠️  Warning: Could not check for GitHub Action updates"
    echo "   This might be due to API rate limits or network issues"
fi

echo "✅ Pre-push checks completed!"
`

	// Write pre-push hook
	prePushPath := filepath.Join(hooksDir, "pre-push")
	// #nosec G306 - Git hooks must be executable (0755) to function properly
	if err := os.WriteFile(prePushPath, []byte(prePushHook), 0755); err != nil {
		return fmt.Errorf("failed to write pre-push hook: %w", err)
	}

	fmt.Printf("✅ Pre-push hook installed at %s\n", prePushPath)

	fmt.Println("\n🎉 Pre-commit hooks successfully installed!")
	fmt.Println("\nThe following hooks are now active:")
	fmt.Println("📋 pre-commit: Runs linting, tests, and SHA verification")
	fmt.Println("🚀 pre-push: Checks for GitHub Action updates")
	fmt.Println("\nTo bypass hooks (not recommended): git commit --no-verify")

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("GitHub CI Hash Updater")
		fmt.Printf("Version: %s (commit: %s, built: %s)\n", Version, GitCommit, BuildTime)
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  github-ci-hash check                    - Check for updates without applying")
		fmt.Println("  github-ci-hash update                   - Update all workflows (with confirmation)")
		fmt.Println("  github-ci-hash update <workflow-file>   - Update specific workflow file")
		fmt.Println("  github-ci-hash verify                   - Verify all actions are pinned to SHAs")
		fmt.Println("  github-ci-hash install-hooks            - Install pre-commit hooks")
		fmt.Println("  github-ci-hash version                  - Show version information")
		fmt.Println("")
		fmt.Println("Environment variables:")
		fmt.Println("  GITHUB_TOKEN or GH_TOKEN - GitHub API token for higher rate limits")
		fmt.Println("  (or authenticate with 'gh auth login' to use gh CLI token)")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "version":
		fmt.Printf("GitHub CI Hash Updater\n")
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Build Time: %s\n", BuildTime)
		fmt.Printf("Go Version: %s\n", strings.TrimPrefix(runtime.Version(), "go"))
		return

	case "check":
		gc := NewGitHubClient()

		fmt.Println("🔍 Scanning workflow files...")
		actions, err := scanWorkflows()
		if err != nil {
			fmt.Printf("Error scanning workflows: %v\n", err)
			os.Exit(1)
		}

		if len(actions) == 0 {
			fmt.Println("No GitHub Actions found in workflow files")
			return
		}

		checkForUpdates(gc, actions)

		printSummary(actions)

	case "update":
		gc := NewGitHubClient()

		var targetWorkflow string
		if len(os.Args) > 2 {
			targetWorkflow = os.Args[2]
			if !strings.HasPrefix(targetWorkflow, ".github/workflows/") {
				targetWorkflow = ".github/workflows/" + targetWorkflow
			}
		}

		fmt.Println("🔍 Scanning workflow files...")
		actions, err := scanWorkflows()
		if err != nil {
			fmt.Printf("Error scanning workflows: %v\n", err)
			os.Exit(1)
		}

		if len(actions) == 0 {
			fmt.Println("No GitHub Actions found in workflow files")
			return
		}

		checkForUpdates(gc, actions)

		if err := updateActions(actions, targetWorkflow); err != nil {
			fmt.Printf("Error updating actions: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n✅ Update process completed!")

	case "verify":
		if err := verifyPinnedSHAs(); err != nil {
			fmt.Printf("Verification failed: %v\n", err)
			os.Exit(1)
		}

	case "install-hooks":
		if err := installPreCommitHooks(); err != nil {
			fmt.Printf("Failed to install hooks: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
