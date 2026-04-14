package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func writeValidateFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeValidateFile: %v", err)
	}
	return p
}

func execValidateCmd(t *testing.T, args []string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	root := &cobra.Command{Use: "envoy"}
	validateCmd.ResetFlags()
	validateCmd.Flags().StringSliceP("require", "r", nil, "")
	validateCmd.Flags().BoolP("strict", "s", false, "")
	root.AddCommand(validateCmd)
	root.SetOut(buf)
	validateCmd.SetOut(buf)
	root.SetArgs(append([]string{"validate"}, args...))
	err := root.Execute()
	return buf.String(), err
}

func TestValidateCmd_PassesWhenRequiredKeysPresent(t *testing.T) {
	p := writeValidateFile(t, "DB_HOST=localhost\nDB_PORT=5432\n")
	out, err := execValidateCmd(t, []string{p, "--require", "DB_HOST,DB_PORT"})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "Validation passed") {
		t.Errorf("expected success message, got: %q", out)
	}
}

func TestValidateCmd_FailsWhenRequiredKeyMissing(t *testing.T) {
	p := writeValidateFile(t, "DB_HOST=localhost\n")
	out, err := execValidateCmd(t, []string{p, "--require", "DB_HOST,DB_PORT"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(out, "missing required key: DB_PORT") {
		t.Errorf("expected missing key message, got: %q", out)
	}
}

func TestValidateCmd_StrictRejectsExtraKeys(t *testing.T) {
	p := writeValidateFile(t, "DB_HOST=localhost\nEXTRA=value\n")
	out, err := execValidateCmd(t, []string{p, "--require", "DB_HOST", "--strict"})
	if err == nil {
		t.Fatal("expected error in strict mode, got nil")
	}
	if !strings.Contains(out, "unexpected key in strict mode: EXTRA") {
		t.Errorf("expected unexpected key message, got: %q", out)
	}
}

func TestValidateCmd_PassesWithNoRequiredKeys(t *testing.T) {
	p := writeValidateFile(t, "FOO=bar\n")
	out, err := execValidateCmd(t, []string{p})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if !strings.Contains(out, "Validation passed") {
		t.Errorf("expected success message, got: %q", out)
	}
}

func TestValidateCmd_FailsOnMissingFile(t *testing.T) {
	_, err := execValidateCmd(t, []string{"/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
