package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"envoy-cli/internal/store"
)

func TestInitCmd_CreatesStore(t *testing.T) {
	dir := t.TempDir()
	storePath := filepath.Join(dir, ".envoy", "store.enc")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	t.Cleanup(func() { rootCmd.SetOut(nil) })

	rootCmd.SetArgs([]string{
		"init",
		"--passphrase", "test-secret",
		"--output", storePath,
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(storePath); os.IsNotExist(err) {
		t.Fatal("expected store file to exist after init")
	}

	env, err := store.Load(storePath, "test-secret")
	if err != nil {
		t.Fatalf("failed to load initialised store: %v", err)
	}
	if len(env) != 0 {
		t.Errorf("expected empty store, got %d entries", len(env))
	}

	output := buf.String()
	if output == "" {
		t.Error("expected success message on stdout")
	}
}

func TestInitCmd_FailsWhenStoreExists(t *testing.T) {
	dir := t.TempDir()
	storePath := filepath.Join(dir, "store.enc")

	// pre-create the file
	if err := store.Save(storePath, "secret", map[string]string{}); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	rootCmd.SetArgs([]string{
		"init",
		"--passphrase", "secret",
		"--output", storePath,
	})

	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error when store already exists, got nil")
	}
}

func TestInitCmd_FailsWithoutPassphrase(t *testing.T) {
	dir := t.TempDir()
	storePath := filepath.Join(dir, "store.enc")

	os.Unsetenv("ENVOY_PASSPHRASE")

	rootCmd.SetArgs([]string{
		"init",
		"--output", storePath,
	})

	if err := rootCmd.Execute(); err == nil {
		t.Fatal("expected error without passphrase, got nil")
	}
}
