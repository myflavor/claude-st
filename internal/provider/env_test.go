package provider

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	content := `# a comment
BLANK=
SIMPLE=hello
QUOTED="with spaces and # hash"
SINGLE='single'
EXPORTED=exported
  SPACED = spaced
PRE_SET=from-env
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	// PRE_SET already exists in the environment; .env must NOT override it.
	os.Setenv("PRE_SET", "from-shell")
	cleanup := func() {
		for _, k := range []string{"PRE_SET", "BLANK", "SIMPLE", "QUOTED", "SINGLE", "EXPORTED", "SPACED"} {
			os.Unsetenv(k)
		}
	}
	t.Cleanup(cleanup)

	n, err := LoadEnv(path)
	if err != nil {
		t.Fatalf("LoadEnv error: %v", err)
	}
	if n != 7 {
		t.Errorf("parsed %d vars, want 7", n)
	}
	if got := os.Getenv("SIMPLE"); got != "hello" {
		t.Errorf("SIMPLE=%q want hello", got)
	}
	if got := os.Getenv("QUOTED"); got != "with spaces and # hash" {
		t.Errorf("QUOTED=%q want 'with spaces and # hash'", got)
	}
	if got := os.Getenv("SINGLE"); got != "single" {
		t.Errorf("SINGLE=%q want single", got)
	}
	if got := os.Getenv("SPACED"); got != "spaced" {
		t.Errorf("SPACED=%q want spaced", got)
	}
	if got := os.Getenv("PRE_SET"); got != "from-shell" {
		t.Errorf("PRE_SET=%q want from-shell (shell must override .env)", got)
	}
}

func TestLoadEnvMissingEquals(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("BADLINE\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadEnv(path); err == nil {
		t.Fatal("expected error for line without '='")
	}
}

func TestLoadEnvMissingFile(t *testing.T) {
	// The env file is optional: a missing file is not an error.
	n, err := LoadEnv(filepath.Join(t.TempDir(), "nope.env"))
	if err != nil {
		t.Fatalf("missing .env should be optional, got error: %v", err)
	}
	if n != 0 {
		t.Errorf("n=%d want 0 for missing file", n)
	}
}

func TestLoadEnvAtomicOnError(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	// GOOD_KEY appears before the malformed line; it must NOT be applied
	// because the whole file failed to parse.
	content := "GOOD_KEY=should_not_apply\nBADLINE\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Unsetenv("GOOD_KEY") })

	n, err := LoadEnv(path)
	if err == nil {
		t.Fatal("expected error for line without '='")
	}
	if n != 0 {
		t.Errorf("n=%d want 0 (nothing applied on parse error)", n)
	}
	if got := os.Getenv("GOOD_KEY"); got != "" {
		t.Errorf("GOOD_KEY=%q want empty (must not apply on parse error)", got)
	}
}
