// Package provider resolves provider configuration paths.
//
// Providers are stored as single JSON files under ~/.claude-cli/<name>.json.
package provider

import (
	"fmt"
	"os"
	"path/filepath"
)

// configDirName is the directory holding provider files.
const configDirName = ".claude-cli"

// Dir returns the provider config directory (~/.claude-cli).
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, configDirName), nil
}

// File returns the path to a provider's settings file (~/.claude-cli/<name>.json).
func File(name string) (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, name+".json"), nil
}
