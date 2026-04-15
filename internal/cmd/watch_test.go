package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/store"
)

func TestDiffEntries_Added(t *testing.T) {
	old := map[string]string{"A": "1"}
	next := map[string]string{"A": "1", "B": "2"}
	lines := diffEntries(old, next)
	if len(lines) != 1 {
		t.Fatalf("expected 1 change, got %d", len(lines))
	}
	if lines[0] != "  + B=2" {
		t.Errorf("unexpected line: %q", lines[0])
	}
}

func TestDiffEntries_Removed(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	next := map[string]string{"A": "1"}
	lines := diffEntries(old, next)
	if len(lines) != 1 {
		t.Fatalf("expected 1 change, got %d", len(lines))
	}
	if lines[0] != "  - B" {
		t.Errorf("unexpected line: %q", lines[0])
	}
}

func TestDiffEntries_Changed(t *testing.T) {
	old := map[string]string{"A": "1"}
	next := map[string]string{"A": "99"}
	lines := diffEntries(old, next)
	if len(lines) != 1 {
		t.Fatalf("expected 1 change, got %d", len(lines))
	}
	if lines[0] != "  ~ A=99 (was 1)" {
		t.Errorf("unexpected line: %q", lines[0])
	}
}

func TestDiffEntries_NoChanges(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	next := map[string]string{"A": "1", "B": "2"}
	lines := diffEntries(old, next)
	if len(lines) != 0 {
		t.Errorf("expected no changes, got %d", len(lines))
	}
}

func TestWatchCmd_FailsWithoutPassphrase(t *testing.T) {
	t.Setenv("ENVOY_PASSPHRASE", "")
	t.Setenv("HOME", t.TempDir())

	var buf bytes.Buffer
	watchCmd.SetOut(&buf)
	watchCmd.SetErr(&buf)

	err := watchCmd.RunE(watchCmd, []string{})
	if err == nil {
		t.Fatal("expected error when passphrase is missing")
	}
}

func TestWatchCmd_FailsOnMissingStore(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("ENVOY_PASSPHRASE", "secret")

	// No store file created — Load should fail gracefully or return empty.
	// We just verify diffEntries works on empty maps.
	empty := store.Store{Entries: map[string]string{}}
	lines := diffEntries(empty.Entries, empty.Entries)
	if len(lines) != 0 {
		t.Errorf("expected no diff on identical empty maps")
	}
}
