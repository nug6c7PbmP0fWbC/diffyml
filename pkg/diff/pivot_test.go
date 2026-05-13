package diff

import (
	"testing"
)

func TestPivot_Empty(t *testing.T) {
	result := Pivot(nil, DefaultPivotOptions())
	if len(result.Keys) != 0 {
		t.Fatalf("expected 0 keys, got %d", len(result.Keys))
	}
	if len(result.Index) != 0 {
		t.Fatalf("expected empty index")
	}
}

func TestPivot_UniqueKeys(t *testing.T) {
	changes := []Change{
		{Path: "a.b", Type: ChangeAdded, Value: "1"},
		{Path: "a.c", Type: ChangeRemoved, Before: "2"},
		{Path: "x.y", Type: ChangeModified, Before: "old", Value: "new"},
	}

	result := Pivot(changes, DefaultPivotOptions())

	if len(result.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(result.Keys))
	}
	if !result.Has("a.b") || !result.Has("a.c") || !result.Has("x.y") {
		t.Fatal("missing expected key")
	}
}

func TestPivot_DuplicateKeys(t *testing.T) {
	changes := []Change{
		{Path: "dup", Type: ChangeAdded, Value: "first"},
		{Path: "dup", Type: ChangeModified, Before: "a", Value: "b"},
	}

	result := Pivot(changes, DefaultPivotOptions())

	if len(result.Keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(result.Keys))
	}
	got := result.Lookup("dup")
	if len(got) != 2 {
		t.Fatalf("expected 2 changes under 'dup', got %d", len(got))
	}
}

func TestPivot_PreservesOrder(t *testing.T) {
	changes := []Change{
		{Path: "z", Type: ChangeAdded},
		{Path: "a", Type: ChangeAdded},
		{Path: "m", Type: ChangeAdded},
	}

	result := Pivot(changes, DefaultPivotOptions())

	expected := []string{"z", "a", "m"}
	for i, k := range result.Keys {
		if k != expected[i] {
			t.Errorf("key[%d]: want %q, got %q", i, expected[i], k)
		}
	}
}

func TestPivot_CustomKeyFunc(t *testing.T) {
	changes := []Change{
		{Path: "a.b", Type: ChangeAdded},
		{Path: "a.c", Type: ChangeAdded},
		{Path: "b.d", Type: ChangeRemoved},
	}

	// Group by top-level prefix
	opts := PivotOptions{
		KeyFunc: func(c Change) string {
			if len(c.Path) >= 1 {
				return string(c.Path[0])
			}
			return c.Path
		},
	}

	result := Pivot(changes, opts)

	if len(result.Keys) != 2 {
		t.Fatalf("expected 2 keys (a, b), got %d", len(result.Keys))
	}
	if len(result.Lookup("a")) != 2 {
		t.Errorf("expected 2 changes under 'a'")
	}
	if len(result.Lookup("b")) != 1 {
		t.Errorf("expected 1 change under 'b'")
	}
}

func TestPivot_Has_Missing(t *testing.T) {
	result := Pivot([]Change{{Path: "x", Type: ChangeAdded}}, DefaultPivotOptions())
	if result.Has("nonexistent") {
		t.Error("Has should return false for missing key")
	}
}
