// claude-cfg copies a provider configuration into ~/.claude/settings.json.
//
// Usage:
//
//	claude-cfg <provider>
//
// It copies ~/.claude-cli/<provider>.json to ~/.claude/settings.json,
// overwriting it. This is destructive: for non-destructive launches prefer
// claude-st <provider>, which passes the provider via --settings instead.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"claude-cli/internal/provider"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "claude-cfg <provider>",
	Short: "Copy a provider configuration into ~/.claude/settings.json",
	Long: `Copy a provider configuration into ~/.claude/settings.json.

This overwrites ~/.claude/settings.json. For non-destructive launches, prefer
claude-st <provider>, which passes the provider via --settings instead.

Example:
  claude-cfg minimax    # copies ~/.claude-cli/minimax.json -> ~/.claude/settings.json`,
	Args: cobra.ExactArgs(1),
	Run:  run,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	name := args[0]

	src, err := provider.File(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if _, err := os.Stat(src); err != nil {
		fmt.Fprintf(os.Stderr, "Error: provider configuration not accessible: %s: %v\n", src, err)
		os.Exit(1)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to get home directory: %v\n", err)
		os.Exit(1)
	}
	dstDir := filepath.Join(home, ".claude")
	dst := filepath.Join(dstDir, "settings.json")

	if err := os.MkdirAll(dstDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create %s: %v\n", dstDir, err)
		os.Exit(1)
	}

	data, err := os.ReadFile(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to read %s: %v\n", src, err)
		os.Exit(1)
	}
	// Write to a temp file in the same directory, then rename over dst.
	// os.WriteFile truncates dst before writing; a partial write (ENOSPC,
	// crash) would leave settings.json truncated with no backup. Rename is
	// atomic on the same filesystem.
	tmp := dst + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to write %s: %v\n", tmp, err)
		os.Exit(1)
	}
	if err := os.Rename(tmp, dst); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to install %s: %v\n", dst, err)
		os.Exit(1)
	}

	fmt.Printf("Copied %s -> %s\n", src, dst)
}
