package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/your-org/envoy-cli/internal/store"
)

func setupSnapshotStore(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	passphrase := "snap-secret"
	t.Setenv("ENVOY_PASSPHRASE", passphrase)
	storePath := filepath.Join(dir, "default.env.enc")
	st, _ := store.Load(storePath, passphrase)
	st.Set("FOO", "bar")
	st.Set("BAZ", "qux")
	_ = store.Save(storePath, st, passphrase)
	return storePath, passphrase
}

func TestSnapshotCreate_CreatesFile(t *testing.T) {
	storePath, _ := setupSnapshotStore(t)
	st, _ := store.Load(storePath, "snap-secret")
	snap, err := CreateSnapshot(storePath, st.Entries(), "test-label", "snap-secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.ID == "" {
		t.Fatal("expected non-empty snapshot ID")
	}
	if snap.Label != "test-label" {
		t.Errorf("expected label 'test-label', got %q", snap.Label)
	}
	if _, err := os.Stat(snapshotFilePath(storePath, snap.ID)); err != nil {
		t.Errorf("snapshot file not found: %v", err)
	}
}

func TestSnapshotList_ReturnsSnapshots(t *testing.T) {
	storePath, _ := setupSnapshotStore(t)
	st, _ := store.Load(storePath, "snap-secret")
	_, _ = CreateSnapshot(storePath, st.Entries(), "first", "snap-secret")
	_, _ = CreateSnapshot(storePath, st.Entries(), "second", "snap-secret")
	snaps, err := LoadSnapshots(storePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snaps) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(snaps))
	}
}

func TestSnapshotRestore_RestoresData(t *testing.T) {
	storePath, passphrase := setupSnapshotStore(t)
	st, _ := store.Load(storePath, passphrase)
	snap, _ := CreateSnapshot(storePath, st.Entries(), "", passphrase)
	// Modify the store after snapshot
	st.Set("NEW_KEY", "new_val")
	_ = store.Save(storePath, st, passphrase)
	// Restore
	if err := RestoreSnapshot(storePath, snap.ID, passphrase); err != nil {
		t.Fatalf("unexpected restore error: %v", err)
	}
	restored, _ := store.Load(storePath, passphrase)
	if v, _ := restored.Get("FOO"); v != "bar" {
		t.Errorf("expected FOO=bar after restore, got %q", v)
	}
	// Ensure key added after snapshot is no longer present
	if v, ok := restored.Get("NEW_KEY"); ok {
		t.Errorf("expected NEW_KEY to be absent after restore, got %q", v)
	}
}

func TestSnapshotRestore_NotFound(t *testing.T) {
	storePath, passphrase := setupSnapshotStore(t)
	err := RestoreSnapshot(storePath, "nonexistent", passphrase)
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestSnapshotListCmd_Empty(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("ENVOY_PASSPHRASE", "pass")
	var buf bytes.Buffer
	snapshotListCmd.SetOut(&buf)
	snapshotListCmd.SetErr(&buf)
	_ = snapshotListCmd.RunE(snapshotListCmd, []string{})
	if !strings.Contains(buf.String(), "No snapshots") {
		t.Errorf("expected 'No snapshots' message, got: %s", buf.String())
	}
}
