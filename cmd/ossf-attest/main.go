// Package main implements the OSSF Security Baseline attestation tool.
// This tool runs comprehensive security checks and generates attestation artifacts
// for compliance with OSSF (Open Source Security Foundation) requirements.
//
// The tool follows idempotence, predictability, and atomicity principles by:
// - Reading all configuration from versions.yml for centralized version management
// - Providing atomic updates through structured configuration
// - Ensuring predictable behavior across environments
// - Supporting concurrent execution with proper synchronization
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// VersionsConfig represents the structure of versions.yml
type VersionsConfig struct {
	Tools       map[string]ToolConfig `yaml:"tools"`
	App         AppConfig             `yaml:"app"`
	Config      GlobalConfig          `yaml:"config"`
	Features    FeatureFlags          `yaml:"features"`
	Environment EnvMappings           `yaml:"environment"`
}

// ToolConfig represents configuration for a single tool
type ToolConfig struct {
	Version         string   `yaml:"version"`
	InstallMethod   string   `yaml:"install_method"`
	InstallURL      string   `yaml:"install_url,omitempty"`
	InstallPackage  string   `yaml:"install_package,omitempty"`
	GitHubRepo      string   `yaml:"github_repo,omitempty"`
	Description     string   `yaml:"description"`
	CheckCommand    []string `yaml:"check_command"`
	RunCommand      []string `yaml:"run_command"`
	OutputFile      string   `yaml:"output_file"`
	SpecialHandling string   `yaml:"special_handling,omitempty"`
	FallbackConfig  string   `yaml:"fallback_config,omitempty"`
	AuthEnv         string   `yaml:"auth_env,omitempty"`
}

// AppConfig represents application metadata
type AppConfig struct {
	Name               string `yaml:"name"`
	Version            string `yaml:"version"`
	Description        string `yaml:"description"`
	DefaultProjectName string `yaml:"default_project_name"`
}

// GlobalConfig represents global configuration options
type GlobalConfig struct {
	DefaultOutputDir   string                 `yaml:"default_output_dir"`
	LocalToolsDir      string                 `yaml:"local_tools_dir"`
	ParallelByDefault  bool                   `yaml:"parallel_by_default"`
	InstallPreferences InstallPreferencesConf `yaml:"install_preferences"`
	Execution          ExecutionConf          `yaml:"execution"`
}

// InstallPreferencesConf represents installation preferences
type InstallPreferencesConf struct {
	PreferLocal       bool `yaml:"prefer_local"`
	InstallTimeout    int  `yaml:"install_timeout"`
	InstallRetryCount int  `yaml:"install_retry_count"`
}

// ExecutionConf represents execution configuration
type ExecutionConf struct {
	MaxToolTimeout    int `yaml:"max_tool_timeout"`
	ChannelBufferSize int `yaml:"channel_buffer_size"`
	MaxWorkers        int `yaml:"max_workers"`
}

// FeatureFlags represents feature toggles
type FeatureFlags struct {
	EnableParallelExecution       bool `yaml:"enable_parallel_execution"`
	EnableLocalToolInstallation   bool `yaml:"enable_local_tool_installation"`
	EnablePipFallbackMethods      bool `yaml:"enable_pip_fallback_methods"`
	EnableVersionChecking         bool `yaml:"enable_version_checking"`
	EnableAutomaticPathManagement bool `yaml:"enable_automatic_path_management"`
	EnableCoverageGeneration      bool `yaml:"enable_coverage_generation"`
	EnableSBOMGeneration          bool `yaml:"enable_sbom_generation"`
	EnableSemgrepAuth             bool `yaml:"enable_semgrep_auth"`
}

// EnvMappings represents environment variable mappings
type EnvMappings struct {
	ProjectName    string `yaml:"project_name"`
	ProjectVersion string `yaml:"project_version"`
	SemgrepToken   string `yaml:"semgrep_token"`
	OutputDir      string `yaml:"output_dir"`
	Verbose        string `yaml:"verbose"`
	Parallel       string `yaml:"parallel"`
}

// Config holds the runtime configuration for the attestation process
type Config struct {
	ProjectName    string
	ProjectVersion string
	OutputDir      string
	LocalToolsDir  string
	Parallel       bool
	Verbose        bool
	VersionsConfig *VersionsConfig
}

// Tool represents a security tool and its configuration
type Tool struct {
	Name        string
	Version     string
	CheckCmd    []string
	InstallFunc func(c *Config) error
	RunFunc     func(c *Config) (*Result, error)
	OutputFile  string
	Background  bool
}

// Result represents the output of running a security tool
type Result struct {
	Tool       string        `json:"tool"`
	Version    string        `json:"version"`
	Success    bool          `json:"success"`
	Output     string        `json:"output"`
	Error      string        `json:"error,omitempty"`
	Duration   time.Duration `json:"duration"`
	OutputFile string        `json:"output_file"`
}

// AttestationReport represents the complete attestation results
type AttestationReport struct {
	ProjectName    string             `json:"project_name"`
	ProjectVersion string             `json:"project_version"`
	Timestamp      time.Time          `json:"timestamp"`
	GoVersion      string             `json:"go_version"`
	OS             string             `json:"os"`
	Arch           string             `json:"arch"`
	Results        map[string]*Result `json:"results"`
	Summary        Summary            `json:"summary"`
}

// Summary provides high-level statistics about the attestation
type Summary struct {
	TotalTools    int           `json:"total_tools"`
	SuccessCount  int           `json:"success_count"`
	FailureCount  int           `json:"failure_count"`
	WarningCount  int           `json:"warning_count"`
	TotalDuration time.Duration `json:"total_duration"`
}

// GitHubRelease represents a GitHub release response
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

// GitHubAPI provides methods to interact with GitHub API
type GitHubAPI struct {
	client *http.Client
}

// NewGitHubAPI creates a new GitHub API client
func NewGitHubAPI() *GitHubAPI {
	return &GitHubAPI{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetLatestRelease fetches the latest release for a GitHub repository
func (g *GitHubAPI) GetLatestRelease(ctx context.Context, repo string) (*GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			// Log the error but don't override the main error
			fmt.Printf("Warning: failed to close response body: %v\n", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &release, nil
}

// loadVersionsConfig loads the versions.yml configuration file
func loadVersionsConfig() (*VersionsConfig, error) {
	configPath := filepath.Join("cmd", "ossf-attest", "versions.yml")

	// Try current directory first, then relative to executable
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Try relative to executable
		if execPath, err := os.Executable(); err == nil {
			configPath = filepath.Join(filepath.Dir(execPath), "versions.yml")
		}
	}

	// Try current directory as fallback
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = "versions.yml"
	}

	// Validate config path to prevent path traversal
	cleanConfigPath := filepath.Clean(configPath)
	if strings.Contains(cleanConfigPath, "..") {
		return nil, fmt.Errorf("invalid config path: path traversal detected in %s", configPath)
	}

	data, err := os.ReadFile(cleanConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read versions config from %s: %w", cleanConfigPath, err)
	}

	var config VersionsConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse versions config: %w", err)
	}

	return &config, nil
}

