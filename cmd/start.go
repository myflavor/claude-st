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
	Long: `Start Claude CLI. Inject provider env vars and start claude.

Short flags (parsed by claude-cli):
  -P <name>    Inject provider env vars from ~/.claude-cli/<name>/settings.json
  -S           Skip permissions (adds -dangerously-skip-permissions)

All other arguments are passed directly to claude.

Examples:
  claude-cli start                           # Start claude directly
  claude-cli start -P minimax                # Inject minimax env and start
  claude-cli start -P minimax -S             # Inject env and skip permissions
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
		switch arg {
		case "-P":
			if i+1 >= len(args) {
				exitWithError("Flag -P requires a value")
			}
			provider = args[i+1]
			i++
		case "-S":
			skipPermissions = true
		case "-H", "--help":
			cmd.Help()
			return
		default:
			claudeArgs = append(claudeArgs, arg)
		}
	}

	// Prepare environment variables
	env := os.Environ()

	// Inject provider env vars if specified
	if provider != "" {
		providerEnv, err := loadProviderEnv(provider)
		if err != nil {
			exitWithError(err.Error())
		}
		env = mergeEnv(env, providerEnv)
		fmt.Printf("Injected env from provider: %s\n", provider)
	}

	// Prepare arguments for claude command
	finalArgs := make([]string, 0, len(claudeArgs)+1)
	if skipPermissions {
		finalArgs = append(finalArgs, "-dangerously-skip-permissions")
	}
	finalArgs = append(finalArgs, claudeArgs...)

	// Execute claude command with injected env
	claudeCmd := exec.Command("claude", finalArgs...)
	claudeCmd.Stdin = os.Stdin
	claudeCmd.Stdout = os.Stdout
	claudeCmd.Stderr = os.Stderr
	claudeCmd.Env = env

	fmt.Printf("Starting: claude %s\n", strings.Join(finalArgs, " "))

	if err := claudeCmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		exitWithError(fmt.Sprintf("Failed to execute claude: %v", err))
	}
}

func loadProviderEnv(provider string) (map[string]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	settingsPath := filepath.Join(homeDir, ".claude-cli", provider, "settings.json")

	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("provider configuration not found: %s", settingsPath)
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read settings: %w", err)
	}

	var settings struct {
		Env map[string]string `json:"env"`
	}
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse settings: %w", err)
	}

	return settings.Env, nil
}

func mergeEnv(base []string, overrides map[string]string) []string {
	envMap := make(map[string]string, len(base)+len(overrides))
	for _, e := range base {
		if idx := strings.Index(e, "="); idx > 0 {
			envMap[e[:idx]] = e[idx+1:]
		}
	}
	for k, v := range overrides {
		envMap[k] = v
	}

	merged := make([]string, 0, len(envMap))
	for k, v := range envMap {
		merged = append(merged, k+"="+v)
	}
	return merged
}
