package diff

import (
	"testing"
)

func TestNormalize_Empty(t *testing.T) {
	out := Normalize(nil, DefaultNormalizeOptions())
	if len(out) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(out))
	}
}

func TestNormalize_TrimSpace(t *testing.T) {
	changes := []Change{
		{Path: "app.name", Type: ChangeModified, OldValue: "  foo  ", NewValue: "  bar  "},
	}
	opts := DefaultNormalizeOptions()
	opts.TrimSpace = true
	out := Normalize(changes, opts)
	if len(out) != 1 {
		t.Fatalf("expected 1 change, got %d", len(out))
	}
	if out[0].OldValue != "foo" {
		t.Errorf("expected OldValue 'foo', got %q", out[0].OldValue)
	}
	if out[0].NewValue != "bar" {
		t.Errorf("expected NewValue 'bar', got %q", out[0].NewValue)
	}
}

func TestNormalize_NoTrimSpace(t *testing.T) {
	changes := []Change{
		{Path: "app.name", Type: ChangeModified, OldValue: "  foo  ", NewValue: "  bar  "},
	}
	opts := DefaultNormalizeOptions()
	opts.TrimSpace = false
	opts.IgnoreEmptyValues = false
	out := Normalize(changes, opts)
	if out[0].OldValue != "  foo  " {
		t.Errorf("expected untrimmed value, got %q", out[0].OldValue)
	}
}

func TestNormalize_LowercaseKeys(t *testing.T) {
	changes := []Change{
		{Path: "App.Name", Type: ChangeAdded, NewValue: "v1"},
	}
	opts := NormalizeOptions{LowercaseKeys: true}
	out := Normalize(changes, opts)
	if out[0].Path != "app.name" {
		t.Errorf("expected 'app.name', got %q", out[0].Path)
	}
}

func TestNormalize_IgnoreEmptyValues(t *testing.T) {
	changes := []Change{
		{Path: "app.desc", Type: ChangeModified, OldValue: "", NewValue: ""},
		{Path: "app.name", Type: ChangeModified, OldValue: "a", NewValue: "b"},
	}
	opts := DefaultNormalizeOptions()
	out := Normalize(changes, opts)
	if len(out) != 1 {
		t.Fatalf("expected 1 change after filtering empty, got %d", len(out))
	}
	if out[0].Path != "app.name" {
		t.Errorf("expected 'app.name', got %q", out[0].Path)
	}
}

func TestNormalize_NonStringValuesUnaffected(t *testing.T) {
	changes := []Change{
		{Path: "app.port", Type: ChangeModified, OldValue: 8080, NewValue: 9090},
	}
	opts := DefaultNormalizeOptions()
	out := Normalize(changes, opts)
	if len(out) != 1 {
		t.Fatalf("expected 1 change, got %d", len(out))
	}
	if out[0].OldValue != 8080 {
		t.Errorf("expected OldValue 8080, got %v", out[0].OldValue)
	}
}
