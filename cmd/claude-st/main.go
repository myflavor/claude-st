// claude-st starts Claude Code with a provider configuration from ~/.claude-st.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"claude-st/internal/provider"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "claude-st <provider> [claude args...]",
	Short: "Start Claude CLI with a provider configuration",
	Long: `Start Claude CLI with a provider configuration.

<provider> is required and selects ~/.claude-st/<provider>.json, which is
passed to claude via --settings. Multiple providers can run side by side in
different terminals without overwriting ~/.claude/settings.json.

If ~/.claude-st/.mcp.json exists, it is automatically passed as --mcp-config
so every provider gets the same global MCP servers (session-scoped; does not
write project .mcp.json).

If ~/.claude-st/.env exists, each KEY=VALUE is exported into the environment
before launching claude. Variables already set in the shell are not overridden,
so common settings like CLAUDE_CODE_NO_FLICKER=1 apply to every provider
without repeating them in each settings file.

If ~/.claude-st/AGENTS.md exists, it is passed via --append-system-prompt-file
so its instructions reach every provider session without being written into a
project file.

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
			// Only claim -H/--help before the provider; afterwards forward to claude.
			if providerName == "" {
				cmd.Help()
				return
			}
			claudeArgs = append(claudeArgs, arg)
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

	settingsPath, err := provider.Resolve(providerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Optional global MCP shared by all providers.
	mcpPath := existingFile(provider.MCPFile())

	// claude rejects two appended prompts; the user's flags win.
	var agentsMdPath string
	if !hasFlag(claudeArgs, "--append-system-prompt") &&
		!hasFlag(claudeArgs, "--append-system-prompt-file") {
		agentsMdPath = existingFile(provider.AgentsMdFile())
	}

	finalArgs := make([]string, 0, len(claudeArgs)+7)
	finalArgs = append(finalArgs, "--settings", settingsPath)
	if mcpPath != "" {
		finalArgs = append(finalArgs, "--mcp-config", mcpPath)
	}
	if agentsMdPath != "" {
		finalArgs = append(finalArgs, "--append-system-prompt-file", agentsMdPath)
	}
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
	if mcpPath != "" {
		fmt.Printf("Using MCP config: %s\n", mcpPath)
	}
	if agentsMdPath != "" {
		fmt.Printf("Using instructions: %s\n", agentsMdPath)
	}

	if envPath, err := provider.EnvFile(); err == nil {
		if n, lerr := provider.LoadEnv(envPath); lerr != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid env file %s: %v\n", envPath, lerr)
			os.Exit(1)
		} else if n > 0 {
			fmt.Printf("Using env: %s\n", envPath)
		}
	}

	// launch replaces this process (Unix) or spawns a child (Windows); see launch_*.go.
	launch(claudeBin, finalArgs)
}

func hasFlag(args []string, flag string) bool {
	for _, arg := range args {
		if arg == flag || strings.HasPrefix(arg, flag+"=") {
			return true
		}
	}
	return false
}

func existingFile(path string, err error) string {
	if err != nil {
		return ""
	}
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	return path
}