// NewConfig creates a new configuration with sensible defaults
func NewConfig() (*Config, error) {
	versionsConfig, err := loadVersionsConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load versions config: %w", err)
	}

	projectName := versionsConfig.App.DefaultProjectName
	if name := os.Getenv(versionsConfig.Environment.ProjectName); name != "" {
		projectName = name
	}

	projectVersion := getProjectVersion()
	if version := os.Getenv(versionsConfig.Environment.ProjectVersion); version != "" {
		projectVersion = version
	}

	outputDir := versionsConfig.Config.DefaultOutputDir
	if dir := os.Getenv(versionsConfig.Environment.OutputDir); dir != "" {
		sanitizedDir, err := sanitizeOutputPath(dir)
		if err != nil {
			return nil, fmt.Errorf("invalid output directory from environment: %w", err)
		}
		outputDir = sanitizedDir
	} else {
		// Sanitize the default output directory as well
		sanitizedDir, err := sanitizeOutputPath(outputDir)
		if err != nil {
			return nil, fmt.Errorf("invalid default output directory: %w", err)
		}
		outputDir = sanitizedDir
	}

	parallel := versionsConfig.Config.ParallelByDefault
	if p := os.Getenv(versionsConfig.Environment.Parallel); p != "" {
		parallel = strings.ToLower(p) == "true"
	}

	verbose := false
	if v := os.Getenv(versionsConfig.Environment.Verbose); v != "" {
		verbose = strings.ToLower(v) == "true"
	}

	return &Config{
		ProjectName:    projectName,
		ProjectVersion: projectVersion,
		OutputDir:      outputDir,
		LocalToolsDir:  filepath.Join(".", versionsConfig.Config.LocalToolsDir),
		Parallel:       parallel,
		Verbose:        verbose,
		VersionsConfig: versionsConfig,
	}, nil
}

// getProjectVersion attempts to get the project version from git
func getProjectVersion() string {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return "dev"
	}
	return strings.TrimSpace(string(output))
}

// ensureLocalToolsDir creates the local tools directory and adds it to PATH
func (c *Config) ensureLocalToolsDir() error {
	if err := os.MkdirAll(c.LocalToolsDir, 0750); err != nil {
		return fmt.Errorf("failed to create local tools directory: %w", err)
	}

	// Add local tools directory to PATH
	currentPath := os.Getenv("PATH")
	newPath := c.LocalToolsDir + string(os.PathListSeparator) + currentPath

	// Also ensure Go bin directory is in PATH
	goBinOutput, err := exec.Command("go", "env", "GOPATH").Output()
	if err == nil {
		goPath := strings.TrimSpace(string(goBinOutput))
		goBinDir := filepath.Join(goPath, "bin")
		if !strings.Contains(newPath, goBinDir) {
			newPath = goBinDir + string(os.PathListSeparator) + newPath
		}
	}

	return os.Setenv("PATH", newPath)
}

// isCommandAvailable checks if a command is available in PATH
func isCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func findCommandInSystemPath(cmd string) string {
	// Common system paths where tools are typically installed
	commonPaths := []string{
		"/usr/local/bin",
		"/usr/bin",
		"/bin",
		"/opt/homebrew/bin",       // Homebrew on Apple Silicon
		"/usr/local/homebrew/bin", // Homebrew on Intel Mac
		"/opt/local/bin",          // MacPorts
	}

	// Also check Go bin directory
	if goBin, err := exec.Command("go", "env", "GOPATH").Output(); err == nil {
		goPath := strings.TrimSpace(string(goBin))
		goBinDir := filepath.Join(goPath, "bin")
		commonPaths = append([]string{goBinDir}, commonPaths...)
	}

	for _, path := range commonPaths {
		fullPath := filepath.Join(path, cmd)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath
		}
	}

	return ""
}

// getLocalToolVersion attempts to get the version of a locally installed tool
func getLocalToolVersion(toolName string, checkCmd []string) (string, error) {
	// First try to find the tool in PATH
	toolPath, err := exec.LookPath(toolName)
	if err != nil {
		// If not found in current PATH, try system paths
		toolPath = findCommandInSystemPath(toolName)
		if toolPath == "" {
			return "", fmt.Errorf("tool %s not found in PATH or system paths", toolName)
		}
	}

	if len(checkCmd) == 0 {
		return "", fmt.Errorf("no check command configured for %s", toolName)
	}

	// Use absolute path to avoid "relative to current directory" issues
	cmd := make([]string, len(checkCmd))
	copy(cmd, checkCmd)
	if filepath.IsAbs(toolPath) {
		// If we have an absolute path, use it for the first command
		cmd[0] = toolPath
	}

	output, err := runCommand(cmd[0], cmd[1:]...)
	if err != nil {
		return "", fmt.Errorf("failed to get version for %s: %w", toolName, err)
	}

	return strings.TrimSpace(output), nil
}

// checkLatestVersion queries GitHub API for the latest version of a tool
func checkLatestVersion(ctx context.Context, githubAPI *GitHubAPI, repo string) (string, error) {
	if repo == "" {
		return "", fmt.Errorf("no GitHub repository configured")
	}

	release, err := githubAPI.GetLatestRelease(ctx, repo)
	if err != nil {
		return "", fmt.Errorf("failed to get latest release: %w", err)
	}

	return release.TagName, nil
}

// atomicToolInstall performs atomic tool installation with retry logic
func atomicToolInstall(toolName string, installFunc func() error, retryCount int) error {
	var lastErr error

	for attempt := 0; attempt <= retryCount; attempt++ {
		if attempt > 0 {
			fmt.Printf("Retrying installation of %s (attempt %d/%d)...\n", toolName, attempt+1, retryCount+1)
			time.Sleep(time.Duration(attempt) * time.Second) // Exponential backoff
		}

		if err := installFunc(); err != nil {
			lastErr = err
			continue
		}

		// Verify installation succeeded by checking both current PATH and Go bin
		if isCommandAvailable(toolName) {
			return nil
		}

		// For Go tools, also check specifically in GOPATH/bin
		if goBinOutput, err := exec.Command("go", "env", "GOPATH").Output(); err == nil {
			goPath := strings.TrimSpace(string(goBinOutput))
			goBinDir := filepath.Join(goPath, "bin")
			toolPath := filepath.Join(goBinDir, toolName)
			if _, err := os.Stat(toolPath); err == nil {
				// Tool exists in Go bin, update PATH and return success
				currentPath := os.Getenv("PATH")
				if !strings.Contains(currentPath, goBinDir) {
					newPath := goBinDir + string(os.PathListSeparator) + currentPath
					if err := os.Setenv("PATH", newPath); err != nil {
						return fmt.Errorf("failed to update PATH: %w", err)
					}
				}
				return nil
			}
		}

		lastErr = fmt.Errorf("installation appeared to succeed but tool not found in PATH or Go bin")
	}

	return fmt.Errorf("failed to install %s after %d attempts: %w", toolName, retryCount+1, lastErr)
}

