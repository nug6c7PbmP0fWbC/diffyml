package diff

import (
	"testing"
)

func TestDedupe_Empty(t *testing.T) {
	out := Dedupe(nil, DefaultDedupeOptions())
	if len(out) != 0 {
		t.Fatalf("expected empty slice, got %d elements", len(out))
	}
}

func TestDedupe_NoDuplicates(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: Added},
		{Path: "b", Type: Removed},
		{Path: "c", Type: Modified},
	}
	out := Dedupe(changes, DefaultDedupeOptions())
	if len(out) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(out))
	}
}

func TestDedupe_PreferFirst(t *testing.T) {
	changes := []Change{
		{Path: "x", Type: Added, NewValue: "first"},
		{Path: "x", Type: Modified, NewValue: "second"},
	}
	out := Dedupe(changes, DedupeOptions{PreferLast: false})
	if len(out) != 1 {
		t.Fatalf("expected 1 change, got %d", len(out))
	}
	if out[0].NewValue != "first" {
		t.Errorf("expected 'first', got %v", out[0].NewValue)
	}
}

func TestDedupe_PreferLast(t *testing.T) {
	changes := []Change{
		{Path: "x", Type: Added, NewValue: "first"},
		{Path: "x", Type: Modified, NewValue: "second"},
	}
	out := Dedupe(changes, DedupeOptions{PreferLast: true})
	if len(out) != 1 {
		t.Fatalf("expected 1 change, got %d", len(out))
	}
	if out[0].NewValue != "second" {
		t.Errorf("expected 'second', got %v", out[0].NewValue)
	}
}

func TestDedupe_PreservesOrder(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: Added},
		{Path: "b", Type: Added},
		{Path: "a", Type: Modified}, // duplicate
		{Path: "c", Type: Removed},
	}
	out := Dedupe(changes, DefaultDedupeOptions())
	if len(out) != 3 {
		t.Fatalf("expected 3 changes, got %d", len(out))
	}
	paths := []string{out[0].Path, out[1].Path, out[2].Path}
	expected := []string{"a", "b", "c"}
	for i, p := range expected {
		if paths[i] != p {
			t.Errorf("position %d: expected %q, got %q", i, p, paths[i])
		}
	}
}

func TestDedupe_MultipleDuplicates(t *testing.T) {
	changes := []Change{
		{Path: "k", Type: Added, NewValue: "v1"},
		{Path: "k", Type: Modified, NewValue: "v2"},
		{Path: "k", Type: Modified, NewValue: "v3"},
	}
	out := Dedupe(changes, DedupeOptions{PreferLast: true})
	if len(out) != 1 {
		t.Fatalf("expected 1 change, got %d", len(out))
	}
	if out[0].NewValue != "v3" {
		t.Errorf("expected 'v3', got %v", out[0].NewValue)
	}
}
