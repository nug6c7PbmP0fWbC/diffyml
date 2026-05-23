package diff

import (
	"testing"
)

func makeSplitChanges() []Change {
	return []Change{
		{Path: "a", Type: ChangeAdded, Value: 1},
		{Path: "b", Type: ChangeRemoved, Before: 2},
		{Path: "c", Type: ChangeModified, Before: 3, Value: 4},
		{Path: "d", Type: ChangeAdded, Value: 5},
		{Path: "e", Type: ChangeRemoved, Before: 6},
	}
}

func TestSplit_Empty(t *testing.T) {
	result := Split(nil, DefaultSplitOptions())
	if len(result) != 0 {
		t.Fatalf("expected empty map, got %d buckets", len(result))
	}
}

func TestSplit_ByType(t *testing.T) {
	changes := makeSplitChanges()
	result := Split(changes, DefaultSplitOptions())

	if len(result["added"]) != 2 {
		t.Errorf("expected 2 added, got %d", len(result["added"]))
	}
	if len(result["removed"]) != 2 {
		t.Errorf("expected 2 removed, got %d", len(result["removed"]))
	}
	if len(result["modified"]) != 1 {
		t.Errorf("expected 1 modified, got %d", len(result["modified"]))
	}
}

func TestSplit_CustomPredicate(t *testing.T) {
	changes := makeSplitChanges()
	opts := DefaultSplitOptions()
	opts.Predicate = func(c Change) string {
		if c.Type == ChangeAdded {
			return "new"
		}
		return "old"
	}
	result := Split(changes, opts)

	if len(result["new"]) != 2 {
		t.Errorf("expected 2 in 'new', got %d", len(result["new"]))
	}
	if len(result["old"]) != 3 {
		t.Errorf("expected 3 in 'old', got %d", len(result["old"]))
	}
}

func TestSplit_MaxBuckets(t *testing.T) {
	changes := makeSplitChanges()
	opts := DefaultSplitOptions()
	opts.MaxBuckets = 2
	result := Split(changes, opts)

	if len(result) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(result))
	}

	total := 0
	for _, v := range result {
		total += len(v)
	}
	if total != len(changes) {
		t.Errorf("expected all %d changes to be present, got %d", len(changes), total)
	}
}

func TestSplit_MaxBucketsOne(t *testing.T) {
	changes := makeSplitChanges()
	opts := DefaultSplitOptions()
	opts.MaxBuckets = 1
	result := Split(changes, opts)

	if len(result) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(result))
	}
	for _, v := range result {
		if len(v) != len(changes) {
			t.Errorf("expected all changes in single bucket, got %d", len(v))
		}
	}
}