// runCommand executes a command and returns its output
// It uses a whitelist-based approach to prevent code injection
func runCommand(name string, args ...string) (string, error) {
	// Use a secure command execution function
	return runSecureCommand(name, args...)
}

// runSecureCommand executes a command using a whitelist-based security model
func runSecureCommand(name string, args ...string) (string, error) {
	// Validate command name using whitelist
	if !isWhitelistedCommand(name) {
		return "", fmt.Errorf("command not in whitelist: %s", name)
	}

	// Validate arguments to prevent injection
	for _, arg := range args {
		if err := validateCommandArg(arg); err != nil {
			return "", fmt.Errorf("invalid command argument: %w", err)
		}
	}

	// Build command using safe construction
	cmd := buildSecureCommand(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// isWhitelistedCommand checks if a command is allowed to execute
func isWhitelistedCommand(name string) bool {
	// Whitelist of allowed commands - only tools we need for OSSF attestation
	allowedCommands := map[string]bool{
		"go":            true,
		"golangci-lint": true,
		"gosec":         true,
		"govulncheck":   true,
		"semgrep":       true,
		"syft":          true,
		"git":           true,
		"curl":          true,
		"sh":            true,
		"bash":          true,
		"pip":           true,
		"pip3":          true,
		"pipx":          true,
		"make":          true,
	}

	// Extract base command name (handle full paths)
	baseName := filepath.Base(name)

	return allowedCommands[baseName]
}

// buildSecureCommand safely constructs an exec.Cmd
func buildSecureCommand(name string, args ...string) *exec.Cmd {
	// Use explicit construction to avoid dynamic command building
	switch filepath.Base(name) {
	case "go":
		return exec.Command("go", args...)
	case "golangci-lint":
		return exec.Command("golangci-lint", args...)
	case "gosec":
		return exec.Command("gosec", args...)
	case "govulncheck":
		return exec.Command("govulncheck", args...)
	case "semgrep":
		return exec.Command("semgrep", args...)
	case "syft":
		return exec.Command("syft", args...)
	case "git":
		return exec.Command("git", args...)
	case "curl":
		return exec.Command("curl", args...)
	case "sh":
		return exec.Command("sh", args...)
	case "bash":
		return exec.Command("bash", args...)
	case "pip":
		return exec.Command("pip", args...)
	case "pip3":
		return exec.Command("pip3", args...)
	case "pipx":
		return exec.Command("pipx", args...)
	case "make":
		return exec.Command("make", args...)
	default:
		// Return an error command instead of dynamic execution
		// This should never be reached due to whitelist check
		return exec.Command("false") // Always fails
	}
}

// validateCommandArg validates that a command argument is safe
func validateCommandArg(arg string) error {
	// For security, we'll be more permissive with arguments but still prevent obvious injection
	// Block null bytes and other dangerous patterns
	if strings.Contains(arg, "\x00") {
		return fmt.Errorf("argument contains null byte")
	}

	// Block some obviously dangerous patterns in arguments
	dangerousPatterns := []string{
		"$(", "`", "${", "||", "&&", ";", "|",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(arg, pattern) {
			return fmt.Errorf("argument contains potentially dangerous pattern: %s", pattern)
		}
	}

	return nil
}

// validateInstallationParams validates parameters used for tool installation
func validateInstallationParams(installURL, localToolsDir, version string) error {
	// Validate install URL
	if installURL == "" {
		return fmt.Errorf("install URL cannot be empty")
	}

	// Basic URL validation - must be HTTPS and from trusted domains
	if !strings.HasPrefix(installURL, "https://") {
		return fmt.Errorf("install URL must use HTTPS")
	}

	// Validate local tools directory path
	if localToolsDir == "" {
		return fmt.Errorf("local tools directory cannot be empty")
	}

	// Ensure no path traversal in directory
	if strings.Contains(localToolsDir, "..") {
		return fmt.Errorf("local tools directory contains path traversal")
	}

	// Validate version string
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}

	// Version should not contain shell metacharacters
	if err := validateCommandArg(version); err != nil {
		return fmt.Errorf("invalid version format: %w", err)
	}

	return nil
}

// installGolangciLintSafely installs golangci-lint using a secure method
func installGolangciLintSafely(installURL, localToolsDir, version string) error {
	// Create a temporary script file with validated content
	tempDir, err := os.MkdirTemp("", "golangci-lint-install")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() {
		if removeErr := os.RemoveAll(tempDir); removeErr != nil {
			fmt.Printf("Warning: failed to remove temp directory %s: %v\n", tempDir, removeErr)
		}
	}()

	scriptPath := filepath.Join(tempDir, "install.sh")

	// Download the install script to a temporary file
	curlCmd := buildSecureCommand("curl", "-sSfL", "-o", scriptPath, installURL)
	if err := curlCmd.Run(); err != nil {
		return fmt.Errorf("failed to download install script: %w", err)
	}

	// Make the script executable (as secure as possible while maintaining functionality)
	// Note: Scripts need execute permission - using 0750 (owner rwx, group rx, others none)
	// #nosec G302 - executable permission required for script execution in temp directory
	if err := os.Chmod(scriptPath, 0750); err != nil {
		return fmt.Errorf("failed to make script executable: %w", err)
	}

	// Execute the script with controlled arguments
	shCmd := buildSecureCommand("sh", scriptPath, "-b", localToolsDir, version)
	output, err := shCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("installation script failed: %w\nOutput: %s", err, output)
	}

	return nil
}

// safeWriteFile safely writes data to a file within the configured output directory
// It validates the file path to prevent directory traversal attacks
func safeWriteFile(outputDir, filename string, data []byte) error {
	// Clean the filename to remove any path separators or traversal attempts
	cleanFilename := filepath.Base(filepath.Clean(filename))

	// Ensure the filename doesn't contain any dangerous characters
	if strings.Contains(cleanFilename, "..") || strings.ContainsAny(cleanFilename, "/\\") {
		return fmt.Errorf("invalid filename: %s", filename)
	}

	// Construct the full path
	fullPath := filepath.Join(outputDir, cleanFilename)

	// Additional check: ensure the final path is still within the output directory
	rel, err := filepath.Rel(outputDir, fullPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("file path outside output directory: %s", filename)
	}

	return os.WriteFile(fullPath, data, 0600)
}

