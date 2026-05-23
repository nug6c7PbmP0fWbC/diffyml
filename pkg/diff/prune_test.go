package diff

import (
	"testing"
)

func TestPrune_Empty(t *testing.T) {
	result := Prune(nil, DefaultPruneOptions())
	if len(result) != 0 {
		t.Fatalf("expected empty, got %d", len(result))
	}
}

func TestPrune_RemoveNilValues(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: Removed, Before: "x", After: nil},
		{Path: "b", Type: Modified, Before: "x", After: "y"},
	}
	opts := DefaultPruneOptions()
	opts.RemoveNilValues = true
	opts.RemoveUnchangedKeys = false
	result := Prune(changes, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Path != "b" {
		t.Errorf("unexpected path %q", result[0].Path)
	}
}

func TestPrune_RemoveZeroValues(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: Modified, Before: "x", After: ""},
		{Path: "b", Type: Modified, Before: 1, After: 0},
		{Path: "c", Type: Modified, Before: false, After: true},
	}
	opts := PruneOptions{RemoveZeroValues: true}
	result := Prune(changes, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Path != "c" {
		t.Errorf("expected path c, got %q", result[0].Path)
	}
}

func TestPrune_RemoveUnchangedKeys(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: Modified, Before: "same", After: "same"},
		{Path: "b", Type: Modified, Before: "old", After: "new"},
	}
	opts := PruneOptions{RemoveUnchangedKeys: true}
	result := Prune(changes, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Path != "b" {
		t.Errorf("expected path b, got %q", result[0].Path)
	}
}

func TestPrune_PathFilter_LimitsScope(t *testing.T) {
	changes := []Change{
		{Path: "secrets.key", Type: Removed, Before: "v", After: nil},
		{Path: "config.key", Type: Removed, Before: "v", After: nil},
	}
	opts := PruneOptions{
		RemoveNilValues: true,
		Paths:           []string{"secrets"},
	}
	result := Prune(changes, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Path != "config.key" {
		t.Errorf("expected config.key, got %q", result[0].Path)
	}
}

func TestPrune_DoesNotMutateInput(t *testing.T) {
	original := []Change{
		{Path: "a", Type: Modified, Before: "x", After: "x"},
	}
	opts := PruneOptions{RemoveUnchangedKeys: true}
	_ = Prune(original, opts)
	if len(original) != 1 {
		t.Error("original slice was mutated")
	}
}

func TestIsZero(t *testing.T) {
	cases := []struct {
		val      interface{}
		expected bool
	}{
		{nil, true},
		{"", true},
		{0, true},
		{float64(0), true},
		{false, true},
		{"hello", false},
		{1, false},
		{true, false},
	}
	for _, tc := range cases {
		if got := isZero(tc.val); got != tc.expected {
			t.Errorf("isZero(%v) = %v, want %v", tc.val, got, tc.expected)
		}
	}
}
