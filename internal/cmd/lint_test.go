package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeLintFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeLintFile: %v", err)
	}
	return p
}

func TestLintCmd_NoIssues(t *testing.T) {
	p := writeLintFile(t, "APP_NAME=envoy\nDEBUG=false\n# comment\n")

	var buf bytes.Buffer
	lintCmd.SetOut(&buf)
	lintCmd.SetErr(&buf)

	err := lintCmd.RunE(lintCmd, []string{p})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(buf.String(), "No issues found") {
		t.Errorf("expected 'No issues found', got: %q", buf.String())
	}
}

func TestLintCmd_MissingEquals(t *testing.T) {
	p := writeLintFile(t, "BADLINE\nGOOD=ok\n")

	var buf bytes.Buffer
	lintCmd.SetOut(&buf)

	err := lintCmd.RunE(lintCmd, []string{p})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
	if !strings.Contains(buf.String(), "missing '=' separator") {
		t.Errorf("expected missing separator message, got: %q", buf.String())
	}
}

func TestLintCmd_DuplicateKey(t *testing.T) {
	p := writeLintFile(t, "FOO=bar\nFOO=baz\n")

	var buf bytes.Buffer
	lintCmd.SetOut(&buf)

	err := lintCmd.RunE(lintCmd, []string{p})
	if err == nil {
		t.Fatal("expected error for duplicate key")
	}
	if !strings.Contains(buf.String(), "duplicate key") {
		t.Errorf("expected duplicate key message, got: %q", buf.String())
	}
}

func TestLintCmd_InvalidKey(t *testing.T) {
	p := writeLintFile(t, "123BAD=value\n")

	var buf bytes.Buffer
	lintCmd.SetOut(&buf)

	err := lintCmd.RunE(lintCmd, []string{p})
	if err == nil {
		t.Fatal("expected error for invalid key")
	}
	if !strings.Contains(buf.String(), "invalid key") {
		t.Errorf("expected invalid key message, got: %q", buf.String())
	}
}

func TestLintCmd_MissingFile(t *testing.T) {
	err := lintCmd.RunE(lintCmd, []string{"/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestIsValidEnvKey(t *testing.T) {
	cases := []struct {
		key   string
		valid bool
	}{
		{"APP_NAME", true},
		{"_PRIVATE", true},
		{"FOO123", true},
		{"123FOO", false},
		{"FOO-BAR", false},
		{"", false},
		{"FOO BAR", false},
	}
	for _, tc := range cases {
		got := isValidEnvKey(tc.key)
		if got != tc.valid {
			t.Errorf("isValidEnvKey(%q) = %v, want %v", tc.key, got, tc.valid)
		}
	}
}
