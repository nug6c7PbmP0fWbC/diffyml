package diff

import (
	"testing"
)

func TestBuildIndex_Empty(t *testing.T) {
	idx := BuildIndex(nil)
	if len(idx) != 0 {
		t.Fatalf("expected empty index, got %d entries", len(idx))
	}
}

func TestBuildIndex_SingleChange(t *testing.T) {
	changes := []Change{
		{Path: "app.name", Type: ChangeModified, Value: "new", Before: "old"},
	}
	idx := BuildIndex(changes)
	if len(idx) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(idx))
	}
	c, ok := idx.Lookup("app.name")
	if !ok {
		t.Fatal("expected path to be found")
	}
	if c.Value != "new" {
		t.Errorf("expected value 'new', got %v", c.Value)
	}
}

func TestBuildIndex_LastWinsOnDuplicate(t *testing.T) {
	changes := []Change{
		{Path: "x", Type: ChangeAdded, Value: "first"},
		{Path: "x", Type: ChangeAdded, Value: "last"},
	}
	idx := BuildIndex(changes)
	c, _ := idx.Lookup("x")
	if c.Value != "last" {
		t.Errorf("expected 'last', got %v", c.Value)
	}
}

func TestIndex_Lookup_Missing(t *testing.T) {
	idx := BuildIndex([]Change{})
	_, ok := idx.Lookup("missing.key")
	if ok {
		t.Error("expected false for missing key")
	}
}

func TestIndex_Paths(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeRemoved},
	}
	idx := BuildIndex(changes)
	paths := idx.Paths()
	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(paths))
	}
}

func TestIndex_Filter_ByType(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeRemoved},
		{Path: "c", Type: ChangeModified},
	}
	idx := BuildIndex(changes)
	filtered := idx.Filter(ChangeAdded, ChangeModified)
	if len(filtered) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(filtered))
	}
	if _, ok := filtered.Lookup("b"); ok {
		t.Error("removed change should have been filtered out")
	}
}

func TestIndex_Filter_NoTypes(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
	}
	idx := BuildIndex(changes)
	filtered := idx.Filter()
	if len(filtered) != 0 {
		t.Errorf("expected empty index, got %d entries", len(filtered))
	}
}
