package diff

import (
	"testing"
)

func makeSliceChanges(n int) []Change {
	out := make([]Change, n)
	for i := 0; i < n; i++ {
		out[i] = Change{
			Path:  fmt.Sprintf("key%d", i),
			Type:  ChangeModified,
			Before: i,
			After:  i + 1,
		}
	}
	return out
}

func TestSlice_Empty(t *testing.T) {
	result := Slice([]Change{}, DefaultSliceOptions())
	if len(result) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(result))
	}
}

func TestSlice_FullRange(t *testing.T) {
	changes := makeSliceChanges(5)
	result := Slice(changes, DefaultSliceOptions())
	if len(result) != 5 {
		t.Fatalf("expected 5, got %d", len(result))
	}
}

func TestSlice_StartOnly(t *testing.T) {
	changes := makeSliceChanges(5)
	result := Slice(changes, SliceOptions{Start: 2, End: 0})
	if len(result) != 3 {
		t.Fatalf("expected 3, got %d", len(result))
	}
	if result[0].Path != "key2" {
		t.Errorf("expected key2, got %s", result[0].Path)
	}
}

func TestSlice_StartAndEnd(t *testing.T) {
	changes := makeSliceChanges(6)
	result := Slice(changes, SliceOptions{Start: 1, End: 4})
	if len(result) != 3 {
		t.Fatalf("expected 3, got %d", len(result))
	}
	if result[0].Path != "key1" {
		t.Errorf("expected key1, got %s", result[0].Path)
	}
	if result[2].Path != "key3" {
		t.Errorf("expected key3, got %s", result[2].Path)
	}
}

func TestSlice_ClampNegativeStart(t *testing.T) {
	changes := makeSliceChanges(3)
	result := Slice(changes, SliceOptions{Start: -5, End: 2})
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestSlice_EndBeyondLen(t *testing.T) {
	changes := makeSliceChanges(3)
	result := Slice(changes, SliceOptions{Start: 0, End: 100})
	if len(result) != 3 {
		t.Fatalf("expected 3, got %d", len(result))
	}
}

func TestSlice_DoesNotMutateOriginal(t *testing.T) {
	changes := makeSliceChanges(4)
	result := Slice(changes, SliceOptions{Start: 0, End: 2})
	result[0].Path = "mutated"
	if changes[0].Path == "mutated" {
		t.Error("Slice mutated the original slice")
	}
}