// sanitizeOutputPath validates and sanitizes an output directory path
// to prevent path traversal attacks
func sanitizeOutputPath(outputPath string) (string, error) {
	if outputPath == "" {
		return "", fmt.Errorf("output path cannot be empty")
	}

	// Clean the path to remove any .. or . components
	cleaned := filepath.Clean(outputPath)

	// Convert to absolute path to detect traversal attempts
	abs, err := filepath.Abs(cleaned)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	// Ensure the output path is within or relative to the current working directory
	// This prevents writing to arbitrary locations on the filesystem
	relPath, err := filepath.Rel(cwd, abs)
	if err != nil {
		return "", fmt.Errorf("failed to determine relative path: %w", err)
	}

	// Check if the relative path tries to escape the current directory
	if strings.HasPrefix(relPath, "..") {
		return "", fmt.Errorf("output path cannot traverse outside current directory: %s", outputPath)
	}

	// Ensure the directory exists and create it if it doesn't (secure permissions)
	if err := os.MkdirAll(cleaned, 0750); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	return cleaned, nil
}

// compareVersions performs a simple semantic version comparison
// Returns true if local version is compatible with expected version
func compareVersions(localVersion, expectedVersion string) bool {
	// Normalize versions by removing 'v' prefix
	local := strings.TrimPrefix(strings.TrimSpace(localVersion), "v")
	expected := strings.TrimPrefix(strings.TrimSpace(expectedVersion), "v")

	// Simple string contains check for now - this covers most cases
	// where the local version string contains the expected version
	if strings.Contains(local, expected) {
		return true
	}

	// Exact match after normalization
	if local == expected {
		return true
	}

	return false
}

// preferLocalTool determines whether to use a local tool installation
func preferLocalTool(c *Config, toolName string, localVersion string, expectedVersion string) bool {
	// Always prefer local if configured to do so
	if c.VersionsConfig.Config.InstallPreferences.PreferLocal {
		if c.Verbose {
			fmt.Printf("Debug: preferLocal=true for %s, using local version\n", toolName)
		}
		return true
	}

	// If version checking is disabled, use local
	if !c.VersionsConfig.Features.EnableVersionChecking {
		return true
	}

	// Check version compatibility
	return compareVersions(localVersion, expectedVersion)
}

// Install functions for each tool
func installGolangciLint(c *Config) error {
	toolConfig := c.VersionsConfig.Tools["golangci-lint"]
	toolName := "golangci-lint"

	// Check if tool is available locally
	if localVersion, err := getLocalToolVersion(toolName, toolConfig.CheckCommand); err == nil {
		if preferLocalTool(c, toolName, localVersion, toolConfig.Version) {
			if c.Verbose {
				fmt.Printf("Using local %s (version: %s)\n", toolName, localVersion)
			}
			return nil
		} else if c.Verbose {
			fmt.Printf("Local %s version %s differs from configured %s, installing configured version\n", toolName, localVersion, toolConfig.Version)
		}
	}

	// Check for latest version from GitHub if configured
	if toolConfig.GitHubRepo != "" && c.Verbose {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		githubAPI := NewGitHubAPI()
		if latestVersion, err := checkLatestVersion(ctx, githubAPI, toolConfig.GitHubRepo); err == nil {
			if latestVersion != toolConfig.Version {
				fmt.Printf("Note: %s latest version is %s (configured: %s)\n", toolName, latestVersion, toolConfig.Version)
			}
		}
	}

	// Perform atomic installation with retry
	installFunc := func() error {
		fmt.Printf("Installing %s %s to %s...\n", toolName, toolConfig.Version, c.LocalToolsDir)

		// Validate installation parameters to prevent injection
		if err := validateInstallationParams(toolConfig.InstallURL, c.LocalToolsDir, toolConfig.Version); err != nil {
			return fmt.Errorf("invalid installation parameters: %w", err)
		}

		// Use a safer installation approach - download and execute separately
		return installGolangciLintSafely(toolConfig.InstallURL, c.LocalToolsDir, toolConfig.Version)
	}

	return atomicToolInstall(toolName, installFunc, c.VersionsConfig.Config.InstallPreferences.InstallRetryCount)
}

func installGosec(c *Config) error {
	toolConfig := c.VersionsConfig.Tools["gosec"]
	toolName := "gosec"

	// Check if tool is available locally
	if localVersion, err := getLocalToolVersion(toolName, toolConfig.CheckCommand); err == nil {
		if preferLocalTool(c, toolName, localVersion, toolConfig.Version) {
			if c.Verbose {
				fmt.Printf("Using local %s (version: %s)\n", toolName, localVersion)
			}
			return nil
		} else if c.Verbose {
			fmt.Printf("Local %s version %s differs from configured %s, installing configured version\n", toolName, localVersion, toolConfig.Version)
		}
	}

	// Check for latest version from GitHub if configured
	if toolConfig.GitHubRepo != "" && c.Verbose {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		githubAPI := NewGitHubAPI()
		if latestVersion, err := checkLatestVersion(ctx, githubAPI, toolConfig.GitHubRepo); err == nil {
			if latestVersion != toolConfig.Version {
				fmt.Printf("Note: %s latest version is %s (configured: %s)\n", toolName, latestVersion, toolConfig.Version)
			}
		}
	}

	// Perform atomic installation with retry
	installFunc := func() error {
		fmt.Printf("Installing %s %s...\n", toolName, toolConfig.Version)

		// Validate inputs to prevent command injection
		if err := validateGoInstallPackage(toolConfig.InstallPackage); err != nil {
			return fmt.Errorf("invalid install package: %w", err)
		}
		if err := validateVersionString(toolConfig.Version); err != nil {
			return fmt.Errorf("invalid version: %w", err)
		}

		// Use separate arguments to prevent injection
		packageAtVersion := toolConfig.InstallPackage + "@" + toolConfig.Version
		cmd := buildSecureCommand("go", "install", packageAtVersion)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("installation command failed: %w\nOutput: %s", err, output)
		}
		return nil
	}

	return atomicToolInstall(toolName, installFunc, c.VersionsConfig.Config.InstallPreferences.InstallRetryCount)
}

