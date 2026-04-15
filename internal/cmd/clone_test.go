package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"envoy-cli/internal/store"
)

// setupCloneStores creates a temporary directory for clone tests, sets the
// required environment variables, and returns the directory path along with
// a cleanup function that removes the directory.
func setupCloneStores(t *testing.T) (string, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "clone-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Setenv("ENVOY_HOME", tmpDir)
	t.Setenv("ENVOY_PASSPHRASE", "test-passphrase")
	return tmpDir, func() { os.RemoveAll(tmpDir) }
}

// saveContext is a test helper that saves a context store file under tmpDir.
func saveContext(t *testing.T, tmpDir, ctx string, data map[string]string) {
	t.Helper()
	p := filepath.Join(tmpDir, fmt.Sprintf("%s.env.enc", ctx))
	if err := store.Save(p, "test-passphrase", data); err != nil {
		t.Fatalf("Save %s: %v", ctx, err)
	}
}

func TestCloneCmd_ClonesContext(t *testing.T) {
	tmpDir, cleanup := setupCloneStores(t)
	defer cleanup()

	saveContext(t, tmpDir, "src", map[string]string{"FOO": "bar", "BAZ": "qux"})

	RootCmd.SetArgs([]string{"clone", "src", "dst"})
	var buf bytes.Buffer
	RootCmd.SetOut(&buf)
	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	dstPath := filepath.Join(tmpDir, "dst.env.enc")
	loaded, err := store.Load(dstPath, "test-passphrase")
	if err != nil {
		t.Fatalf("Load dst: %v", err)
	}
	if loaded["FOO"] != "bar" || loaded["BAZ"] != "qux" {
		t.Errorf("unexpected cloned data: %v", loaded)
	}
	if got := buf.String(); got == "" {
		t.Error("expected output, got empty")
	}
}

func TestCloneCmd_SameContextFails(t *testing.T) {
	_, cleanup := setupCloneStores(t)
	defer cleanup()

	RootCmd.SetArgs([]string{"clone", "same", "same"})
	if err := RootCmd.Execute(); err == nil {
		t.Error("expected error for same source and destination")
	}
}

func TestCloneCmd_FailsWhenDestinationExists(t *testing.T) {
	tmpDir, cleanup := setupCloneStores(t)
	defer cleanup()

	for _, ctx := range []string{"alpha", "beta"} {
		saveContext(t, tmpDir, ctx, map[string]string{"K": "v"})
	}

	RootCmd.SetArgs([]string{"clone", "alpha", "beta"})
	if err := RootCmd.Execute(); err == nil {
		t.Error("expected error when destination exists without --overwrite")
	}
}

func TestCloneCmd_OverwriteFlag(t *testing.T) {
	tmpDir, cleanup := setupCloneStores(t)
	defer cleanup()

	saveContext(t, tmpDir, "src", map[string]string{"NEW": "value"})
	saveContext(t, tmpDir, "dst", map[string]string{"OLD": "stale"})

	RootCmd.SetArgs([]string{"clone", "src", "dst", "--overwrite"})
	if err := RootCmd.Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	dstPath := filepath.Join(tmpDir, "dst.env.enc")
	loaded, _ := store.Load(dstPath, "test-passphrase")
	if loaded["NEW"] != "value" {
		t.Errorf("expected overwritten data, got %v", loaded)
	}
	if _, ok := loaded["OLD"]; ok {
		t.Error("expected OLD key to be gone after overwrite")
	}
}
