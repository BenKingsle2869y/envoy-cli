package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func setupAliasPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "default.aliases.json")
}

func TestAddAlias_AddsSuccessfully(t *testing.T) {
	aliases := map[string]string{}
	if err := AddAlias(aliases, "db", "DATABASE_URL"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if aliases["db"] != "DATABASE_URL" {
		t.Errorf("expected DATABASE_URL, got %s", aliases["db"])
	}
}

func TestAddAlias_RejectsDuplicate(t *testing.T) {
	aliases := map[string]string{"db": "DATABASE_URL"}
	if err := AddAlias(aliases, "db", "OTHER"); err == nil {
		t.Fatal("expected error for duplicate alias")
	}
}

func TestRemoveAlias_RemovesExisting(t *testing.T) {
	aliases := map[string]string{"db": "DATABASE_URL"}
	if err := RemoveAlias(aliases, "db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := aliases["db"]; ok {
		t.Error("alias should have been removed")
	}
}

func TestRemoveAlias_FailsWhenNotFound(t *testing.T) {
	aliases := map[string]string{}
	if err := RemoveAlias(aliases, "missing"); err == nil {
		t.Fatal("expected error for missing alias")
	}
}

func TestResolveAlias_ReturnsMappedKey(t *testing.T) {
	aliases := map[string]string{"db": "DATABASE_URL"}
	if got := ResolveAlias(aliases, "db"); got != "DATABASE_URL" {
		t.Errorf("expected DATABASE_URL, got %s", got)
	}
}

func TestResolveAlias_PassthroughWhenUnknown(t *testing.T) {
	aliases := map[string]string{}
	if got := ResolveAlias(aliases, "KEY"); got != "KEY" {
		t.Errorf("expected KEY, got %s", got)
	}
}

func TestSaveAndLoadAliases(t *testing.T) {
	path := setupAliasPath(t)
	aliases := map[string]string{"db": "DATABASE_URL", "port": "APP_PORT"}
	if err := SaveAliases(path, aliases); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadAliases(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded["db"] != "DATABASE_URL" || loaded["port"] != "APP_PORT" {
		t.Errorf("loaded aliases mismatch: %v", loaded)
	}
}

func TestLoadAliases_EmptyWhenNotExist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.aliases.json")
	aliases, err := LoadAliases(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(aliases) != 0 {
		t.Errorf("expected empty map, got %v", aliases)
	}
}

func execAliasCmd(t *testing.T, home string, args ...string) (string, error) {
	t.Helper()
	t.Setenv("HOME", home)
	var buf bytes.Buffer
	aliasCmd.SetOut(&buf)
	aliasCmd.SetErr(&buf)
	aliasCmd.SetArgs(args)
	err := aliasCmd.Execute()
	return buf.String(), err
}

func TestAliasListCmd_Empty(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if err := os.MkdirAll(filepath.Join(home, ".envoy"), 0o700); err != nil {
		t.Fatal(err)
	}
	out, err := execAliasCmd(t, home, "list")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Log("got empty output (no aliases defined message)")
	}
}
