package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func setupSchemaPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "schema_test.json")
}

func TestLoadSchema_EmptyWhenNotExist(t *testing.T) {
	entries, err := LoadSchema("/nonexistent/path/schema.json")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %d entries", len(entries))
	}
}

func TestSaveAndLoadSchema(t *testing.T) {
	path := setupSchemaPath(t)
	input := []SchemaEntry{
		{Key: "DB_URL", Required: true, Description: "Database connection URL"},
		{Key: "PORT", Required: false, Default: "8080"},
	}
	if err := SaveSchema(path, input); err != nil {
		t.Fatalf("SaveSchema failed: %v", err)
	}
	loaded, err := LoadSchema(path)
	if err != nil {
		t.Fatalf("LoadSchema failed: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(loaded))
	}
	if loaded[0].Key != "DB_URL" || !loaded[0].Required {
		t.Errorf("unexpected first entry: %+v", loaded[0])
	}
}

func TestValidateAgainstSchema_MissingRequired(t *testing.T) {
	schema := []SchemaEntry{
		{Key: "DB_URL", Required: true},
		{Key: "PORT", Required: false},
	}
	env := map[string]string{"PORT": "8080"}
	missing := ValidateAgainstSchema(schema, env)
	if len(missing) != 1 || missing[0] != "DB_URL" {
		t.Errorf("expected [DB_URL] missing, got %v", missing)
	}
}

func TestValidateAgainstSchema_AllPresent(t *testing.T) {
	schema := []SchemaEntry{
		{Key: "DB_URL", Required: true},
	}
	env := map[string]string{"DB_URL": "postgres://localhost/db"}
	missing := ValidateAgainstSchema(schema, env)
	if len(missing) != 0 {
		t.Errorf("expected no missing keys, got %v", missing)
	}
}

func TestSchemaAddCmd_AddsEntry(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	os.MkdirAll(filepath.Join(home, ".envoy"), 0o700)

	root := newRootCmdForTest()
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"schema", "add", "API_KEY", "--required", "--desc", "API access key"})
	if err := root.Execute(); err != nil {
		t.Fatalf("schema add failed: %v", err)
	}

	path := filepath.Join(home, ".envoy", "schema_default.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("schema file not created: %v", err)
	}
	var entries []SchemaEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(entries) != 1 || entries[0].Key != "API_KEY" || !entries[0].Required {
		t.Errorf("unexpected entries: %+v", entries)
	}
}
