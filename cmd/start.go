package cmd

import (
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
  -S           Skip permissions (adds --dangerously-skip-permissions)

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

	// Prepare arguments for claude command
	finalArgs := make([]string, 0, len(claudeArgs)+2)

	// If provider is specified, pass it via --settings
	if provider != "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			exitWithError(fmt.Sprintf("Failed to get home directory: %v", err))
		}
		settingsPath := filepath.Join(homeDir, ".claude-cli", provider, "settings.json")
		if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
			exitWithError(fmt.Sprintf("Provider configuration not found: %s", settingsPath))
		}
		finalArgs = append(finalArgs, "--settings", settingsPath)
		fmt.Printf("Using provider: %s\n", provider)
	}

	// Add skip permissions flag if set
	if skipPermissions {
		finalArgs = append(finalArgs, "--dangerously-skip-permissions")
	}

	finalArgs = append(finalArgs, claudeArgs...)

	// Locate claude binary (LookPath handles .exe/.cmd/.bat on Windows)
	claudeBin, err := exec.LookPath("claude")
	if err != nil {
		exitWithError(fmt.Sprintf("claude command not found in PATH: %v", err))
	}

	claudeCmd := exec.Command(claudeBin, finalArgs...)
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
