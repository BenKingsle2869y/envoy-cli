package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"

	"envoy-cli/internal/store"
)

func TestRotateCmd_Success(t *testing.T) {
	dir := t.TempDir()
	storePath := filepath.Join(dir, "test.env.enc")

	currentPass := "old-secret"
	newPass := "new-secret"

	// Create initial store
	s := &store.Store{Envs: map[string]string{"FOO": "bar"}}
	if err := store.Save(storePath, s, currentPass); err != nil {
		t.Fatalf("setup: save store: %v", err)
	}

	// Set env vars for passphrases
	t.Setenv("ENVOY_PASSPHRASE", currentPass)
	t.Setenv("ENVOY_NEW_PASSPHRASE", newPass)

	cmd := &cobra.Command{}
	cmd.Flags().String("store", storePath, "")
	cmd.Flags().String("new-passphrase", "", "")

	if err := runRotate(cmd, nil); err != nil {
		t.Fatalf("runRotate: %v", err)
	}

	// Verify store is readable with new passphrase
	loaded, err := store.Load(storePath, newPass)
	if err != nil {
		t.Fatalf("load with new passphrase: %v", err)
	}
	if loaded.Envs["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", loaded.Envs["FOO"])
	}
}

func TestRotateCmd_FailsSamePassphrase(t *testing.T) {
	dir := t.TempDir()
	storePath := filepath.Join(dir, "test.env.enc")

	pass := "same-secret"
	s := &store.Store{Envs: map[string]string{}}
	if err := store.Save(storePath, s, pass); err != nil {
		t.Fatalf("setup: %v", err)
	}

	t.Setenv("ENVOY_PASSPHRASE", pass)
	t.Setenv("ENVOY_NEW_PASSPHRASE", pass)

	cmd := &cobra.Command{}
	cmd.Flags().String("store", storePath, "")
	cmd.Flags().String("new-passphrase", "", "")

	err := runRotate(cmd, nil)
	if err == nil {
		t.Fatal("expected error when passphrases are identical")
	}
}

func TestRotateCmd_FailsWrongCurrentPassphrase(t *testing.T) {
	dir := t.TempDir()
	storePath := filepath.Join(dir, "test.env.enc")

	s := &store.Store{Envs: map[string]string{}}
	if err := store.Save(storePath, s, "correct"); err != nil {
		t.Fatalf("setup: %v", err)
	}

	os.Setenv("ENVOY_PASSPHRASE", "wrong")
	os.Setenv("ENVOY_NEW_PASSPHRASE", "newpass")
	t.Cleanup(func() {
		os.Unsetenv("ENVOY_PASSPHRASE")
		os.Unsetenv("ENVOY_NEW_PASSPHRASE")
	})

	cmd := &cobra.Command{}
	cmd.Flags().String("store", storePath, "")
	cmd.Flags().String("new-passphrase", "", "")

	if err := runRotate(cmd, nil); err == nil {
		t.Fatal("expected error with wrong current passphrase")
	}
}
