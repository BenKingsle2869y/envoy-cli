package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func setupMaskPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "default.masks.json")
}

func TestMarkMasked_AddsKey(t *testing.T) {
	masks := map[string]bool{}
	if err := MarkMasked(masks, "SECRET_KEY"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !masks["SECRET_KEY"] {
		t.Error("expected SECRET_KEY to be masked")
	}
}

func TestMarkMasked_RejectsDuplicate(t *testing.T) {
	masks := map[string]bool{"SECRET_KEY": true}
	if err := MarkMasked(masks, "SECRET_KEY"); err == nil {
		t.Error("expected error for duplicate mask")
	}
}

func TestMarkMasked_NilMapReturnsError(t *testing.T) {
	if err := MarkMasked(nil, "KEY"); err == nil {
		t.Error("expected error for nil map")
	}
}

func TestUnmarkMasked_RemovesKey(t *testing.T) {
	masks := map[string]bool{"API_TOKEN": true}
	if err := UnmarkMasked(masks, "API_TOKEN"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if masks["API_TOKEN"] {
		t.Error("expected API_TOKEN to be removed")
	}
}

func TestUnmarkMasked_FailsWhenNotMasked(t *testing.T) {
	masks := map[string]bool{}
	if err := UnmarkMasked(masks, "MISSING"); err == nil {
		t.Error("expected error when key is not masked")
	}
}

func TestIsMasked(t *testing.T) {
	masks := map[string]bool{"DB_PASS": true}
	if !IsMasked(masks, "DB_PASS") {
		t.Error("expected DB_PASS to be masked")
	}
	if IsMasked(masks, "OTHER") {
		t.Error("expected OTHER not to be masked")
	}
}

func TestSaveAndLoadMasks(t *testing.T) {
	path := setupMaskPath(t)
	masks := map[string]bool{"TOKEN": true, "SECRET": true}
	if err := SaveMasks(path, masks); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadMasks(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if !loaded["TOKEN"] || !loaded["SECRET"] {
		t.Error("expected both keys to be loaded")
	}
}

func TestLoadMasks_EmptyWhenNotExist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.masks.json")
	masks, err := LoadMasks(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(masks) != 0 {
		t.Error("expected empty masks")
	}
}

func TestMaskedKeys_ReturnsSorted(t *testing.T) {
	masks := map[string]bool{"ZEBRA": true, "ALPHA": true, "MANGO": true}
	keys := MaskedKeys(masks)
	if keys[0] != "ALPHA" || keys[1] != "MANGO" || keys[2] != "ZEBRA" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestLoadMasks_ReturnsErrorOnCorruptFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.masks.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o600)
	if _, err := LoadMasks(path); err == nil {
		t.Error("expected parse error")
	}
}
