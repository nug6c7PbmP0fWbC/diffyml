package diff

import (
	"testing"
)

func TestSort_Empty(t *testing.T) {
	result := Sort(nil, DefaultSortOptions())
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestSort_Ascending(t *testing.T) {
	changes := []Change{
		{Path: "z.key", Type: ChangeAdded},
		{Path: "a.key", Type: ChangeRemoved},
		{Path: "m.key", Type: ChangeModified},
	}
	result := Sort(changes, DefaultSortOptions())
	if result[0].Path != "a.key" || result[1].Path != "m.key" || result[2].Path != "z.key" {
		t.Fatalf("unexpected order: %v", result)
	}
}

func TestSort_Descending(t *testing.T) {
	changes := []Change{
		{Path: "a.key", Type: ChangeAdded},
		{Path: "z.key", Type: ChangeRemoved},
		{Path: "m.key", Type: ChangeModified},
	}
	opts := SortOptions{Order: SortDescending, Stable: true}
	result := Sort(changes, opts)
	if result[0].Path != "z.key" || result[1].Path != "m.key" || result[2].Path != "a.key" {
		t.Fatalf("unexpected order: %v", result)
	}
}

func TestSort_ByType(t *testing.T) {
	changes := []Change{
		{Path: "b", Type: ChangeRemoved},
		{Path: "a", Type: ChangeModified},
		{Path: "c", Type: ChangeAdded},
	}
	opts := SortOptions{Order: SortByType, Stable: true}
	result := Sort(changes, opts)
	if result[0].Type != ChangeAdded {
		t.Fatalf("expected first to be added, got %v", result[0].Type)
	}
	if result[1].Type != ChangeModified {
		t.Fatalf("expected second to be modified, got %v", result[1].Type)
	}
	if result[2].Type != ChangeRemoved {
		t.Fatalf("expected third to be removed, got %v", result[2].Type)
	}
}

func TestSort_DoesNotMutateOriginal(t *testing.T) {
	changes := []Change{
		{Path: "z", Type: ChangeAdded},
		{Path: "a", Type: ChangeRemoved},
	}
	originalFirst := changes[0].Path
	Sort(changes, DefaultSortOptions())
	if changes[0].Path != originalFirst {
		t.Fatal("Sort mutated the original slice")
	}
}

func TestSort_ByType_TieBreakByPath(t *testing.T) {
	changes := []Change{
		{Path: "z", Type: ChangeAdded},
		{Path: "a", Type: ChangeAdded},
	}
	opts := SortOptions{Order: SortByType, Stable: true}
	result := Sort(changes, opts)
	if result[0].Path != "a" {
		t.Fatalf("expected 'a' first within same type, got %v", result[0].Path)
	}
}
