package diff

import (
	"testing"
)

func makeBucketChanges() []Change {
	return []Change{
		{Type: ChangeAdded, Path: "a.x", After: 1},
		{Type: ChangeRemoved, Path: "b.y", Before: 2},
		{Type: ChangeModified, Path: "c.z", Before: 3, After: 4},
		{Type: ChangeAdded, Path: "d.w", After: 5},
	}
}

func TestBucketChanges_Empty(t *testing.T) {
	buckets := BucketChanges(nil, DefaultBucketOptions())
	if len(buckets) != 0 {
		t.Fatalf("expected 0 buckets, got %d", len(buckets))
	}
}

func TestBucketChanges_ByType(t *testing.T) {
	changes := makeBucketChanges()
	buckets := BucketChanges(changes, DefaultBucketOptions())

	if len(buckets) != 3 {
		t.Fatalf("expected 3 buckets, got %d", len(buckets))
	}

	// First bucket should be "added" with 2 entries
	if buckets[0].Name != string(ChangeAdded) {
		t.Errorf("expected first bucket name %q, got %q", string(ChangeAdded), buckets[0].Name)
	}
	if len(buckets[0].Changes) != 2 {
		t.Errorf("expected 2 added changes, got %d", len(buckets[0].Changes))
	}
}

func TestBucketChanges_CustomKeyFunc(t *testing.T) {
	changes := makeBucketChanges()
	opts := DefaultBucketOptions()
	opts.KeyFunc = func(c Change) string {
		if c.Type == ChangeAdded || c.Type == ChangeModified {
			return "present"
		}
		return "absent"
	}

	buckets := BucketChanges(changes, opts)
	if len(buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(buckets))
	}
	if buckets[0].Name != "present" {
		t.Errorf("expected first bucket 'present', got %q", buckets[0].Name)
	}
	if len(buckets[0].Changes) != 3 {
		t.Errorf("expected 3 present changes, got %d", len(buckets[0].Changes))
	}
	if len(buckets[1].Changes) != 1 {
		t.Errorf("expected 1 absent change, got %d", len(buckets[1].Changes))
	}
}

func TestBucketChanges_PreservesInsertionOrder(t *testing.T) {
	changes := makeBucketChanges()
	buckets := BucketChanges(changes, DefaultBucketOptions())

	expected := []string{string(ChangeAdded), string(ChangeRemoved), string(ChangeModified)}
	for i, b := range buckets {
		if b.Name != expected[i] {
			t.Errorf("bucket[%d]: expected %q, got %q", i, expected[i], b.Name)
		}
	}
}

func TestBucketChanges_SingleType(t *testing.T) {
	changes := []Change{
		{Type: ChangeAdded, Path: "x", After: 1},
		{Type: ChangeAdded, Path: "y", After: 2},
	}
	buckets := BucketChanges(changes, DefaultBucketOptions())
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(buckets))
	}
	if len(buckets[0].Changes) != 2 {
		t.Errorf("expected 2 changes in bucket, got %d", len(buckets[0].Changes))
	}
}
