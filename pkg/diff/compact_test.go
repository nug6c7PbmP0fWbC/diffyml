package diff

import (
	"testing"
)

func makeCompactChanges() []Change {
	return []Change{
		{Path: "a.b", Type: ChangeAdded, Value: "1"},
		{Path: "a.b", Type: ChangeModified, Before: "1", Value: "2"},
		{Path: "a.b", Type: ChangeModified, Before: "2", Value: "3"},
		{Path: "x.y", Type: ChangeRemoved, Before: "old"},
		{Path: "x.y", Type: ChangeAdded, Value: "new"},
	}
}

func TestCompact_Empty(t *testing.T) {
	out := Compact(nil, DefaultCompactOptions())
	if len(out) != 0 {
		t.Fatalf("expected empty, got %d", len(out))
	}
}

func TestCompact_KeepLast(t *testing.T) {
	opts := DefaultCompactOptions() // MaxPerPath=1, KeepLast=true
	out := Compact(makeCompactChanges(), opts)
	if len(out) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(out))
	}
	if out[0].Path != "a.b" || out[0].Value != "3" {
		t.Errorf("expected last a.b change (value=3), got %+v", out[0])
	}
	if out[1].Path != "x.y" || out[1].Type != ChangeAdded {
		t.Errorf("expected last x.y change (added), got %+v", out[1])
	}
}

func TestCompact_KeepFirst(t *testing.T) {
	opts := CompactOptions{MaxPerPath: 1, KeepLast: false}
	out := Compact(makeCompactChanges(), opts)
	if len(out) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(out))
	}
	if out[0].Path != "a.b" || out[0].Type != ChangeAdded {
		t.Errorf("expected first a.b change (added), got %+v", out[0])
	}
	if out[1].Path != "x.y" || out[1].Type != ChangeRemoved {
		t.Errorf("expected first x.y change (removed), got %+v", out[1])
	}
}

func TestCompact_MaxPerPathTwo(t *testing.T) {
	opts := CompactOptions{MaxPerPath: 2, KeepLast: true}
	out := Compact(makeCompactChanges(), opts)
	// a.b has 3 entries → keep last 2; x.y has 2 → keep both
	if len(out) != 4 {
		t.Fatalf("expected 4 changes, got %d", len(out))
	}
}

func TestCompact_NoLimit(t *testing.T) {
	opts := CompactOptions{MaxPerPath: 0}
	changes := makeCompactChanges()
	out := Compact(changes, opts)
	if len(out) != len(changes) {
		t.Fatalf("expected %d changes unchanged, got %d", len(changes), len(out))
	}
}

func TestCompact_PreservesPathOrder(t *testing.T) {
	changes := []Change{
		{Path: "z", Type: ChangeAdded, Value: "1"},
		{Path: "a", Type: ChangeAdded, Value: "2"},
		{Path: "z", Type: ChangeModified, Before: "1", Value: "2"},
	}
	opts := DefaultCompactOptions()
	out := Compact(changes, opts)
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
	if out[0].Path != "z" {
		t.Errorf("expected first path to be z, got %s", out[0].Path)
	}
	if out[1].Path != "a" {
		t.Errorf("expected second path to be a, got %s", out[1].Path)
	}
}
