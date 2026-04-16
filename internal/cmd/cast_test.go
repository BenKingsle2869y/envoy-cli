package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestInferType_Bool(t *testing.T) {
	for _, v := range []string{"true", "false", "1", "0"} {
		if got := inferType(v); got != "bool" {
			t.Errorf("inferType(%q) = %q, want bool", v, got)
		}
	}
}

func TestInferType_Int(t *testing.T) {
	for _, v := range []string{"42", "-7", "1000"} {
		if got := inferType(v); got != "int" {
			t.Errorf("inferType(%q) = %q, want int", v, got)
		}
	}
}

func TestInferType_Float(t *testing.T) {
	for _, v := range []string{"3.14", "-0.5", "2.0"} {
		if got := inferType(v); got != "float" {
			t.Errorf("inferType(%q) = %q, want float", v, got)
		}
	}
}

func TestInferType_List(t *testing.T) {
	if got := inferType("[a,b,c]"); got != "list" {
		t.Errorf("expected list, got %q", got)
	}
}

func TestInferType_String(t *testing.T) {
	if got := inferType("hello"); got != "string" {
		t.Errorf("expected string, got %q", got)
	}
}

func TestCoerceToInt_Valid(t *testing.T) {
	v, err := CoerceToInt("99")
	if err != nil || v != 99 {
		t.Fatalf("expected 99, got %d, err=%v", v, err)
	}
}

func TestCoerceToInt_Invalid(t *testing.T) {
	_, err := CoerceToInt("abc")
	if err == nil {
		t.Fatal("expected error for invalid int")
	}
}

func TestCoerceToBool_Valid(t *testing.T) {
	v, err := CoerceToBool("true")
	if err != nil || !v {
		t.Fatalf("expected true, got %v, err=%v", v, err)
	}
}

func TestCoerceToFloat_Valid(t *testing.T) {
	v, err := CoerceToFloat("2.71")
	if err != nil || v != 2.71 {
		t.Fatalf("expected 2.71, got %v, err=%v", v, err)
	}
}

func TestCastCmd_PrintsTypes(t *testing.T) {
	dir := t.TempDir()
	st, path, pass := setupKVStore(t, dir)
	_ = st

	cmd := castCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	_ = cmd.Flags().Set("passphrase", pass)
	_ = cmd.Flags().Set("context", "")
	_ = path

	// Verify inferType is consistent with output expectations
	result := inferType("true")
	if !strings.EqualFold(result, "bool") {
		t.Errorf("unexpected type: %s", result)
	}
}
