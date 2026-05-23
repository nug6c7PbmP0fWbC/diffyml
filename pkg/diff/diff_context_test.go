package diff

import (
	"testing"
)

func makeContextChanges(paths ...string) []Change {
	out := make([]Change, len(paths))
	for i, p := range paths {
		out[i] = Change{Path: p, Type: ChangeTypeModified}
	}
	return out
}

func TestWithContext_Empty(t *testing.T) {
	result := WithContext(nil, nil, DefaultContextOptions())
	if len(result) != 0 {
		t.Fatalf("expected empty, got %d", len(result))
	}
}

func TestWithContext_NoNeighbours(t *testing.T) {
	all := makeContextChanges("a", "b", "c", "d", "e")
	changed := makeContextChanges("c")
	opts := ContextOptions{Lines: 0}
	result := WithContext(changed, all, opts)
	// Lines=0 returns original changes unchanged
	if len(result) != 1 || result[0].Path != "c" {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestWithContext_FullWindow(t *testing.T) {
	all := makeContextChanges("a", "b", "c", "d", "e")
	changed := makeContextChanges("c")
	opts := ContextOptions{Lines: 2}
	result := WithContext(changed, all, opts)
	if len(result) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(result))
	}
}

func TestWithContext_PreservesOrder(t *testing.T) {
	all := makeContextChanges("a", "b", "c", "d", "e")
	changed := makeContextChanges("c")
	opts := ContextOptions{Lines: 1}
	result := WithContext(changed, all, opts)
	want := []string{"b", "c", "d"}
	if len(result) != len(want) {
		t.Fatalf("expected %d entries, got %d", len(want), len(result))
	}
	for i, w := range want {
		if result[i].Path != w {
			t.Errorf("index %d: want %q got %q", i, w, result[i].Path)
		}
	}
}

func TestWithContext_EdgeBoundary(t *testing.T) {
	all := makeContextChanges("a", "b", "c")
	changed := makeContextChanges("a")
	opts := ContextOptions{Lines: 2}
	result := WithContext(changed, all, opts)
	// a + neighbours b, c (no out-of-bounds)
	if len(result) != 3 {
		t.Fatalf("expected 3, got %d", len(result))
	}
}

func TestWithContext_MultipleChanges_NoOverlap(t *testing.T) {
	all := makeContextChanges("a", "b", "c", "d", "e", "f", "g")
	changed := makeContextChanges("b", "f")
	opts := ContextOptions{Lines: 1}
	result := WithContext(changed, all, opts)
	// b -> a,b,c; f -> e,f,g => a,b,c,e,f,g
	if len(result) != 6 {
		t.Fatalf("expected 6, got %d: %+v", len(result), result)
	}
}