func installGovulncheck(c *Config) error {
	toolConfig := c.VersionsConfig.Tools["govulncheck"]
	toolName := "govulncheck"

	// Check if tool is available locally
	if _, err := getLocalToolVersion(toolName, toolConfig.CheckCommand); err == nil {
		// For govulncheck, also check the Go module version if version checking is enabled
		if c.VersionsConfig.Features.EnableVersionChecking && !c.VersionsConfig.Config.InstallPreferences.PreferLocal {
			// Validate package name before using in command
			if err := validateGoInstallPackage(toolConfig.InstallPackage); err != nil {
				return fmt.Errorf("invalid package name for version check: %w", err)
			}

			cmd := buildSecureCommand("go", "list", "-m", "-f", "{{.Version}}", toolConfig.InstallPackage)
			if output, err := cmd.Output(); err == nil {
				moduleVersion := strings.TrimSpace(string(output))
				if compareVersions(moduleVersion, toolConfig.Version) {
					if c.Verbose {
						fmt.Printf("Using local %s (module version: %s)\n", toolName, moduleVersion)
					}
					return nil
				} else if c.Verbose {
					fmt.Printf("Local %s module version %s differs from configured %s, installing configured version\n", toolName, moduleVersion, toolConfig.Version)
				}
			}
		} else {
			if c.Verbose {
				fmt.Printf("Using local %s\n", toolName)
			}
			return nil
		}
	}

	// Check for latest version from GitHub if configured
	if toolConfig.GitHubRepo != "" && c.Verbose {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		githubAPI := NewGitHubAPI()
		if latestVersion, err := checkLatestVersion(ctx, githubAPI, toolConfig.GitHubRepo); err == nil {
			if latestVersion != toolConfig.Version {
				fmt.Printf("Note: %s latest version is %s (configured: %s)\n", toolName, latestVersion, toolConfig.Version)
			}
		}
	}

	// Perform atomic installation with retry
	installFunc := func() error {
		fmt.Printf("Installing %s %s...\n", toolName, toolConfig.Version)

		// Validate inputs to prevent command injection
		if err := validateGoInstallPackage(toolConfig.InstallPackage); err != nil {
			return fmt.Errorf("invalid install package: %w", err)
		}
		if err := validateVersionString(toolConfig.Version); err != nil {
			return fmt.Errorf("invalid version: %w", err)
		}

		// Use separate arguments to prevent injection
		packageAtVersion := toolConfig.InstallPackage + "@" + toolConfig.Version
		cmd := buildSecureCommand("go", "install", packageAtVersion)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("installation command failed: %w\nOutput: %s", err, output)
		}
		return nil
	}

	return atomicToolInstall(toolName, installFunc, c.VersionsConfig.Config.InstallPreferences.InstallRetryCount)
}

func installSyft(c *Config) error {
	toolConfig := c.VersionsConfig.Tools["syft"]
	toolName := "syft"

	// Check if tool is available locally
	if localVersion, err := getLocalToolVersion(toolName, toolConfig.CheckCommand); err == nil {
		if preferLocalTool(c, toolName, localVersion, toolConfig.Version) {
			if c.Verbose {
				fmt.Printf("Using local %s\n", toolName)
			}
			return nil
		} else if c.Verbose {
			fmt.Printf("Local %s version differs from configured %s, installing configured version\n", toolName, toolConfig.Version)
		}
	} else if c.Verbose {
		fmt.Printf("Debug: Failed to get local %s version: %v\n", toolName, err)
	}

	// Check for latest version from GitHub if configured
	if toolConfig.GitHubRepo != "" && c.Verbose {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		githubAPI := NewGitHubAPI()
		if latestVersion, err := checkLatestVersion(ctx, githubAPI, toolConfig.GitHubRepo); err == nil {
			if latestVersion != toolConfig.Version {
				fmt.Printf("Note: %s latest version is %s (configured: %s)\n", toolName, latestVersion, toolConfig.Version)
			}
		}
	}

	// Perform atomic installation with retry
	installFunc := func() error {
		fmt.Printf("Installing %s %s...\n", toolName, toolConfig.Version)

		// Validate inputs to prevent command injection
		if err := validateGoInstallPackage(toolConfig.InstallPackage); err != nil {
			return fmt.Errorf("invalid install package: %w", err)
		}
		if err := validateVersionString(toolConfig.Version); err != nil {
			return fmt.Errorf("invalid version: %w", err)
		}

		// Use separate arguments to prevent injection
		packageAtVersion := toolConfig.InstallPackage + "@" + toolConfig.Version
		cmd := buildSecureCommand("go", "install", packageAtVersion)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("installation command failed: %w\nOutput: %s", err, output)
		}
		return nil
	}

	return atomicToolInstall(toolName, installFunc, c.VersionsConfig.Config.InstallPreferences.InstallRetryCount)
}

func installSemgrep(c *Config) error {
	toolConfig := c.VersionsConfig.Tools["semgrep"]
	toolName := "semgrep"

	// Check if tool is available locally
	if localVersion, err := getLocalToolVersion(toolName, toolConfig.CheckCommand); err == nil {
		if preferLocalTool(c, toolName, localVersion, toolConfig.Version) {
			if c.Verbose {
				fmt.Printf("Using local %s (version: %s)\n", toolName, localVersion)
			}
			return nil
		} else if c.Verbose {
			fmt.Printf("Local %s version %s differs from configured %s, installing configured version\n", toolName, localVersion, toolConfig.Version)
		}
	}

	if !c.VersionsConfig.Features.EnablePipFallbackMethods {
		return fmt.Errorf("pip installation methods are disabled for %s", toolName)
	}

	// Perform atomic installation with retry
	installFunc := func() error {
		fmt.Printf("Installing %s %s...\n", toolName, toolConfig.Version)

		// Validate version before using in commands
		if err := validateVersionString(toolConfig.Version); err != nil {
			return fmt.Errorf("invalid version for semgrep installation: %w", err)
		}

		// Try pipx first, then pip3, then pip
		semgrepPackage := "semgrep==" + toolConfig.Version
		if isCommandAvailable("pipx") {
			cmd := buildSecureCommand("pipx", "install", "--force", semgrepPackage)
			if err := cmd.Run(); err == nil {
				return nil
			}
		}

		if isCommandAvailable("pip3") {
			cmd := buildSecureCommand("pip3", "install", "--user", "--force-reinstall", semgrepPackage)
			if err := cmd.Run(); err == nil {
				// Add ~/.local/bin to PATH if needed
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get user home directory: %w", err)
				}
				localBin := filepath.Join(homeDir, ".local", "bin")
				currentPath := os.Getenv("PATH")
				if !strings.Contains(currentPath, localBin) {
					newPath := localBin + string(os.PathListSeparator) + currentPath
					if err := os.Setenv("PATH", newPath); err != nil {
						return fmt.Errorf("failed to update PATH: %w", err)
					}
				}
				return nil
			}
		}

		if isCommandAvailable("pip") {
			cmd := buildSecureCommand("pip", "install", "--user", "--force-reinstall", semgrepPackage)
			return cmd.Run()
		}

		return fmt.Errorf("pipx, pip3, or pip is required to install semgrep")
	}

	return atomicToolInstall(toolName, installFunc, c.VersionsConfig.Config.InstallPreferences.InstallRetryCount)
}

// Run functions for each tool
func runGolangciLint(c *Config) (*Result, error) {
	toolConfig := c.VersionsConfig.Tools["golangci-lint"]
	start := time.Now()
	result := &Result{
		Tool:       "golangci-lint",
		OutputFile: toolConfig.OutputFile,
	}

	// Get version
	if version, err := runCommand("golangci-lint", "--version"); err == nil {
		result.Version = strings.TrimSpace(version)
	}

	// Run linting using configured command
	output, err := runCommand(toolConfig.RunCommand[0], toolConfig.RunCommand[1:]...)
	result.Output = output
	result.Duration = time.Since(start)
	result.Success = err == nil

	if err != nil {
		result.Error = err.Error()
	}

	// Write output to file
	if err := safeWriteFile(c.OutputDir, result.OutputFile, []byte(output)); err != nil {
		return result, fmt.Errorf("failed to write output file: %w", err)
	}

	return result, nil
}

