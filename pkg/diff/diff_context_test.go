package diff

import (
	"testing"
)

func makeContextChanges() []Change {
	keys := []string{"a", "b", "c", "d", "e", "f", "g"}
	out := make([]Change, len(keys))
	for i, k := range keys {
		out[i] = Change{Path: k, Type: ChangeAdded, After: i}
	}
	// Mark only index 3 as the "real" change; rest become context candidates
	// For test purposes we'll mix types.
	out[0].Type = "unchanged"
	out[1].Type = "unchanged"
	out[2].Type = "unchanged"
	// out[3] stays ChangeAdded — the real change
	out[4].Type = "unchanged"
	out[5].Type = "unchanged"
	out[6].Type = "unchanged"
	return out
}

func TestWithContext_Empty(t *testing.T) {
	result := WithContext(nil, DefaultContextOptions())
	if len(result) != 0 {
		t.Fatalf("expected empty, got %d", len(result))
	}
}

func TestWithContext_NoNeighbours(t *testing.T) {
	changes := makeContextChanges()
	opts := ContextOptions{Before: 0, After: 0}
	result := WithContext(changes, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result))
	}
	if result[0].Path != "d" {
		t.Errorf("expected path 'd', got %q", result[0].Path)
	}
}

func TestWithContext_FullWindow(t *testing.T) {
	changes := makeContextChanges()
	opts := ContextOptions{Before: 2, After: 2}
	result := WithContext(changes, opts)
	// indices 1,2,3,4,5 should be included
	if len(result) != 5 {
		t.Fatalf("expected 5 changes, got %d", len(result))
	}
}

func TestWithContext_PreservesOrder(t *testing.T) {
	changes := makeContextChanges()
	opts := ContextOptions{Before: 1, After: 1}
	result := WithContext(changes, opts)
	for i := 1; i < len(result); i++ {
		if result[i].Path < result[i-1].Path {
			t.Errorf("order not preserved at index %d", i)
		}
	}
}

func TestWithContext_ContextMetadataSet(t *testing.T) {
	changes := makeContextChanges()
	opts := ContextOptions{Before: 1, After: 1}
	result := WithContext(changes, opts)
	for _, c := range result {
		if c.Type == ChangeAdded || c.Type == ChangeRemoved || c.Type == ChangeModified {
			if c.Metadata != nil {
				if _, ok := c.Metadata["context"]; ok {
					t.Errorf("real change should not have context=true")
				}
			}
		} else {
			if c.Metadata == nil || c.Metadata["context"] != true {
				t.Errorf("context change missing metadata for path %q", c.Path)
			}
		}
	}
}

func TestWithContext_ClampsBounds(t *testing.T) {
	changes := makeContextChanges()
	// Real change is at index 3; request 10 before — should not panic
	opts := ContextOptions{Before: 10, After: 10}
	result := WithContext(changes, opts)
	if len(result) != len(changes) {
		t.Fatalf("expected all %d changes, got %d", len(changes), len(result))
	}
}
