// claude-st starts Claude Code with a provider configuration.
//
// Usage:
//
//	claude-st <provider> [claude args...]
//	claude-st <provider> -S [claude args...]
//
// <provider> is required and selects ~/.claude-cli/<provider>.json, which is
// passed to claude via --settings. -S adds --dangerously-skip-permissions.
// All other arguments are passed through to claude unchanged.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"claude-cli/internal/provider"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "claude-st <provider> [claude args...]",
	Short: "Start Claude CLI with a provider configuration",
	Long: `Start Claude CLI with a provider configuration.

<provider> is required and selects ~/.claude-cli/<provider>.json, which is
passed to claude via --settings. Multiple providers can run side by side in
different terminals without overwriting ~/.claude/settings.json.

Short flags (parsed by claude-st):
  -S    Skip permissions (adds --dangerously-skip-permissions)
  -H    Show help

All other arguments are passed directly to claude.

Examples:
  claude-st minimax                  # Start with minimax provider
  claude-st minimax -S               # ...and skip permissions
  claude-st minimax --model opus     # Pass args to claude
  claude-st minimax -S -c            # Skip permissions + continue
  claude-st -S deepseek              # Flags may precede the provider`,
	DisableFlagParsing: true,
	Run:                run,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	var providerName string
	skipPermissions := false
	claudeArgs := make([]string, 0, len(args))

	// Parse our short flags; the first non-flag token is the provider.
	// Everything else is passed through to claude.
	for _, arg := range args {
		switch arg {
		case "-S":
			skipPermissions = true
		case "-H", "--help":
			cmd.Help()
			return
		default:
			if providerName == "" && !strings.HasPrefix(arg, "-") {
				providerName = arg
				continue
			}
			claudeArgs = append(claudeArgs, arg)
		}
	}

	if providerName == "" {
		fmt.Fprintln(os.Stderr, "Error: provider is required")
		fmt.Fprintln(os.Stderr, "Usage: claude-st <provider> [claude args...]")
		os.Exit(1)
	}

	settingsPath, err := provider.File(providerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: provider configuration not found: %s\n", settingsPath)
		os.Exit(1)
	}

	// Build the claude argument list: --settings first, then optional -S, then passthrough.
	finalArgs := make([]string, 0, len(claudeArgs)+3)
	finalArgs = append(finalArgs, "--settings", settingsPath)
	if skipPermissions {
		finalArgs = append(finalArgs, "--dangerously-skip-permissions")
	}
	finalArgs = append(finalArgs, claudeArgs...)

	// LookPath handles .exe/.cmd/.bat resolution on Windows.
	claudeBin, err := exec.LookPath("claude")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: claude command not found in PATH: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Using provider: %s\n", providerName)

	claudeCmd := exec.Command(claudeBin, finalArgs...)
	claudeCmd.Stdin = os.Stdin
	claudeCmd.Stdout = os.Stdout
	claudeCmd.Stderr = os.Stderr

	if err := claudeCmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "Error: failed to execute claude: %v\n", err)
		os.Exit(1)
	}
}