func runGosec(c *Config) (*Result, error) {
	toolConfig := c.VersionsConfig.Tools["gosec"]
	start := time.Now()
	result := &Result{
		Tool:       "gosec",
		OutputFile: toolConfig.OutputFile,
	}

	// Get version
	if version, err := runCommand("gosec", "--version"); err == nil {
		result.Version = strings.TrimSpace(version)
	}

	// Run gosec using configured command
	output, err := runCommand(toolConfig.RunCommand[0], toolConfig.RunCommand[1:]...)
	result.Output = output
	result.Duration = time.Since(start)
	result.Success = err == nil

	if err != nil {
		result.Error = err.Error()
	}

	// Write output to file
	if err := safeWriteFile(c.OutputDir, result.OutputFile, []byte(output)); err != nil {
		return result, fmt.Errorf("failed to write output file: %w", err)
	}

	return result, nil
}

func runGovulncheck(c *Config) (*Result, error) {
	toolConfig := c.VersionsConfig.Tools["govulncheck"]
	start := time.Now()
	result := &Result{
		Tool:       "govulncheck",
		OutputFile: toolConfig.OutputFile,
	}

	// Get version
	if version, err := runCommand("govulncheck", "version"); err == nil {
		result.Version = strings.TrimSpace(version)
	}

	// Run govulncheck using configured command
	output, err := runCommand(toolConfig.RunCommand[0], toolConfig.RunCommand[1:]...)
	result.Output = output
	result.Duration = time.Since(start)
	result.Success = err == nil

	if err != nil {
		result.Error = err.Error()
	}

	// Write output to file
	if err := safeWriteFile(c.OutputDir, result.OutputFile, []byte(output)); err != nil {
		return result, fmt.Errorf("failed to write output file: %w", err)
	}

	return result, nil
}

func runSemgrep(c *Config) (*Result, error) {
	toolConfig := c.VersionsConfig.Tools["semgrep"]
	start := time.Now()
	result := &Result{
		Tool:       "semgrep",
		OutputFile: toolConfig.OutputFile,
	}

	// Get version
	if version, err := runCommand("semgrep", "--version"); err == nil {
		result.Version = strings.TrimSpace(version)
	}

	// Check for existing Semgrep authentication or login if token is available and auth is enabled
	if c.VersionsConfig.Features.EnableSemgrepAuth {
		// First check if user is already authenticated by checking settings file
		homeDir, err := os.UserHomeDir()
		if err != nil {
			result.Error = fmt.Sprintf("Failed to get user home directory: %v", err)
			return result, nil
		}
		semgrepSettingsPath := filepath.Join(homeDir, ".semgrep", "settings.yml")

		isAuthenticated := false
		if _, err := os.Stat(semgrepSettingsPath); err == nil {
			// Settings file exists, check if it contains an API token
			// Validate the settings path to prevent path traversal
			cleanSettingsPath := filepath.Clean(semgrepSettingsPath)
			if strings.Contains(cleanSettingsPath, "..") {
				result.Error = "Invalid semgrep settings path"
				return result, nil
			}

			if settingsData, err := os.ReadFile(cleanSettingsPath); err == nil {
				if strings.Contains(string(settingsData), "api_token") {
					isAuthenticated = true
					if c.Verbose {
						fmt.Println("Semgrep: Using existing authentication from ~/.semgrep/settings.yml")
					}
				}
			}
		}

		// Only attempt login if not already authenticated
		if !isAuthenticated {
			if token := os.Getenv(toolConfig.AuthEnv); token != "" {
				if c.Verbose {
					fmt.Println("Semgrep: Logging in with environment token")
				}
				if _, err := runCommand("semgrep", "login", "--token", token); err != nil {
					if c.Verbose {
						fmt.Printf("Warning: Semgrep login failed: %v\n", err)
					}
				}
			}
		}
	}

	// Build command with configured base command
	cmd := make([]string, len(toolConfig.RunCommand))
	copy(cmd, toolConfig.RunCommand)

	// Try configured config first, fallback if specified
	if toolConfig.FallbackConfig != "" {
		// Test if the primary config works
		testCmd := []string{"semgrep", "scan", "--config", cmd[3], "--dryrun", "."}
		if _, err := runCommand(testCmd[0], testCmd[1:]...); err != nil {
			fmt.Printf("Semgrep config %s not found, falling back to %s\n", cmd[3], toolConfig.FallbackConfig)
			cmd[3] = toolConfig.FallbackConfig
		}
	}

	// Run semgrep
	output, err := runCommand(cmd[0], cmd[1:]...)
	result.Output = output
	result.Duration = time.Since(start)
	result.Success = err == nil

	if err != nil {
		result.Error = err.Error()
	}

	// Write output to file
	if err := safeWriteFile(c.OutputDir, result.OutputFile, []byte(output)); err != nil {
		return result, fmt.Errorf("failed to write output file: %w", err)
	}

	return result, nil
}

func runTests(c *Config) (*Result, error) {
	toolConfig := c.VersionsConfig.Tools["go-test"]
	start := time.Now()
	result := &Result{
		Tool:       "go-test",
		OutputFile: toolConfig.OutputFile,
	}

	// Get Go version
	if version, err := runCommand("go", "version"); err == nil {
		result.Version = strings.TrimSpace(version)
	}

	// Build command with coverage if enabled
	cmd := make([]string, len(toolConfig.RunCommand))
	copy(cmd, toolConfig.RunCommand)

	if c.VersionsConfig.Features.EnableCoverageGeneration {
		coverageFile := filepath.Join(c.OutputDir, "coverage.txt")
		// Update coverage file path in command
		for i, arg := range cmd {
			if strings.HasPrefix(arg, "-coverprofile=") {
				cmd[i] = "-coverprofile=" + coverageFile
				break
			}
		}
	}

	// Run tests
	output, err := runCommand(cmd[0], cmd[1:]...)
	result.Output = output
	result.Duration = time.Since(start)
	result.Success = err == nil

	if err != nil {
		result.Error = err.Error()
	}

	// Write output to file
	if err := safeWriteFile(c.OutputDir, result.OutputFile, []byte(output)); err != nil {
		return result, fmt.Errorf("failed to write output file: %w", err)
	}

	return result, nil
}

