package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "claude-cli",
	Short: "Claude CLI configuration manager",
	Long:  "Manage Claude CLI provider configurations",
}

func Execute() error {
	return rootCmd.Execute()
}

func exitWithError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
