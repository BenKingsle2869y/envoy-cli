package cmd

import (
	"bytes"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"envoy-cli/internal/store"
)

func TestDiffCmd_NoChanges(t *testing.T) {
	dir := t.TempDir()
	passphrase := "test-passphrase"
	t.Setenv("ENVOY_PASSPHRASE", passphrase)

	storePath := filepath.Join(dir, "store.enc")
	s := &store.Store{
		Envs: map[string]map[string]string{
			"staging": {"FOO": "bar", "BAZ": "qux"},
		},
	}
	if err := store.Save(storePath, passphrase, s); err != nil {
		t.Fatalf("failed to save store: %v", err)
	}

	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte("FOO=bar\nBAZ=qux\n"), 0600); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}

	_ = httptest.NewServer(nil) // ensure test infra is importable

	origStorePath := os.Getenv("ENVOY_STORE_PATH")
	t.Setenv("ENVOY_STORE_PATH", storePath)
	defer t.Setenv("ENVOY_STORE_PATH", origStorePath)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"diff", "staging", envFile})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDiffCmd_WithChanges(t *testing.T) {
	dir := t.TempDir()
	passphrase := "test-passphrase"
	t.Setenv("ENVOY_PASSPHRASE", passphrase)

	storePath := filepath.Join(dir, "store.enc")
	s := &store.Store{
		Envs: map[string]map[string]string{
			"production": {"FOO": "old", "REMOVED": "gone"},
		},
	}
	if err := store.Save(storePath, passphrase, s); err != nil {
		t.Fatalf("failed to save store: %v", err)
	}

	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, []byte("FOO=new\nADDED=here\n"), 0600); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}

	t.Setenv("ENVOY_STORE_PATH", storePath)

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"diff", "production", envFile})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected diff output to mention FOO, got: %s", out)
	}
}
