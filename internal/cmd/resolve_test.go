package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/store"
)

func TestExpandVars_BraceStyle(t *testing.T) {
	lookup := map[string]string{"HOST": "localhost", "PORT": "5432"}
	out := expandVars("postgres://${HOST}:${PORT}/db", lookup)
	if out != "postgres://localhost:5432/db" {
		t.Fatalf("unexpected: %s", out)
	}
}

func TestExpandVars_DollarStyle(t *testing.T) {
	lookup := map[string]string{"NAME": "world"}
	out := expandVars("hello $NAME", lookup)
	if out != "hello world" {
		t.Fatalf("unexpected: %s", out)
	}
}

func TestExpandVars_UnknownLeftAsIs(t *testing.T) {
	lookup := map[string]string{}
	out := expandVars("${MISSING}", lookup)
	if out != "${MISSING}" {
		t.Fatalf("should be unchanged, got: %s", out)
	}
}

func TestResolveReferences_DetectsChanges(t *testing.T) {
	entries := map[string]store.Entry{
		"BASE_URL": {Value: "http://localhost"},
		"API_URL":  {Value: "${BASE_URL}/api"},
		"PLAIN":    {Value: "no-ref"},
	}
	_, changed := resolveReferences(entries)
	if _, ok := changed["API_URL"]; !ok {
		t.Fatal("expected API_URL to be in changed map")
	}
	if _, ok := changed["PLAIN"]; ok {
		t.Fatal("PLAIN should not be in changed map")
	}
	if _, ok := changed["BASE_URL"]; ok {
		t.Fatal("BASE_URL should not be in changed map")
	}
}

func TestResolveReferences_NoRefs(t *testing.T) {
	entries := map[string]store.Entry{
		"FOO": {Value: "bar"},
		"BAZ": {Value: "qux"},
	}
	_, changed := resolveReferences(entries)
	if len(changed) != 0 {
		t.Fatalf("expected no changes, got %d", len(changed))
	}
}

func TestResolveReferences_MultipleRefs(t *testing.T) {
	entries := map[string]store.Entry{
		"SCHEME": {Value: "https"},
		"HOST":   {Value: "example.com"},
		"URL":    {Value: "${SCHEME}://${HOST}"},
	}
	_, changed := resolveReferences(entries)
	if v, ok := changed["URL"]; !ok || v != "https://example.com" {
		t.Fatalf("expected resolved URL, got %q", v)
	}
}

// Ensure unused import is satisfied
var _ = bytes.NewBuffer
