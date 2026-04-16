package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func setupSecretStore(t *testing.T) (string, func()) {
	t.Helper()
	dir := t.TempDir()
	old := os.Getenv("HOME")
	os.Setenv("HOME", dir)
	os.MkdirAll(filepath.Join(dir, ".envoy"), 0755)
	return dir, func() { os.Setenv("HOME", old) }
}

func execSecretCmd(args ...string) (string, error) {
	root := &cobra.Command{Use: "envoy"}
	root.AddCommand(secretCmd)
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)
	err := root.Execute()
	return buf.String(), err
}

func TestSecretMarkCmd_MarksKey(t *testing.T) {
	dir, cleanup := setupSecretStore(t)
	defer cleanup()
	sfPath := filepath.Join(dir, ".envoy", "secrets.json")

	_, err := execSecretCmd("secret", "mark", "API_KEY", "--passphrase", "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s, err := LoadSecrets(sfPath)
	if err != nil {
		t.Fatal(err)
	}
	if !s.Keys["API_KEY"] {
		t.Error("expected API_KEY to be marked")
	}
}

func TestSecretUnmarkCmd_UnmarksKey(t *testing.T) {
	dir, cleanup := setupSecretStore(t)
	defer cleanup()
	sfPath := filepath.Join(dir, ".envoy", "secrets.json")

	s := &SecretStore{Keys: map[string]bool{"TOKEN": true}}
	_ = SaveSecrets(sfPath, s)

	out, err := execSecretCmd("secret", "unmark", "TOKEN", "--passphrase", "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "unmarked") {
		t.Errorf("expected unmarked message, got %q", out)
	}
}

func TestSecretListCmd_ShowsKeys(t *testing.T) {
	dir, cleanup := setupSecretStore(t)
	defer cleanup()
	sfPath := filepath.Join(dir, ".envoy", "secrets.json")

	s := &SecretStore{Keys: map[string]bool{"DB_PASS": true, "API_KEY": true}}
	_ = SaveSecrets(sfPath, s)

	out, err := execSecretCmd("secret", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "DB_PASS") || !strings.Contains(out, "API_KEY") {
		t.Errorf("expected keys in output, got %q", out)
	}
}

func TestSecretListCmd_EmptyMessage(t *testing.T) {
	_, cleanup := setupSecretStore(t)
	defer cleanup()

	out, err := execSecretCmd("secret", "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No secret keys") {
		t.Errorf("expected empty message, got %q", out)
	}
}
