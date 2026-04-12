package store_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/envoy/internal/store"
)

const testPassphrase = "test-passphrase-123"

func tempStorePath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), ".envoy")
}

func TestLoadReturnsEmptyStoreWhenFileNotExist(t *testing.T) {
	path := tempStorePath(t)
	s, err := store.Load(path, testPassphrase)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(s.Environments) != 0 {
		t.Errorf("expected empty environments, got %d entries", len(s.Environments))
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := tempStorePath(t)

	s := &store.EnvStore{Environments: make(map[string]string)}
	s.PutEnv("production", map[string]string{
		"DB_HOST": "db.example.com",
		"DB_PORT": "5432",
	})

	if err := store.Save(path, testPassphrase, s); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := store.Load(path, testPassphrase)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	vars, err := loaded.GetEnv("production")
	if err != nil {
		t.Fatalf("GetEnv failed: %v", err)
	}

	if vars["DB_HOST"] != "db.example.com" {
		t.Errorf("expected DB_HOST=db.example.com, got %q", vars["DB_HOST"])
	}
	if vars["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", vars["DB_PORT"])
	}
}

func TestLoadFailsWithWrongPassphrase(t *testing.T) {
	path := tempStorePath(t)

	s := &store.EnvStore{Environments: make(map[string]string)}
	s.PutEnv("dev", map[string]string{"KEY": "value"})

	if err := store.Save(path, testPassphrase, s); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	_, err := store.Load(path, "wrong-passphrase")
	if err == nil {
		t.Error("expected error when loading with wrong passphrase, got nil")
	}
}

func TestGetEnvNotFound(t *testing.T) {
	s := &store.EnvStore{Environments: make(map[string]string)}
	_, err := s.GetEnv("nonexistent")
	if err == nil {
		t.Error("expected error for missing environment, got nil")
	}
}

func TestSaveCreatesFileWithRestrictedPermissions(t *testing.T) {
	path := tempStorePath(t)
	s := &store.EnvStore{Environments: make(map[string]string)}
	s.PutEnv("staging", map[string]string{"APP_ENV": "staging"})

	if err := store.Save(path, testPassphrase, s); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected file mode 0600, got %v", info.Mode().Perm())
	}
}
