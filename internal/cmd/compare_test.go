package cmd

import (
	"bytes"
	"testing"

	"envoy-cli/internal/store"
)

func setupCompareStores(t *testing.T, passphraseA, passphraseB string) (string, string) {
	t.Helper()
	tmpA := t.TempDir() + "/a.env"
	tmpB := t.TempDir() + "/b.env"

	stA := &store.Store{Entries: map[string]string{"FOO": "bar", "SHARED": "same", "DIFF": "aval"}}
	stB := &store.Store{Entries: map[string]string{"BAZ": "qux", "SHARED": "same", "DIFF": "bval"}}

	if err := store.Save(tmpA, stA, passphraseA); err != nil {
		t.Fatal(err)
	}
	if err := store.Save(tmpB, stB, passphraseB); err != nil {
		t.Fatal(err)
	}
	return tmpA, tmpB
}

func TestCompareStores_OnlyInA(t *testing.T) {
	a := map[string]string{"FOO": "1", "SHARED": "x"}
	b := map[string]string{"SHARED": "x"}
	r := compareStores(a, b)
	if len(r.OnlyInA) != 1 || r.OnlyInA[0] != "FOO" {
		t.Errorf("expected FOO only in A, got %v", r.OnlyInA)
	}
}

func TestCompareStores_OnlyInB(t *testing.T) {
	a := map[string]string{"SHARED": "x"}
	b := map[string]string{"BAR": "2", "SHARED": "x"}
	r := compareStores(a, b)
	if len(r.OnlyInB) != 1 || r.OnlyInB[0] != "BAR" {
		t.Errorf("expected BAR only in B, got %v", r.OnlyInB)
	}
}

func TestCompareStores_Different(t *testing.T) {
	a := map[string]string{"KEY": "aval"}
	b := map[string]string{"KEY": "bval"}
	r := compareStores(a, b)
	if len(r.Different) != 1 || r.Different[0] != "KEY" {
		t.Errorf("expected KEY to differ, got %v", r.Different)
	}
}

func TestCompareStores_Same(t *testing.T) {
	a := map[string]string{"KEY": "val"}
	b := map[string]string{"KEY": "val"}
	r := compareStores(a, b)
	if len(r.Same) != 1 || r.Same[0] != "KEY" {
		t.Errorf("expected KEY to be same, got %v", r.Same)
	}
	if len(r.Different)+len(r.OnlyInA)+len(r.OnlyInB) != 0 {
		t.Errorf("unexpected differences")
	}
}

func TestCompareStores_Empty(t *testing.T) {
	r := compareStores(map[string]string{}, map[string]string{})
	if len(r.OnlyInA)+len(r.OnlyInB)+len(r.Different)+len(r.Same) != 0 {
		t.Errorf("expected empty result")
	}
}

func TestSorted(t *testing.T) {
	in := []string{"z", "a", "m"}
	out := sorted(in)
	if out[0] != "a" || out[1] != "m" || out[2] != "z" {
		t.Errorf("unexpected sort order: %v", out)
	}
}

// Ensure compareStores is usable without stdout side-effects
func TestCompareStores_NoSideEffects(t *testing.T) {
	var buf bytes.Buffer
	_ = buf
	a := map[string]string{"A": "1"}
	b := map[string]string{"B": "2"}
	r := compareStores(a, b)
	if len(r.OnlyInA) != 1 || len(r.OnlyInB) != 1 {
		t.Errorf("unexpected result: %+v", r)
	}
}
