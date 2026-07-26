//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
)

// launch runs claude as a child process. Windows has no exec(2), so we spawn
// claude and proxy its exit code; the console is shared via stdio passthrough.
func launch(claudeBin string, args []string) {
	c := exec.Command(claudeBin, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "Error: failed to execute claude: %v\n", err)
		os.Exit(1)
	}
}
