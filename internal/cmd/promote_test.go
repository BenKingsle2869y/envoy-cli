package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/store"
)

func setupPromoteStores(t *testing.T, srcVars, dstVars map[string]string) (srcPath, dstPath, passphrase string) {
	t.Helper()
	passphrase = "test-pass"

	srcPath = t.TempDir() + "/src.env"
	dstPath = t.TempDir() + "/dst.env"

	src := &store.Store{Vars: srcVars}
	dst := &store.Store{Vars: dstVars}

	if err := store.Save(srcPath, passphrase, src); err != nil {
		t.Fatal(err)
	}
	if err := store.Save(dstPath, passphrase, dst); err != nil {
		t.Fatal(err)
	}
	return
}

func execPromoteCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"promote"}, args...))
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestPromoteCmd_SkipsExistingByDefault(t *testing.T) {
	saveContext(t, "src")
	src := StorePathForContext("src")
	dst := StorePathForContext("dst")
	pass := "pass"

	_ = store.Save(src, pass, &store.Store{Vars: map[string]string{"A": "1", "B": "2"}})
	_ = store.Save(dst, pass, &store.Store{Vars: map[string]string{"A": "existing"}})

	_, err := execPromoteCmd(t, "src", "dst", "--passphrase", pass)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, _ := store.Load(dst, pass)
	if loaded.Vars["A"] != "existing" {
		t.Errorf("expected A to remain 'existing', got %q", loaded.Vars["A"])
	}
	if loaded.Vars["B"] != "2" {
		t.Errorf("expected B=2, got %q", loaded.Vars["B"])
	}
}

func TestPromoteCmd_OverwriteFlag(t *testing.T) {
	pass := "pass"
	src := StorePathForContext("src2")
	dst := StorePathForContext("dst2")

	_ = store.Save(src, pass, &store.Store{Vars: map[string]string{"A": "new"}})
	_ = store.Save(dst, pass, &store.Store{Vars: map[string]string{"A": "old"}})

	_, err := execPromoteCmd(t, "src2", "dst2", "--passphrase", pass, "--overwrite")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, _ := store.Load(dst, pass)
	if loaded.Vars["A"] != "new" {
		t.Errorf("expected A=new after overwrite, got %q", loaded.Vars["A"])
	}
}

func TestPromoteCmd_SameContextFails(t *testing.T) {
	_, err := execPromoteCmd(t, "staging", "staging", "--passphrase", "x")
	if err == nil {
		t.Error("expected error for same context, got nil")
	}
}
