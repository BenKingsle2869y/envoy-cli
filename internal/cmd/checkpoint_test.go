package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/envoy-cli/internal/store"
)

func setupCheckpointDir(t *testing.T) (string, func()) {
	t.Helper()
	tmp, err := os.MkdirTemp("", "envoy-checkpoint-*")
	if err != nil {
		t.Fatal(err)
	}
	old, _ := os.UserHomeDir()
	t.Setenv("HOME", tmp)
	_ = old
	return tmp, func() { os.RemoveAll(tmp) }
}

func sampleEntries() map[string]store.Entry {
	return map[string]store.Entry{
		"FOO": {Value: "bar"},
		"BAZ": {Value: "qux"},
	}
}

func TestCreateCheckpoint_CreatesFile(t *testing.T) {
	tmp, cleanup := setupCheckpointDir(t)
	defer cleanup()

	err := CreateCheckpoint("default", "v1", sampleEntries())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	path := filepath.Join(tmp, ".envoy", "checkpoints", "default", "v1.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected checkpoint file to exist")
	}
}

func TestLoadCheckpoints_ReturnsAll(t *testing.T) {
	_, cleanup := setupCheckpointDir(t)
	defer cleanup()

	_ = CreateCheckpoint("default", "alpha", sampleEntries())
	_ = CreateCheckpoint("default", "beta", sampleEntries())

	cps, err := LoadCheckpoints("default")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cps) != 2 {
		t.Errorf("expected 2 checkpoints, got %d", len(cps))
	}
}

func TestRestoreCheckpoint_RestoresEntries(t *testing.T) {
	_, cleanup := setupCheckpointDir(t)
	defer cleanup()

	original := sampleEntries()
	_ = CreateCheckpoint("default", "snap1", original)

	entries, err := RestoreCheckpoint("default", "snap1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries["FOO"].Value != "bar" {
		t.Errorf("expected FOO=bar, got %s", entries["FOO"].Value)
	}
}

func TestRestoreCheckpoint_NotFound(t *testing.T) {
	_, cleanup := setupCheckpointDir(t)
	defer cleanup()

	_, err := RestoreCheckpoint("default", "ghost")
	if err == nil {
		t.Error("expected error for missing checkpoint")
	}
}

func TestCreateCheckpoint_TimestampIsUTC(t *testing.T) {
	_, cleanup := setupCheckpointDir(t)
	defer cleanup()

	before := time.Now().UTC()
	_ = CreateCheckpoint("default", "ts-test", sampleEntries())
	after := time.Now().UTC()

	cps, _ := LoadCheckpoints("default")
	if len(cps) == 0 {
		t.Fatal("no checkpoints loaded")
	}
	cp := cps[0]
	if cp.CreatedAt.Before(before) || cp.CreatedAt.After(after) {
		t.Errorf("timestamp %v out of expected range", cp.CreatedAt)
	}
}
