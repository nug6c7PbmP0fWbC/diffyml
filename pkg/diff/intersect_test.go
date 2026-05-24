package diff

import (
	"testing"
)

func TestIntersect_Empty(t *testing.T) {
	result := Intersect()
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestIntersect_SingleSet(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeRemoved},
	}
	result := Intersect(changes)
	if len(result) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(result))
	}
}

func TestIntersect_CommonPaths(t *testing.T) {
	a := []Change{
		{Path: "x", Type: ChangeAdded},
		{Path: "y", Type: ChangeModified},
		{Path: "z", Type: ChangeRemoved},
	}
	b := []Change{
		{Path: "x", Type: ChangeAdded},
		{Path: "z", Type: ChangeRemoved},
	}
	result := Intersect(a, b)
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	paths := map[string]bool{}
	for _, c := range result {
		paths[c.Path] = true
	}
	if !paths["x"] || !paths["z"] {
		t.Errorf("expected paths x and z, got %v", paths)
	}
}

func TestIntersect_NoCommonPaths(t *testing.T) {
	a := []Change{{Path: "a", Type: ChangeAdded}}
	b := []Change{{Path: "b", Type: ChangeAdded}}
	result := Intersect(a, b)
	if len(result) != 0 {
		t.Fatalf("expected 0, got %d", len(result))
	}
}

func TestIntersect_ThreeSets(t *testing.T) {
	a := []Change{
		{Path: "p", Type: ChangeAdded},
		{Path: "q", Type: ChangeModified},
	}
	b := []Change{
		{Path: "p", Type: ChangeAdded},
		{Path: "q", Type: ChangeAdded},
	}
	c := []Change{
		{Path: "p", Type: ChangeRemoved},
	}
	result := Intersect(a, b, c)
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Path != "p" {
		t.Errorf("expected path p, got %s", result[0].Path)
	}
	// Values should come from first set
	if result[0].Type != ChangeAdded {
		t.Errorf("expected type from first set (Added), got %v", result[0].Type)
	}
}

func TestIntersect_PreservesOrder(t *testing.T) {
	a := []Change{
		{Path: "c", Type: ChangeAdded},
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeAdded},
	}
	b := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeAdded},
		{Path: "c", Type: ChangeAdded},
	}
	result := Intersect(a, b)
	if len(result) != 3 {
		t.Fatalf("expected 3, got %d", len(result))
	}
	if result[0].Path != "c" || result[1].Path != "a" || result[2].Path != "b" {
		t.Errorf("order not preserved: %v", result)
	}
}
