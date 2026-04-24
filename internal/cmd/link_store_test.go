package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddLink_AddsSuccessfully(t *testing.T) {
	links := make(Links)
	if err := AddLink(links, "DB_URL", "production"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if links["DB_URL"] != "production" {
		t.Errorf("expected production, got %s", links["DB_URL"])
	}
}

func TestAddLink_RejectsDuplicate(t *testing.T) {
	links := Links{"DB_URL": "production"}
	if err := AddLink(links, "DB_URL", "staging"); err == nil {
		t.Error("expected error for duplicate link")
	}
}

func TestAddLink_NilMapReturnsError(t *testing.T) {
	if err := AddLink(nil, "KEY", "ctx"); err == nil {
		t.Error("expected error for nil map")
	}
}

func TestRemoveLink_RemovesExisting(t *testing.T) {
	links := Links{"DB_URL": "production"}
	if err := RemoveLink(links, "DB_URL"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := links["DB_URL"]; ok {
		t.Error("expected key to be removed")
	}
}

func TestRemoveLink_FailsWhenNotFound(t *testing.T) {
	links := make(Links)
	if err := RemoveLink(links, "MISSING"); err == nil {
		t.Error("expected error for missing key")
	}
}

func TestResolveLink_ReturnsTarget(t *testing.T) {
	links := Links{"API_KEY": "production"}
	target, ok := ResolveLink(links, "API_KEY")
	if !ok || target != "production" {
		t.Errorf("expected production, got %s (ok=%v)", target, ok)
	}
}

func TestSaveAndLoadLinks(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.links.json")
	links := Links{"DB_URL": "production", "API_KEY": "staging"}
	if err := SaveLinks(path, links); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadLinks(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded["DB_URL"] != "production" {
		t.Errorf("expected production, got %s", loaded["DB_URL"])
	}
}

func TestLoadLinks_EmptyWhenNotExist(t *testing.T) {
	links, err := LoadLinks("/tmp/nonexistent-links-file.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(links) != 0 {
		t.Errorf("expected empty links, got %d", len(links))
	}
	_ = os.Remove("/tmp/nonexistent-links-file.json")
}
