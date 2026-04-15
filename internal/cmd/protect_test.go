package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupProtectStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	return filepath.Join(dir, "default.env.enc")
}

func execProtectCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestMarkProtected_AddsKey(t *testing.T) {
	ps := &ProtectedStore{Keys: make(map[string]bool)}
	err := MarkProtected(ps, "SECRET")
	require.NoError(t, err)
	assert.True(t, ps.Keys["SECRET"])
}

func TestMarkProtected_RejectsDuplicate(t *testing.T) {
	ps := &ProtectedStore{Keys: map[string]bool{"SECRET": true}}
	err := MarkProtected(ps, "SECRET")
	assert.Error(t, err)
}

func TestUnmarkProtected_RemovesKey(t *testing.T) {
	ps := &ProtectedStore{Keys: map[string]bool{"SECRET": true}}
	err := UnmarkProtected(ps, "SECRET")
	require.NoError(t, err)
	assert.False(t, ps.Keys["SECRET"])
}

func TestUnmarkProtected_FailsIfNotProtected(t *testing.T) {
	ps := &ProtectedStore{Keys: make(map[string]bool)}
	err := UnmarkProtected(ps, "MISSING")
	assert.Error(t, err)
}

func TestIsProtected_ReturnsTrueWhenSet(t *testing.T) {
	ps := &ProtectedStore{Keys: map[string]bool{"API_KEY": true}}
	assert.True(t, IsProtected(ps, "API_KEY"))
	assert.False(t, IsProtected(ps, "OTHER"))
}

func TestSaveAndLoadProtected(t *testing.T) {
	dir := t.TempDir()
	storePath := filepath.Join(dir, "default.env.enc")

	ps := &ProtectedStore{Keys: map[string]bool{"DB_PASS": true, "TOKEN": true}}
	err := SaveProtected(storePath, ps)
	require.NoError(t, err)

	loaded, err := LoadProtected(storePath)
	require.NoError(t, err)
	assert.True(t, loaded.Keys["DB_PASS"])
	assert.True(t, loaded.Keys["TOKEN"])
}

func TestLoadProtected_ReturnsEmptyWhenMissing(t *testing.T) {
	dir := t.TempDir()
	ps, err := LoadProtected(filepath.Join(dir, "nope.env.enc"))
	require.NoError(t, err)
	assert.Empty(t, ps.Keys)
}

func TestProtectedKeys_ReturnsSorted(t *testing.T) {
	ps := &ProtectedStore{Keys: map[string]bool{"ZEBRA": true, "ALPHA": true, "MANGO": true}}
	keys := ProtectedKeys(ps)
	assert.Equal(t, []string{"ALPHA", "MANGO", "ZEBRA"}, keys)
}

func TestProtectedFilePath(t *testing.T) {
	path := protectedFilePath("/home/user/.envoy/default.env.enc")
	assert.Equal(t, "/home/user/.envoy/default.env.protected.json", path)
	_ = os.Getenv("HOME") // suppress unused import warning
}
