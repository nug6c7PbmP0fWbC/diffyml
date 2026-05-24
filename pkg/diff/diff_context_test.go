package diff

import (
	"testing"
)

func makeContextChanges() []Change {
	return []Change{
		{Path: "a", Type: Unchanged},
		{Path: "b", Type: Unchanged},
		{Path: "c", Type: Added},
		{Path: "d", Type: Unchanged},
		{Path: "e", Type: Unchanged},
		{Path: "f", Type: Unchanged},
		{Path: "g", Type: Modified},
		{Path: "h", Type: Unchanged},
		{Path: "i", Type: Unchanged},
	}
}

func TestWithContext_Empty(t *testing.T) {
	out := WithContext(nil, DefaultContextOptions())
	if len(out) != 0 {
		t.Fatalf("expected empty, got %d", len(out))
	}
}

func TestWithContext_NoNeighbours(t *testing.T) {
	changes := makeContextChanges()
	out := WithContext(changes, ContextOptions{Before: 0, After: 0})
	for _, c := range out {
		if c.Type != Added && c.Type != Modified {
			t.Errorf("expected only real changes, got %s at %s", c.Type, c.Path)
		}
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(out))
	}
}

func TestWithContext_FullWindow(t *testing.T) {
	changes := makeContextChanges()
	out := WithContext(changes, ContextOptions{Before: 2, After: 2})
	paths := make(map[string]bool)
	for _, c := range out {
		paths[c.Path] = true
	}
	// "c" (Added) should pull in "a","b" before and "d","e" after
	for _, p := range []string{"a", "b", "c", "d", "e"} {
		if !paths[p] {
			t.Errorf("expected path %q to be included", p)
		}
	}
	// "f" is 3 away from "c" and 1 before "g", so it should appear via "g"
	if !paths["f"] {
		t.Errorf("expected path f to be included as context before g")
	}
}

func TestWithContext_PreservesOrder(t *testing.T) {
	changes := makeContextChanges()
	out := WithContext(changes, DefaultContextOptions())
	for i := 1; i < len(out); i++ {
		if out[i].Path < out[i-1].Path {
			t.Errorf("order not preserved at index %d", i)
		}
	}
}

func TestWithContext_ContextMarker(t *testing.T) {
	changes := makeContextChanges()
	out := WithContext(changes, ContextOptions{Before: 1, After: 1})
	for _, c := range out {
		if c.Type == Added || c.Type == Removed || c.Type == Modified {
			if c.Metadata != nil && c.Metadata["context"] == true {
				t.Errorf("real change %s should not be marked as context", c.Path)
			}
		} else {
			if c.Metadata == nil || c.Metadata["context"] != true {
				t.Errorf("context change %s should be marked with context=true", c.Path)
			}
		}
	}
}

func TestWithContext_BoundaryClamp(t *testing.T) {
	changes := []Change{
		{Path: "x", Type: Removed},
		{Path: "y", Type: Unchanged},
	}
	out := WithContext(changes, ContextOptions{Before: 5, After: 5})
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}
