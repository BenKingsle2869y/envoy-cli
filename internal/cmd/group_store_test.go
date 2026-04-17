package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddToGroup_AddsKey(t *testing.T) {
	groups := map[string][]string{}
	if err := AddToGroup(groups, "backend", "DB_URL"); err != nil {
		t.Fatal(err)
	}
	if len(groups["backend"]) != 1 || groups["backend"][0] != "DB_URL" {
		t.Errorf("expected DB_URL in backend group")
	}
}

func TestAddToGroup_RejectsDuplicate(t *testing.T) {
	groups := map[string][]string{"backend": {"DB_URL"}}
	err := AddToGroup(groups, "backend", "DB_URL")
	if err == nil {
		t.Fatal("expected error for duplicate key")
	}
}

func TestRemoveFromGroup_RemovesKey(t *testing.T) {
	groups := map[string][]string{"backend": {"DB_URL", "API_KEY"}}
	if err := RemoveFromGroup(groups, "backend", "DB_URL"); err != nil {
		t.Fatal(err)
	}
	if len(groups["backend"]) != 1 {
		t.Errorf("expected 1 key remaining")
	}
}

func TestRemoveFromGroup_DeletesEmptyGroup(t *testing.T) {
	groups := map[string][]string{"backend": {"DB_URL"}}
	if err := RemoveFromGroup(groups, "backend", "DB_URL"); err != nil {
		t.Fatal(err)
	}
	if _, ok := groups["backend"]; ok {
		t.Errorf("expected group to be deleted when empty")
	}
}

func TestRemoveFromGroup_FailsWhenGroupMissing(t *testing.T) {
	groups := map[string][]string{}
	err := RemoveFromGroup(groups, "missing", "KEY")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSaveAndLoadGroups(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.groups.json")
	groups := map[string][]string{"infra": {"HOST", "PORT"}}
	if err := SaveGroups(path, groups); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadGroups(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded["infra"]) != 2 {
		t.Errorf("expected 2 keys in infra group")
	}
}

func TestLoadGroups_EmptyWhenNotExist(t *testing.T) {
	groups, err := LoadGroups("/nonexistent/path.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 0 {
		t.Errorf("expected empty map")
	}
	_ = os.Remove("/nonexistent/path.json")
}
