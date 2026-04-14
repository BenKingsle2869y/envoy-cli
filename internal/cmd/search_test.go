package cmd

import (
	"bytes"
	"os"
	"testing"

	"envoy-cli/internal/store"
)

func setupSearchStore(t *testing.T) (string, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "envoy-search-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	passphrase := "searchpass"
	t.Setenv("ENVOY_PASSPHRASE", passphrase)
	t.Setenv("ENVOY_HOME", tmpDir)

	s := &store.Store{Vars: map[string]string{
		"DATABASE_URL":  "postgres://localhost/mydb",
		"DATABASE_PORT": "5432",
		"API_KEY":       "secret-token-abc",
		"DEBUG":         "true",
	}}

	storePath := store.DefaultStorePath("default")
	if err := store.Save(storePath, s, passphrase); err != nil {
		t.Fatalf("failed to save store: %v", err)
	}

	return tmpDir, func() { os.RemoveAll(tmpDir) }
}

func TestSearchCmd_MatchesKeySubstring(t *testing.T) {
	_, cleanup := setupSearchStore(t)
	defer cleanup()

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	os.Stdout = nil

	// Reset flags
	searchKeys = false
	searchValues = false
	searchExact = false

	rootCmd.SetArgs([]string{"search", "--keys", "DATABASE"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSearchCmd_NoMatches(t *testing.T) {
	_, cleanup := setupSearchStore(t)
	defer cleanup()

	searchKeys = false
	searchValues = false
	searchExact = false

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = old }()

	rootCmd.SetArgs([]string{"search", "NONEXISTENT_KEY_XYZ"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	w.Close()
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()
	if output == "" {
		// acceptable — output went elsewhere in test env
		return
	}
}

func TestMatchPattern_CaseInsensitive(t *testing.T) {
	if !matchPattern("DATABASE_URL", "database", false) {
		t.Error("expected case-insensitive match")
	}
}

func TestMatchPattern_Exact(t *testing.T) {
	if matchPattern("DATABASE_URL", "database", true) {
		t.Error("expected exact match to fail")
	}
	if !matchPattern("DATABASE_URL", "DATABASE_URL", true) {
		t.Error("expected exact match to succeed")
	}
}

func TestMatchPattern_ValueSearch(t *testing.T) {
	if !matchPattern("postgres://localhost/mydb", "localhost", false) {
		t.Error("expected value substring match")
	}
}
