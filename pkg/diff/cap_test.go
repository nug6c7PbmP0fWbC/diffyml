package diff

import (
	"testing"
)

func makeCap(n int, ct ChangeType) []Change {
	out := make([]Change, n)
	for i := range out {
		out[i] = Change{Path: fmt.Sprintf("key%d", i), Type: ct}
	}
	return out
}

func TestCap_Empty(t *testing.T) {
	result := Cap(nil, DefaultCapOptions())
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestCap_NoLimit(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeRemoved},
	}
	result := Cap(changes, DefaultCapOptions())
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestCap_MaxCaps(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeAdded},
		{Path: "c", Type: ChangeAdded},
	}
	opts := DefaultCapOptions()
	opts.MaxChanges = 2
	result := Cap(changes, opts)
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	if result[0].Path != "a" || result[1].Path != "b" {
		t.Fatalf("unexpected paths: %v", result)
	}
}

func TestCap_PriorityPreservesHigherTier(t *testing.T) {
	changes := []Change{
		{Path: "r1", Type: ChangeRemoved},
		{Path: "a1", Type: ChangeAdded},
		{Path: "r2", Type: ChangeRemoved},
		{Path: "a2", Type: ChangeAdded},
	}
	opts := CapOptions{
		MaxChanges: 3,
		Priority:   []ChangeType{ChangeAdded, ChangeRemoved},
	}
	result := Cap(changes, opts)
	if len(result) != 3 {
		t.Fatalf("expected 3, got %d", len(result))
	}
	// Both added changes come first.
	if result[0].Type != ChangeAdded || result[1].Type != ChangeAdded {
		t.Fatalf("expected added changes first, got %v", result)
	}
	if result[2].Type != ChangeRemoved {
		t.Fatalf("expected removed last, got %v", result[2])
	}
}

func TestCap_DoesNotMutateOriginal(t *testing.T) {
	changes := []Change{
		{Path: "x", Type: ChangeModified},
		{Path: "y", Type: ChangeModified},
		{Path: "z", Type: ChangeModified},
	}
	opts := CapOptions{MaxChanges: 2}
	_ = Cap(changes, opts)
	if len(changes) != 3 {
		t.Fatalf("original slice mutated")
	}
}

func TestCap_ExactLimit(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeRemoved},
	}
	opts := CapOptions{MaxChanges: 2}
	result := Cap(changes, opts)
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}
