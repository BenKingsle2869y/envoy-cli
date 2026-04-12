package crypto_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/envoy/internal/crypto"
)

func TestResolvePassphraseFromEnv(t *testing.T) {
	t.Setenv("ENVOY_PASSPHRASE", "env-passphrase")

	pass, err := crypto.ResolvePassphrase()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pass != "env-passphrase" {
		t.Fatalf("expected 'env-passphrase', got %q", pass)
	}
}

func TestSaveAndResolvePassphraseFromFile(t *testing.T) {
	// Ensure env var is not set
	os.Unsetenv("ENVOY_PASSPHRASE")

	tmpDir := t.TempDir()
	err := crypto.SavePassphraseToFile(tmpDir, "file-passphrase")
	if err != nil {
		t.Fatalf("SavePassphraseToFile failed: %v", err)
	}

	// Verify file permissions
	info, err := os.Stat(filepath.Join(tmpDir, ".envoy_key"))
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("expected file mode 0600, got %v", info.Mode().Perm())
	}
}

func TestResolvePassphraseNotFound(t *testing.T) {
	os.Unsetenv("ENVOY_PASSPHRASE")

	// Change to a temp dir with no key file
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)
	os.Chdir(tmpDir)

	_, err := crypto.ResolvePassphrase()
	if err == nil {
		t.Fatal("expected error when no passphrase is available")
	}
}
