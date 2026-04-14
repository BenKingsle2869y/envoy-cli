package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/internal/store"
)

func setupMergeStores(t *testing.T) (string, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "merge-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Setenv("ENVOY_HOME", tmpDir)
	t.Setenv("ENVOY_PASSPHRASE", "testpass")
	return tmpDir, func() { os.RemoveAll(tmpDir) }
}

func TestMergeCmd_MergesNewKeys(t *testing.T) {
	tmpDir, cleanup := setupMergeStores(t)
	defer cleanup()

	srcPath := filepath.Join(tmpDir, "src.env.enc")
	dstPath := filepath.Join(tmpDir, "dst.env.enc")

	srcStore := &store.Store{Entries: map[string]store.Entry{"FOO": {Value: "bar"}, "BAZ": {Value: "qux"}}}
	dstStore := &store.Store{Entries: map[string]store.Entry{"EXISTING": {Value: "keep"}}}

	_ = store.Save(srcPath, "testpass", srcStore)
	_ = store.Save(dstPath, "testpass", dstStore)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"merge", "src", "dst"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, err := store.Load(dstPath, "testpass")
	if err != nil {
		t.Fatalf("failed to load dst store: %v", err)
	}
	if _, ok := loaded.Entries["FOO"]; !ok {
		t.Error("expected FOO to be merged into dst")
	}
	if _, ok := loaded.Entries["EXISTING"]; !ok {
		t.Error("expected EXISTING to be preserved in dst")
	}
}

func TestMergeCmd_SkipsExistingWithoutOverwrite(t *testing.T) {
	tmpDir, cleanup := setupMergeStores(t)
	defer cleanup()

	srcPath := filepath.Join(tmpDir, "src.env.enc")
	dstPath := filepath.Join(tmpDir, "dst.env.enc")

	srcStore := &store.Store{Entries: map[string]store.Entry{"KEY": {Value: "from-src"}}}
	dstStore := &store.Store{Entries: map[string]store.Entry{"KEY": {Value: "from-dst"}}}

	_ = store.Save(srcPath, "testpass", srcStore)
	_ = store.Save(dstPath, "testpass", dstStore)

	rootCmd.SetArgs([]string{"merge", "src", "dst"})
	_ = rootCmd.Execute()

	loaded, _ := store.Load(dstPath, "testpass")
	if loaded.Entries["KEY"].Value != "from-dst" {
		t.Errorf("expected KEY to remain 'from-dst', got %q", loaded.Entries["KEY"].Value)
	}
}

func TestMergeCmd_OverwriteFlag(t *testing.T) {
	tmpDir, cleanup := setupMergeStores(t)
	defer cleanup()

	srcPath := filepath.Join(tmpDir, "src.env.enc")
	dstPath := filepath.Join(tmpDir, "dst.env.enc")

	srcStore := &store.Store{Entries: map[string]store.Entry{"KEY": {Value: "from-src"}}}
	dstStore := &store.Store{Entries: map[string]store.Entry{"KEY": {Value: "from-dst"}}}

	_ = store.Save(srcPath, "testpass", srcStore)
	_ = store.Save(dstPath, "testpass", dstStore)

	rootCmd.SetArgs([]string{"merge", "src", "dst", "--overwrite"})
	_ = rootCmd.Execute()

	loaded, _ := store.Load(dstPath, "testpass")
	if loaded.Entries["KEY"].Value != "from-src" {
		t.Errorf("expected KEY to be overwritten with 'from-src', got %q", loaded.Entries["KEY"].Value)
	}
}

func TestMergeCmd_SameContextFails(t *testing.T) {
	_, cleanup := setupMergeStores(t)
	defer cleanup()

	rootCmd.SetArgs([]string{"merge", "prod", "prod"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when source and destination are the same")
	}
}
