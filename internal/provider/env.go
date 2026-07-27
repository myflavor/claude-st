package provider

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EnvFile returns the path to the global env file (~/.claude-cli/.env).
// It is optional: if the file does not exist, claude-st skips it silently.
func EnvFile() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ".env"), nil
}

// LoadEnv parses a dotenv file at path and exports each KEY=VALUE into the
// process environment. Variables already set in the environment are NOT
// overridden, so shell exports take precedence over .env defaults.
//
// Supported syntax: blank lines and # comments; an optional "export " prefix;
// surrounding single or double quotes around the value (so values may contain
// '#' or spaces).
//
// The file is parsed fully before any variable is applied: on a parse error
// nothing is exported, so claude is never launched with a partially-applied
// environment. It returns the count of variables parsed (including ones
// skipped because they were already set).
func LoadEnv(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	type entry struct{ key, val string; line int }
	var entries []entry
	sc := bufio.NewScanner(f)
	ln := 0
	for sc.Scan() {
		ln++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		eq := strings.IndexByte(line, '=')
		if eq < 0 {
			return 0, fmt.Errorf("%s:%d: missing '='", path, ln)
		}
		key := strings.TrimSpace(line[:eq])
		val := strings.TrimSpace(line[eq+1:])
		if key == "" {
			return 0, fmt.Errorf("%s:%d: empty key", path, ln)
		}
		if len(val) >= 2 && (val[0] == '"' && val[len(val)-1] == '"' || val[0] == '\'' && val[len(val)-1] == '\'') {
			val = val[1 : len(val)-1]
		}
		entries = append(entries, entry{key, val, ln})
	}
	if err := sc.Err(); err != nil {
		return 0, err
	}
	// Apply only after the whole file parsed cleanly.
	n := 0
	for _, e := range entries {
		if _, ok := os.LookupEnv(e.key); !ok {
			if err := os.Setenv(e.key, e.val); err != nil {
				return n, fmt.Errorf("%s:%d: %w", path, e.line, err)
			}
		}
		n++
	}
	return n, nil
}
