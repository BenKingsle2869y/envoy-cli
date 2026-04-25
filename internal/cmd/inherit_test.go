package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"envoy-cli/internal/store"
)

func setupInheritStores(t *testing.T) (string, func()) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	os.MkdirAll(filepath.Join(dir, ".envoy"), 0700)
	return dir, func() {}
}

func execInheritCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestInheritCmd_InheritsNewKeys(t *testing.T) {
	setupInheritStores(t)
	pass := "testpass"

	parentPath := StorePathForContext("production")
	parentStore, _ := store.Load(parentPath, pass)
	parentStore.Set("DB_HOST", "prod.db")
	parentStore.Set("API_KEY", "secret")
	store.Save(parentPath, pass, parentStore)

	activePath := StorePathForContext("staging")
	activeStore, _ := store.Load(activePath, pass)
	activeStore.Set("DB_HOST", "staging.db")
	store.Save(activePath, pass, activeStore)

	SetActiveContext("staging")

	out, err := execInheritCmd(t, "inherit", "run", "production", "--passphrase", pass)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Contains([]byte(out), []byte("1 key(s)")) {
		t.Errorf("expected inherited count in output, got: %s", out)
	}

	updated, _ := store.Load(activePath, pass)
	v, ok := updated.Get("API_KEY")
	if !ok || v != "secret" {
		t.Errorf("expected API_KEY to be inherited, got %q ok=%v", v, ok)
	}
	v2, _ := updated.Get("DB_HOST")
	if v2 != "staging.db" {
		t.Errorf("expected DB_HOST to remain unchanged, got %q", v2)
	}
}

func TestInheritCmd_OverwriteFlag(t *testing.T) {
	setupInheritStores(t)
	pass := "testpass"

	parentPath := StorePathForContext("production")
	parentStore, _ := store.Load(parentPath, pass)
	parentStore.Set("DB_HOST", "prod.db")
	store.Save(parentPath, pass, parentStore)

	activePath := StorePathForContext("staging")
	activeStore, _ := store.Load(activePath, pass)
	activeStore.Set("DB_HOST", "staging.db")
	store.Save(activePath, pass, activeStore)

	SetActiveContext("staging")

	_, err := execInheritCmd(t, "inherit", "run", "production", "--passphrase", pass, "--overwrite")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, _ := store.Load(activePath, pass)
	v, _ := updated.Get("DB_HOST")
	if v != "prod.db" {
		t.Errorf("expected DB_HOST to be overwritten to prod.db, got %q", v)
	}
}

func TestInheritCmd_SameContextFails(t *testing.T) {
	setupInheritStores(t)
	SetActiveContext("staging")

	_, err := execInheritCmd(t, "inherit", "run", "staging", "--passphrase", "pass")
	if err == nil {
		t.Error("expected error when parent and active context are the same")
	}
}

func TestInheritCmd_FailsWithoutPassphrase(t *testing.T) {
	setupInheritStores(t)
	os.Unsetenv("ENVOY_PASSPHRASE")
	SetActiveContext("staging")

	_, err := execInheritCmd(t, "inherit", "run", "production")
	if err == nil {
		t.Error("expected error without passphrase")
	}
}
