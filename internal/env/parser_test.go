package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseString_BasicKeyValue(t *testing.T) {
	input := "FOO=bar\nBAZ=qux\n"
	envMap, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if envMap["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", envMap["FOO"])
	}
	if envMap["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", envMap["BAZ"])
	}
}

func TestParseString_SkipsCommentsAndBlanks(t *testing.T) {
	input := "# this is a comment\n\nKEY=value\n"
	envMap, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(envMap) != 1 {
		t.Errorf("expected 1 entry, got %d", len(envMap))
	}
	if envMap["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", envMap["KEY"])
	}
}

func TestParseString_QuotedValues(t *testing.T) {
	input := `SINGLE='hello world'
DOUBLE="goodbye world"
`
	envMap, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if envMap["SINGLE"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", envMap["SINGLE"])
	}
	if envMap["DOUBLE"] != "goodbye world" {
		t.Errorf("expected 'goodbye world', got %q", envMap["DOUBLE"])
	}
}

func TestParseString_InvalidLine(t *testing.T) {
	_, err := ParseString("INVALID_LINE_NO_EQUALS")
	if err == nil {
		t.Error("expected error for invalid line, got nil")
	}
}

func TestParseFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := "APP_ENV=production\nDEBUG=false\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	envMap, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if envMap["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", envMap["APP_ENV"])
	}
}

func TestSerialize(t *testing.T) {
	envMap := EnvMap{"KEY": "simple"}
	out := Serialize(envMap)
	if out != "KEY=simple\n" {
		t.Errorf("unexpected serialization: %q", out)
	}

	envMap2 := EnvMap{"MSG": "hello world"}
	out2 := Serialize(envMap2)
	if out2 != `MSG="hello world"`+"\n" {
		t.Errorf("unexpected serialization for spaced value: %q", out2)
	}
}
