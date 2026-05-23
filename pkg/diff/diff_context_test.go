package diff

import (
	"testing"
)

func makeContextChanges(paths ...string) []Change {
	out := make([]Change, len(paths))
	for i, p := range paths {
		out[i] = Change{Path: p, Type: Modified, Before: "a", After: "b"}
	}
	return out
}

func TestWithContext_Empty(t *testing.T) {
	result := WithContext(nil, DefaultContextOptions())
	if len(result) != 0 {
		t.Fatalf("expected empty, got %d", len(result))
	}
}

func TestWithContext_NoNeighbours(t *testing.T) {
	changes := makeContextChanges("a", "b", "c", "d", "e")
	opts := ContextOptions{Lines: 0}
	result := WithContext(changes, opts)
	// With Lines=0 every element is its own window; all elements included once.
	if len(result) != len(changes) {
		t.Fatalf("expected %d, got %d", len(changes), len(result))
	}
}

func TestWithContext_FullWindow(t *testing.T) {
	changes := makeContextChanges("a", "b", "c", "d", "e")
	opts := ContextOptions{Lines: 2}
	result := WithContext(changes, opts)
	// With Lines=2 the windows overlap and the whole slice is returned once.
	if len(result) != len(changes) {
		t.Fatalf("expected %d, got %d", len(changes), len(result))
	}
}

func TestWithContext_PreservesOrder(t *testing.T) {
	changes := makeContextChanges("a", "b", "c", "d", "e")
	opts := ContextOptions{Lines: 1}
	result := WithContext(changes, opts)
	for i := 1; i < len(result); i++ {
		if result[i].Path < result[i-1].Path {
			t.Fatalf("order not preserved at index %d", i)
		}
	}
}

func TestWithContext_NegativeLinesClampedToZero(t *testing.T) {
	changes := makeContextChanges("x", "y", "z")
	opts := ContextOptions{Lines: -5}
	result := WithContext(changes, opts)
	if len(result) != len(changes) {
		t.Fatalf("expected %d, got %d", len(changes), len(result))
	}
}

func TestDefaultContextOptions(t *testing.T) {
	opts := DefaultContextOptions()
	if opts.Lines != 2 {
		t.Fatalf("expected Lines=2, got %d", opts.Lines)
	}
}