func runSyft(c *Config) (*Result, error) {
	toolConfig := c.VersionsConfig.Tools["syft"]
	start := time.Now()
	result := &Result{
		Tool:       "syft",
		OutputFile: toolConfig.OutputFile,
	}

	if !c.VersionsConfig.Features.EnableSBOMGeneration {
		result.Success = true
		result.Output = "SBOM generation disabled by feature flag"
		result.Duration = time.Since(start)
		return result, nil
	}

	// Get version
	if version, err := runCommand("syft", "version"); err == nil {
		result.Version = strings.TrimSpace(version)
	}

	// Build command with project metadata
	cmd := make([]string, len(toolConfig.RunCommand))
	copy(cmd, toolConfig.RunCommand)

	// Add source name and version
	cmd = append(cmd, "--source-name", c.ProjectName, "--source-version", c.ProjectVersion)

	// Generate SBOM
	output, err := runCommand(cmd[0], cmd[1:]...)
	result.Output = output
	result.Duration = time.Since(start)
	result.Success = err == nil

	if err != nil {
		result.Error = err.Error()
	} else {
		// Write SBOM to file
		if err := safeWriteFile(c.OutputDir, result.OutputFile, []byte(output)); err != nil {
			return result, fmt.Errorf("failed to write SBOM file: %w", err)
		}
	}

	return result, nil
}

// setupTools installs all required tools based on configuration with concurrent execution
func setupTools(c *Config) error {
	fmt.Println("=== Setting up tools ===")

	// Get tools that need installation from configuration
	installFunctions := map[string]func(*Config) error{
		"golangci-lint": installGolangciLint,
		"gosec":         installGosec,
		"govulncheck":   installGovulncheck,
		"syft":          installSyft,
		"semgrep":       installSemgrep,
	}

	if c.Parallel && c.VersionsConfig.Features.EnableParallelExecution {
		// Install tools in parallel with controlled concurrency
		var wg sync.WaitGroup
		var mu sync.Mutex
		semaphore := make(chan struct{}, c.VersionsConfig.Config.Execution.MaxWorkers)
		errors := make([]error, 0)

		for toolName := range c.VersionsConfig.Tools {
			if toolName == "go-test" {
				continue // Built-in tool, no installation needed
			}

			installFunc, exists := installFunctions[toolName]
			if !exists {
				fmt.Printf("Warning: No installation function for tool %s\n", toolName)
				continue
			}

			wg.Add(1)
			go func(name string, installFunc func(*Config) error) {
				defer wg.Done()

				// Acquire semaphore
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				fmt.Printf("Checking %s...\n", name)
				if err := installFunc(c); err != nil {
					mu.Lock()
					errors = append(errors, fmt.Errorf("failed to setup %s: %w", name, err))
					mu.Unlock()
				}
			}(toolName, installFunc)
		}

		wg.Wait()

		if len(errors) > 0 {
			var errorStrs []string
			for _, err := range errors {
				errorStrs = append(errorStrs, err.Error())
			}
			return fmt.Errorf("tool setup failed: %s", strings.Join(errorStrs, "; "))
		}
	} else {
		// Install tools sequentially
		for toolName := range c.VersionsConfig.Tools {
			if toolName == "go-test" {
				continue // Built-in tool, no installation needed
			}

			installFunc, exists := installFunctions[toolName]
			if !exists {
				fmt.Printf("Warning: No installation function for tool %s\n", toolName)
				continue
			}

			fmt.Printf("Checking %s...\n", toolName)
			if err := installFunc(c); err != nil {
				return fmt.Errorf("failed to setup %s: %w", toolName, err)
			}
		}
	}

	return nil
}

// runAllTools executes all security tools based on configuration
func runAllTools(c *Config) *AttestationReport {
	// Build tool execution functions dynamically from configuration
	runFunctions := map[string]func(*Config) (*Result, error){
		"golangci-lint": runGolangciLint,
		"gosec":         runGosec,
		"govulncheck":   runGovulncheck,
		"semgrep":       runSemgrep,
		"go-test":       runTests,
		"syft":          runSyft,
	}

	report := &AttestationReport{
		ProjectName:    c.ProjectName,
		ProjectVersion: c.ProjectVersion,
		Timestamp:      time.Now(),
		GoVersion:      runtime.Version(),
		OS:             runtime.GOOS,
		Arch:           runtime.GOARCH,
		Results:        make(map[string]*Result),
	}

	fmt.Println("=== Running OSSF Security Baseline Checks ===")

	if c.Parallel && c.VersionsConfig.Features.EnableParallelExecution {
		// Run tools in parallel with controlled concurrency
		var wg sync.WaitGroup
		var mu sync.Mutex
		semaphore := make(chan struct{}, c.VersionsConfig.Config.Execution.MaxWorkers)

		for toolName := range c.VersionsConfig.Tools {
			runFunc, exists := runFunctions[toolName]
			if !exists {
				fmt.Printf("Warning: No run function for tool %s\n", toolName)
				continue
			}

			wg.Add(1)
			go func(name string, runFunc func(*Config) (*Result, error)) {
				defer wg.Done()

				// Acquire semaphore
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				fmt.Printf("Running %s...\n", name)
				result, err := runFunc(c)
				if err != nil {
					fmt.Printf("Error running %s: %v\n", name, err)
				}

				mu.Lock()
				report.Results[name] = result
				mu.Unlock()
			}(toolName, runFunc)
		}

		wg.Wait()
	} else {
		// Run tools sequentially
		for toolName := range c.VersionsConfig.Tools {
			runFunc, exists := runFunctions[toolName]
			if !exists {
				fmt.Printf("Warning: No run function for tool %s\n", toolName)
				continue
			}

			fmt.Printf("Running %s...\n", toolName)
			result, err := runFunc(c)
			if err != nil {
				fmt.Printf("Error running %s: %v\n", toolName, err)
			}
			report.Results[toolName] = result
		}
	}

	// Calculate summary
	var totalDuration time.Duration
	for _, result := range report.Results {
		report.Summary.TotalTools++
		totalDuration += result.Duration

		if result.Success {
			report.Summary.SuccessCount++
		} else {
			report.Summary.FailureCount++
		}
	}
	report.Summary.TotalDuration = totalDuration

	return report
}

// generateReport generates and saves the attestation report
func generateReport(report *AttestationReport, outputPath string) error {
	// Validate output path to prevent path traversal
	safeOutputPath, err := sanitizeOutputPath(outputPath)
	if err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	}

	// Generate JSON report
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	if err := safeWriteFile(safeOutputPath, "ossf-attestation-report.json", jsonData); err != nil {
		return fmt.Errorf("failed to write JSON report: %w", err)
	}

	jsonFile := filepath.Join(safeOutputPath, "ossf-attestation-report.json")
	fmt.Printf("✅ JSON report saved to: %s\n", jsonFile)

	// Generate human-readable summary
	summary := generateTextSummary(report)
	if err := safeWriteFile(safeOutputPath, "ossf-attestation-summary.txt", []byte(summary)); err != nil {
		return fmt.Errorf("failed to write text summary: %w", err)
	}

	textFile := filepath.Join(safeOutputPath, "ossf-attestation-summary.txt")
	fmt.Printf("✅ Text summary saved to: %s\n", textFile)
	return nil
}

