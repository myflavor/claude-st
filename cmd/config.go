package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config <provider>",
	Short: "Switch Claude CLI provider configuration",
	Args:  cobra.ExactArgs(1),
	Run:   runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfig(cmd *cobra.Command, args []string) {
	provider := args[0]

	homeDir, err := os.UserHomeDir()
	if err != nil {
		exitWithError("Failed to get home directory")
	}

	srcDir := filepath.Join(homeDir, ".claude-cli", provider)
	dstDir := filepath.Join(homeDir, ".claude")

	// Check if source directory exists
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		exitWithError(fmt.Sprintf("Provider configuration not found: %s", srcDir))
	}

	// Copy all files from source to destination
	if err := copyDir(srcDir, dstDir); err != nil {
		exitWithError(fmt.Sprintf("Failed to copy config: %v", err))
	}

	fmt.Printf("Switched to %s configuration\n", provider)
}

func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		srcFile := filepath.Join(src, entry.Name())
		dstFile := filepath.Join(dst, entry.Name())

		data, err := os.ReadFile(srcFile)
		if err != nil {
			return err
		}

		if err := os.WriteFile(dstFile, data, 0644); err != nil {
			return err
		}
	}

	return nil
}
