package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/store"
)

func setupTouchStore(t *testing.T) (string, string, string) {
	t.Helper()
	h := t.TempDir()
	t.Setenv("HOME", h)
	t.Setenv("ENVOY_PASSPHRASE", "touchpass")

	ctx := "default"
	if err := SetActiveContext(ctx); err != nil {
		t.Fatalf("SetActiveContext: %v", err)
	}

	path := StorePathForContext(ctx)
	s := &store.Store{
		Vars: map[string]string{"FOO": "bar", "BAZ": "qux"},
	}
	if err := store.Save(path, "touchpass", s); err != nil {
		t.Fatalf("Save: %v", err)
	}
	return h, ctx, path
}

func execTouchCmd(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(append([]string{"touch"}, args...))
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestTouchCmd_Success(t *testing.T) {
	_, _, _ = setupTouchStore(t)

	out, err := execTouchCmd("FOO")
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	if !contains(out, "touched") || !contains(out, "FOO") {
		t.Errorf("expected touch confirmation, got: %s", out)
	}
}

func TestTouchCmd_MissingKey(t *testing.T) {
	_, _, _ = setupTouchStore(t)

	_, err := execTouchCmd("NONEXISTENT")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestTouchCmd_PreservesValue(t *testing.T) {
	_, _, path := setupTouchStore(t)

	_, err := execTouchCmd("BAZ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s, err := store.Load(path, "touchpass")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if s.Vars["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux after touch, got %q", s.Vars["BAZ"])
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsHelper(s, sub))
}

func containsHelper(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
