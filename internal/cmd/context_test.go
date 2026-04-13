package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func withTempHome(t *testing.T) (string, func()) {
	t.Helper()
	tmp := t.TempDir()
	orig := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	return tmp, func() { os.Setenv("HOME", orig) }
}

func TestActiveContext_DefaultsWhenNoFile(t *testing.T) {
	_, cleanup := withTempHome(t)
	defer cleanup()

	ctx, err := ActiveContext()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx != defaultContext {
		t.Errorf("expected %q, got %q", defaultContext, ctx)
	}
}

func TestSetAndActiveContext(t *testing.T) {
	_, cleanup := withTempHome(t)
	defer cleanup()

	if err := SetActiveContext("staging"); err != nil {
		t.Fatalf("SetActiveContext: %v", err)
	}
	ctx, err := ActiveContext()
	if err != nil {
		t.Fatalf("ActiveContext: %v", err)
	}
	if ctx != "staging" {
		t.Errorf("expected %q, got %q", "staging", ctx)
	}
}

func TestSetActiveContext_EmptyNameFails(t *testing.T) {
	_, cleanup := withTempHome(t)
	defer cleanup()

	if err := SetActiveContext(""); err == nil {
		t.Error("expected error for empty context name")
	}
}

func TestListContexts_ReturnsStoreFiles(t *testing.T) {
	tmp, cleanup := withTempHome(t)
	defer cleanup()

	dir := filepath.Join(tmp, ".envoy")
	os.MkdirAll(dir, 0700)
	for _, name := range []string{"development.enc", "staging.enc", "context"} {
		os.WriteFile(filepath.Join(dir, name), []byte("x"), 0600)
	}

	ctxs, err := ListContexts()
	if err != nil {
		t.Fatalf("ListContexts: %v", err)
	}
	if len(ctxs) != 2 {
		t.Errorf("expected 2 contexts, got %d: %v", len(ctxs), ctxs)
	}
}
