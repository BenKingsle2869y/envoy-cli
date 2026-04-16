package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMarkSecret_AddsKey(t *testing.T) {
	s := &SecretStore{Keys: map[string]bool{}}
	if err := MarkSecret(s, "API_KEY"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Keys["API_KEY"] {
		t.Error("expected API_KEY to be marked")
	}
}

func TestMarkSecret_RejectsDuplicate(t *testing.T) {
	s := &SecretStore{Keys: map[string]bool{"API_KEY": true}}
	if err := MarkSecret(s, "API_KEY"); err == nil {
		t.Error("expected error for duplicate")
	}
}

func TestUnmarkSecret_RemovesKey(t *testing.T) {
	s := &SecretStore{Keys: map[string]bool{"API_KEY": true}}
	if err := UnmarkSecret(s, "API_KEY"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Keys["API_KEY"] {
		t.Error("expected API_KEY to be removed")
	}
}

func TestUnmarkSecret_FailsWhenNotMarked(t *testing.T) {
	s := &SecretStore{Keys: map[string]bool{}}
	if err := UnmarkSecret(s, "MISSING"); err == nil {
		t.Error("expected error")
	}
}

func TestIsSecret(t *testing.T) {
	s := &SecretStore{Keys: map[string]bool{"TOKEN": true}}
	if !IsSecret(s, "TOKEN") {
		t.Error("expected TOKEN to be secret")
	}
	if IsSecret(s, "OTHER") {
		t.Error("expected OTHER to not be secret")
	}
}

func TestMaskIfSecret(t *testing.T) {
	s := &SecretStore{Keys: map[string]bool{"TOKEN": true}}
	if got := MaskIfSecret(s, "TOKEN", "abc123"); got != "********" {
		t.Errorf("expected masked, got %q", got)
	}
	if got := MaskIfSecret(s, "OTHER", "visible"); got != "visible" {
		t.Errorf("expected visible, got %q", got)
	}
}

func TestSaveAndLoadSecrets(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secrets.json")
	s := &SecretStore{Keys: map[string]bool{"A": true, "B": true}}
	if err := SaveSecrets(path, s); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadSecrets(path)
	if err != nil {
		t.Fatal(err)
	}
	if !loaded.Keys["A"] || !loaded.Keys["B"] {
		t.Error("expected A and B to be loaded")
	}
}

func TestLoadSecrets_EmptyWhenNotExist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	s, err := LoadSecrets(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Keys) != 0 {
		t.Error("expected empty store")
	}
	_ = os.Remove(path)
}
