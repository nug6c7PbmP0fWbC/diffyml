package diff

import (
	"testing"
)

func TestRollup_Empty(t *testing.T) {
	buckets := Rollup(nil, DefaultRollupOptions())
	if len(buckets) != 0 {
		t.Fatalf("expected 0 buckets, got %d", len(buckets))
	}
}

func TestRollup_SingleBucket(t *testing.T) {
	changes := []Change{
		{Path: "server.host", Type: ChangeModified},
		{Path: "server.port", Type: ChangeAdded},
	}
	buckets := Rollup(changes, DefaultRollupOptions())
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(buckets))
	}
	b := buckets[0]
	if b.Prefix != "server" {
		t.Errorf("expected prefix 'server', got %q", b.Prefix)
	}
	if b.Total != 2 {
		t.Errorf("expected total 2, got %d", b.Total)
	}
	if b.Modified != 1 || b.Added != 1 {
		t.Errorf("unexpected counts: modified=%d added=%d", b.Modified, b.Added)
	}
}

func TestRollup_MultipleBuckets(t *testing.T) {
	changes := []Change{
		{Path: "server.host", Type: ChangeModified},
		{Path: "database.url", Type: ChangeAdded},
		{Path: "database.pool", Type: ChangeRemoved},
	}
	buckets := Rollup(changes, DefaultRollupOptions())
	if len(buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d", len(buckets))
	}
	if buckets[0].Prefix != "server" {
		t.Errorf("expected first bucket 'server', got %q", buckets[0].Prefix)
	}
	if buckets[1].Prefix != "database" {
		t.Errorf("expected second bucket 'database', got %q", buckets[1].Prefix)
	}
	if buckets[1].Total != 2 {
		t.Errorf("expected database total 2, got %d", buckets[1].Total)
	}
}

func TestRollup_Depth2(t *testing.T) {
	changes := []Change{
		{Path: "a.b.c", Type: ChangeAdded},
		{Path: "a.b.d", Type: ChangeRemoved},
		{Path: "a.x.y", Type: ChangeModified},
	}
	opts := RollupOptions{Depth: 2}
	buckets := Rollup(changes, opts)
	if len(buckets) != 2 {
		t.Fatalf("expected 2 buckets at depth 2, got %d", len(buckets))
	}
	if buckets[0].Prefix != "a.b" {
		t.Errorf("expected 'a.b', got %q", buckets[0].Prefix)
	}
}

func TestRollup_RootPath(t *testing.T) {
	changes := []Change{
		{Path: "", Type: ChangeAdded},
	}
	buckets := Rollup(changes, DefaultRollupOptions())
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(buckets))
	}
	if buckets[0].Prefix != "(root)" {
		t.Errorf("expected '(root)', got %q", buckets[0].Prefix)
	}
}

func TestRollupPrefix_NoDot(t *testing.T) {
	if got := rollupPrefix("toplevel", 1); got != "toplevel" {
		t.Errorf("expected 'toplevel', got %q", got)
	}
}
