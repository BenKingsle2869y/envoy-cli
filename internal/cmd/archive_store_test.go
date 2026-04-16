package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/envoy/internal/store"
)

func setupArchiveDir(t *testing.T) (string, func()) {
	t.Helper()
	tmp, err := os.MkdirTemp("", "envoy-archive-*")
	if err != nil {
		t.Fatal(err)
	}
	old, _ := os.UserHomeDir()
	t.Setenv("HOME", tmp)
	_ = old
	return tmp, func() { os.RemoveAll(tmp) }
}

func sampleArchiveEntries() []store.Entry {
	return []store.Entry{
		{Key: "FOO", Value: "bar"},
		{Key: "BAZ", Value: "qux"},
	}
}

func TestCreateArchive_CreatesFile(t *testing.T) {
	tmp, cleanup := setupArchiveDir(t)
	defer cleanup()

	err := CreateArchive("default", "test-archive", sampleArchiveEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	path := filepath.Join(tmp, ".envoy", "archives", "default", "test-archive.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("archive file not created")
	}
}

func TestLoadArchives_ReturnsAll(t *testing.T) {
	_, cleanup := setupArchiveDir(t)
	defer cleanup()

	_ = CreateArchive("default", "arch-1", sampleArchiveEntries())
	_ = CreateArchive("default", "arch-2", sampleArchiveEntries())

	archives, err := LoadArchives("default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(archives) != 2 {
		t.Fatalf("expected 2 archives, got %d", len(archives))
	}
}

func TestLoadArchives_EmptyWhenNoneExist(t *testing.T) {
	_, cleanup := setupArchiveDir(t)
	defer cleanup()

	archives, err := LoadArchives("default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(archives) != 0 {
		t.Fatalf("expected 0 archives, got %d", len(archives))
	}
}

func TestRestoreArchive_RestoresEntries(t *testing.T) {
	_, cleanup := setupArchiveDir(t)
	defer cleanup()

	original := sampleArchiveEntries()
	_ = CreateArchive("default", "my-arch", original)

	entries, err := RestoreArchive("default", "my-arch")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != len(original) {
		t.Fatalf("expected %d entries, got %d", len(original), len(entries))
	}
	if entries[0].Key != "FOO" || entries[0].Value != "bar" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestRestoreArchive_NotFound(t *testing.T) {
	_, cleanup := setupArchiveDir(t)
	defer cleanup()

	_, err := RestoreArchive("default", "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing archive")
	}
}
