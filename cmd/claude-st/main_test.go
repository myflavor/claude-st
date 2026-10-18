package main

import "testing"

func TestHasFlag(t *testing.T) {
	cases := []struct {
		name string
		args []string
		flag string
		want bool
	}{
		{"exact", []string{"-S", "--append-system-prompt-file", "x"}, "--append-system-prompt-file", true},
		{"equals", []string{"--append-system-prompt-file=x"}, "--append-system-prompt-file", true},
		{"absent", []string{"--append-system-prompt", "x"}, "--append-system-prompt-file", false},
		{"prefix only", []string{"--append-system-prompt-file-extra"}, "--append-system-prompt-file", false},
		{"empty", nil, "--append-system-prompt-file", false},
	}
	for _, tc := range cases {
		if got := hasFlag(tc.args, tc.flag); got != tc.want {
			t.Errorf("%s: hasFlag(%q) = %v, want %v", tc.name, tc.args, got, tc.want)
		}
	}
}
