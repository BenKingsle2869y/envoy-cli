package env

import (
	"testing"
)

func TestDiff_AddedKeys(t *testing.T) {
	base := map[string]string{"A": "1"}
	target := map[string]string{"A": "1", "B": "2"}

	result := Diff(base, target)

	if len(result.Added) != 1 || result.Added["B"] != "2" {
		t.Errorf("expected B=2 in Added, got %v", result.Added)
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected no removals, got %v", result.Removed)
	}
}

func TestDiff_RemovedKeys(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	target := map[string]string{"A": "1"}

	result := Diff(base, target)

	if len(result.Removed) != 1 || result.Removed["B"] != "2" {
		t.Errorf("expected B=2 in Removed, got %v", result.Removed)
	}
}

func TestDiff_ChangedKeys(t *testing.T) {
	base := map[string]string{"A": "old"}
	target := map[string]string{"A": "new"}

	result := Diff(base, target)

	pair, ok := result.Changed["A"]
	if !ok {
		t.Fatal("expected A in Changed")
	}
	if pair[0] != "old" || pair[1] != "new" {
		t.Errorf("expected [old new], got %v", pair)
	}
}

func TestDiff_UnchangedKeys(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	target := map[string]string{"A": "1", "B": "2"}

	result := Diff(base, target)

	if result.HasChanges() {
		t.Error("expected no changes")
	}
	if len(result.Unchanged) != 2 {
		t.Errorf("expected 2 unchanged, got %d", len(result.Unchanged))
	}
}

func TestApply_MergesTargetIntoBase(t *testing.T) {
	base := map[string]string{"A": "1", "B": "old"}
	target := map[string]string{"B": "new", "C": "3"}

	result := Apply(base, target)

	if result["A"] != "1" {
		t.Errorf("expected A=1, got %s", result["A"])
	}
	if result["B"] != "new" {
		t.Errorf("expected B=new, got %s", result["B"])
	}
	if result["C"] != "3" {
		t.Errorf("expected C=3, got %s", result["C"])
	}
}

func TestApply_DoesNotMutateBase(t *testing.T) {
	base := map[string]string{"A": "1"}
	target := map[string]string{"A": "2"}

	_ = Apply(base, target)

	if base["A"] != "1" {
		t.Error("Apply should not mutate the base map")
	}
}