func generateTextSummary(report *AttestationReport) string {
	var sb strings.Builder

	sb.WriteString("OSSF Security Baseline Attestation Report\n")
	sb.WriteString("==========================================\n\n")

	sb.WriteString(fmt.Sprintf("Project: %s\n", report.ProjectName))
	sb.WriteString(fmt.Sprintf("Version: %s\n", report.ProjectVersion))
	sb.WriteString(fmt.Sprintf("Timestamp: %s\n", report.Timestamp.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Go Version: %s\n", report.GoVersion))
	sb.WriteString(fmt.Sprintf("Platform: %s/%s\n\n", report.OS, report.Arch))

	sb.WriteString("Summary:\n")
	sb.WriteString(fmt.Sprintf("  Total Tools: %d\n", report.Summary.TotalTools))
	sb.WriteString(fmt.Sprintf("  Successful: %d\n", report.Summary.SuccessCount))
	sb.WriteString(fmt.Sprintf("  Failed: %d\n", report.Summary.FailureCount))
	sb.WriteString(fmt.Sprintf("  Total Duration: %v\n\n", report.Summary.TotalDuration))

	sb.WriteString("Tool Results:\n")
	for name, result := range report.Results {
		status := "✅ PASS"
		if !result.Success {
			status = "❌ FAIL"
		}

		sb.WriteString(fmt.Sprintf("  %s: %s (Duration: %v)\n", name, status, result.Duration))
		if result.Version != "" {
			sb.WriteString(fmt.Sprintf("    Version: %s\n", result.Version))
		}
		if result.OutputFile != "" {
			sb.WriteString(fmt.Sprintf("    Output: %s\n", result.OutputFile))
		}
		if !result.Success && result.Error != "" {
			sb.WriteString(fmt.Sprintf("    Error: %s\n", result.Error))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// validateGoInstallPackage validates that a Go package path is safe
func validateGoInstallPackage(pkg string) error {
	if pkg == "" {
		return fmt.Errorf("package name cannot be empty")
	}

	// Allow only alphanumeric characters, dots, slashes, hyphens for Go packages
	for _, char := range pkg {
		if (char < 'a' || char > 'z') &&
			(char < 'A' || char > 'Z') &&
			(char < '0' || char > '9') &&
			char != '.' && char != '/' && char != '-' && char != '_' {
			return fmt.Errorf("package name contains invalid character: %c", char)
		}
	}

	// Must start with a domain-like pattern (basic validation)
	if !strings.Contains(pkg, ".") || strings.HasPrefix(pkg, ".") {
		return fmt.Errorf("package name must be a valid Go module path")
	}

	return nil
}

// validateVersionString validates that a version string is safe
func validateVersionString(version string) error {
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}

	// Allow only alphanumeric characters, dots, hyphens for versions
	for _, char := range version {
		if (char < 'a' || char > 'z') &&
			(char < 'A' || char > 'Z') &&
			(char < '0' || char > '9') &&
			char != '.' && char != '-' && char != 'v' {
			return fmt.Errorf("version contains invalid character: %c", char)
		}
	}

	return nil
}

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// Parse command line arguments
	for i, arg := range os.Args[1:] {
		switch arg {
		case "-v", "--verbose":
			config.Verbose = true
		case "-s", "--sequential":
			config.Parallel = false
		case "-o", "--output":
			if i+1 < len(os.Args[1:]) {
				sanitizedPath, err := sanitizeOutputPath(os.Args[i+2])
				if err != nil {
					log.Fatalf("Invalid output directory: %v", err)
				}
				config.OutputDir = sanitizedPath
			}
		case "-h", "--help":
			fmt.Printf("OSSF Security Baseline Attestation Tool v%s\n", config.VersionsConfig.App.Version)
			fmt.Printf("Description: %s\n\n", config.VersionsConfig.App.Description)
			fmt.Println("Usage: ossf-attest [options]")
			fmt.Println("Options:")
			fmt.Println("  -v, --verbose     Enable verbose output")
			fmt.Println("  -s, --sequential  Run tools sequentially instead of parallel")
			fmt.Println("  -o, --output DIR  Output directory for reports")
			fmt.Println("  -h, --help        Show this help message")
			fmt.Println("\nConfigured Tools:")
			for name, tool := range config.VersionsConfig.Tools {
				fmt.Printf("  %-15s %s (%s)\n", name, tool.Version, tool.Description)
			}
			fmt.Println("\nEnvironment Variables:")
			fmt.Printf("  %s    Project name (default: %s)\n", config.VersionsConfig.Environment.ProjectName, config.ProjectName)
			fmt.Printf("  %s Project version\n", config.VersionsConfig.Environment.ProjectVersion)
			fmt.Printf("  %s    Output directory\n", config.VersionsConfig.Environment.OutputDir)
			fmt.Printf("  %s     Enable verbose output\n", config.VersionsConfig.Environment.Verbose)
			fmt.Printf("  %s    Run in parallel mode\n", config.VersionsConfig.Environment.Parallel)
			if config.VersionsConfig.Features.EnableSemgrepAuth {
				fmt.Printf("  %s Semgrep authentication token\n", config.VersionsConfig.Tools["semgrep"].AuthEnv)
			}
			os.Exit(0)
		}
	}

	if config.Verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		fmt.Printf("Configuration loaded from versions.yml\n")
		fmt.Printf("Project: %s v%s\n", config.ProjectName, config.ProjectVersion)
		fmt.Printf("Output Directory: %s\n", config.OutputDir)
		fmt.Printf("Local Tools Directory: %s\n", config.LocalToolsDir)
		fmt.Printf("Parallel Execution: %t\n", config.Parallel)
		fmt.Printf("Tools configured: %d\n", len(config.VersionsConfig.Tools))
	}

	// Setup local tools directory
	if config.VersionsConfig.Features.EnableLocalToolInstallation {
		if err := config.ensureLocalToolsDir(); err != nil {
			log.Fatalf("Failed to setup local tools directory: %v", err)
		}
	}

	// Install tools
	if err := setupTools(config); err != nil {
		log.Fatalf("Failed to setup tools: %v", err)
	}

	// Run attestation
	report := runAllTools(config)

	// Generate reports
	if err := generateReport(report, config.OutputDir); err != nil {
		log.Fatalf("Failed to generate reports: %v", err)
	}

	// Print summary to console
	fmt.Println("\n" + generateTextSummary(report))

	// Exit with error code if any tools failed
	if report.Summary.FailureCount > 0 {
		fmt.Printf("\n❌ Attestation completed with %d failures\n", report.Summary.FailureCount)
		os.Exit(1)
	}

	fmt.Printf("\n✅ OSSF Security Baseline Attestation completed successfully!")
	fmt.Printf("\n📊 %d tools executed in %v\n", report.Summary.TotalTools, report.Summary.TotalDuration)
}
