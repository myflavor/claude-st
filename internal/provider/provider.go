// Package provider resolves provider configuration paths.
//
// Providers are stored as single JSON files under ~/.claude-cli/<name>.json.
// A shared MCP config may live at ~/.claude-cli/.mcp.json and is injected for
// every provider start.
package provider

import (
	"fmt"
	"os"
	"path/filepath"
)

// configDirName is the directory holding provider files.
const configDirName = ".claude-cli"

// mcpFileName is the global MCP config shared by all providers.
const mcpFileName = ".mcp.json"

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

// Resolve returns the provider config path, erroring if it is not accessible.
// Combines File with an existence check so callers don't repeat the stat+error.
func Resolve(name string) (string, error) {
	path, err := File(name)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("provider configuration not accessible: %s: %w", path, err)
	}
	return path, nil
}

// MCPFile returns the path to the global MCP config (~/.claude-cli/.mcp.json).
func MCPFile() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, mcpFileName), nil
}
