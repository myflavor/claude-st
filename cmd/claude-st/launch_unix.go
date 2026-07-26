//go:build !windows

package main

import (
	"fmt"
	"os"
	"syscall"
)

// launch replaces this process with claude via exec(2). claude becomes the
// direct foreground process, exactly as if the user had run `claude` directly,
// so terminal control (raw mode, mouse tracking, signal delivery) stays with
// claude instead of running it behind a Go parent that only injects args.
func launch(claudeBin string, args []string) {
	argv := append([]string{claudeBin}, args...)
	if err := syscall.Exec(claudeBin, argv, os.Environ()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to exec claude: %v\n", err)
		os.Exit(1)
	}
}
