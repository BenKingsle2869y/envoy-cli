package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"envoy-cli/internal/store"
)

func setupKVStore(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.env.enc")
	passphrase := "test-passphrase"
	s := &store.Store{Env: map[string]string{}}
	if err := store.Save(path, passphrase, s); err != nil {
		t.Fatalf("failed to create test store: %v", err)
	}
	return path, passphrase
}

func TestSetCmd_AddsKey(t *testing.T) {
	path, pass := setupKVStore(t)
	os.Setenv("ENVOY_PASSPHRASE", pass)
	defer os.Unsetenv("ENVOY_PASSPHRASE")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"set", "FOO=bar", "--file", path})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("set command failed: %v", err)
	}

	s, err := store.Load(path, pass)
	if err != nil {
		t.Fatalf("failed to load store: %v", err)
	}
	if s.Env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", s.Env["FOO"])
	}
}

func TestGetCmd_ReturnsValue(t *testing.T) {
	path, pass := setupKVStore(t)
	os.Setenv("ENVOY_PASSPHRASE", pass)
	defer os.Unsetenv("ENVOY_PASSPHRASE")

	s := &store.Store{Env: map[string]string{"MY_KEY": "my_value"}}
	if err := store.Save(path, pass, s); err != nil {
		t.Fatalf("failed to save store: %v", err)
	}

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"get", "MY_KEY", "--file", path})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("get command failed: %v", err)
	}
}

func TestUnsetCmd_RemovesKey(t *testing.T) {
	path, pass := setupKVStore(t)
	os.Setenv("ENVOY_PASSPHRASE", pass)
	defer os.Unsetenv("ENVOY_PASSPHRASE")

	s := &store.Store{Env: map[string]string{"TO_REMOVE": "value"}}
	if err := store.Save(path, pass, s); err != nil {
		t.Fatalf("failed to save store: %v", err)
	}

	rootCmd.SetArgs([]string{"unset", "TO_REMOVE", "--file", path})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unset command failed: %v", err)
	}

	loaded, err := store.Load(path, pass)
	if err != nil {
		t.Fatalf("failed to load store: %v", err)
	}
	if _, ok := loaded.Env["TO_REMOVE"]; ok {
		t.Error("expected TO_REMOVE to be deleted")
	}
}

func TestListCmd_PrintsKeys(t *testing.T) {
	path, pass := setupKVStore(t)
	os.Setenv("ENVOY_PASSPHRASE", pass)
	defer os.Unsetenv("ENVOY_PASSPHRASE")

	s := &store.Store{Env: map[string]string{"ALPHA": "1", "BETA": "2"}}
	if err := store.Save(path, pass, s); err != nil {
		t.Fatalf("failed to save store: %v", err)
	}

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"list", "--file", path})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("list command failed: %v", err)
	}
}
