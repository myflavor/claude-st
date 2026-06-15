package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [args...]",
	Short: "Start Claude CLI with optional provider switch",
	Long: `Start Claude CLI. Switch provider configuration and start claude.

Short flags (parsed by claude-cli):
  -P <name>    Use provider config from ~/.claude-cli/<name>/settings.json
  -S           Skip permissions (adds -dangerously-skip-permissions)

All other arguments are passed directly to claude.

Examples:
  claude-cli start                           # Start claude directly
  claude-cli start -P minimax                # Switch to minimax and start
  claude-cli start -P minimax -S             # Switch and skip permissions
  claude-cli start -P minimax --model opus   # Pass args to claude
  claude-cli start -S --model opus           # Skip permissions with model`,
	DisableFlagParsing: true,
	Run:                runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func runStart(cmd *cobra.Command, args []string) {
	var provider string
	skipPermissions := false
	claudeArgs := make([]string, 0, len(args))

	// Parse our short flags, pass everything else to claude
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-P":
			if i+1 >= len(args) {
				exitWithError("Flag -P requires a value")
			}
			provider = args[i+1]
			i++
		case arg == "-S":
			skipPermissions = true
		case arg == "-H", arg == "--help":
			cmd.Help()
			return
		default:
			claudeArgs = append(claudeArgs, arg)
		}
	}

	// Prepare arguments for claude command
	finalArgs := make([]string, 0, len(claudeArgs)+2)

	// If provider is specified, generate a merged settings file and pass via --settings
	if provider != "" {
		settingsPath, err := buildProviderSettings(provider)
		if err != nil {
			exitWithError(err.Error())
		}
		finalArgs = append(finalArgs, "--settings", settingsPath)
		fmt.Printf("Using provider: %s (settings: %s)\n", provider, settingsPath)
	}

	// Add skip permissions flag if set
	if skipPermissions {
		finalArgs = append(finalArgs, "-dangerously-skip-permissions")
	}

	// Add remaining arguments
	finalArgs = append(finalArgs, claudeArgs...)

	// Execute claude command
	claudeCmd := exec.Command("claude", finalArgs...)
	claudeCmd.Stdin = os.Stdin
	claudeCmd.Stdout = os.Stdout
	claudeCmd.Stderr = os.Stderr

	fmt.Printf("Starting: claude %s\n", strings.Join(finalArgs, " "))

	if err := claudeCmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		exitWithError(fmt.Sprintf("Failed to execute claude: %v", err))
	}
}

func buildProviderSettings(provider string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	providerPath := filepath.Join(homeDir, ".claude-cli", provider, "settings.json")
	if _, err := os.Stat(providerPath); os.IsNotExist(err) {
		return "", fmt.Errorf("provider configuration not found: %s", providerPath)
	}

	// Read provider settings
	providerData, err := os.ReadFile(providerPath)
	if err != nil {
		return "", fmt.Errorf("failed to read provider settings: %w", err)
	}

	var providerSettings map[string]any
	if err := json.Unmarshal(providerData, &providerSettings); err != nil {
		return "", fmt.Errorf("failed to parse provider settings: %w", err)
	}

	// Read existing user settings (if any) and merge
	userSettingsPath := filepath.Join(homeDir, ".claude", "settings.json")
	if userData, err := os.ReadFile(userSettingsPath); err == nil {
		var userSettings map[string]any
		if err := json.Unmarshal(userData, &userSettings); err == nil {
			// Merge: user settings as base, provider settings override
			merged := mergeSettings(userSettings, providerSettings)

			// Write merged settings to a temp file
			tmpPath := filepath.Join(os.TempDir(), fmt.Sprintf("claude-cli-settings-%s.json", provider))
			mergedData, err := json.MarshalIndent(merged, "", "  ")
			if err != nil {
				return "", fmt.Errorf("failed to marshal settings: %w", err)
			}
			if err := os.WriteFile(tmpPath, mergedData, 0644); err != nil {
				return "", fmt.Errorf("failed to write settings: %w", err)
			}
			return tmpPath, nil
		}
	}

	// No user settings, just use provider settings
	tmpPath := filepath.Join(os.TempDir(), fmt.Sprintf("claude-cli-settings-%s.json", provider))
	prettyData, err := json.MarshalIndent(providerSettings, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal settings: %w", err)
	}
	if err := os.WriteFile(tmpPath, prettyData, 0644); err != nil {
		return "", fmt.Errorf("failed to write settings: %w", err)
	}
	return tmpPath, nil
}

func mergeSettings(base, overrides map[string]any) map[string]any {
	merged := make(map[string]any, len(base))
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range overrides {
		merged[k] = v
	}
	return merged
}
