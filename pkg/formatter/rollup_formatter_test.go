package formatter

import (
	"strings"
	"testing"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

func TestRollupFormatter_NoBuckets(t *testing.T) {
	var sb strings.Builder
	f := NewRollupFormatter(&sb)
	if err := f.Format(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "No changes") {
		t.Errorf("expected 'No changes' message, got: %s", sb.String())
	}
}

func TestRollupFormatter_HeaderPresent(t *testing.T) {
	var sb strings.Builder
	f := NewRollupFormatter(&sb)
	buckets := []diff.RollupBucket{
		{Prefix: "server", Added: 1, Total: 1},
	}
	if err := f.Format(buckets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	for _, col := range []string{"PREFIX", "ADDED", "REMOVED", "MODIFIED", "TOTAL"} {
		if !strings.Contains(out, col) {
			t.Errorf("missing column header %q in output:\n%s", col, out)
		}
	}
}

func TestRollupFormatter_SingleBucket(t *testing.T) {
	var sb strings.Builder
	f := NewRollupFormatter(&sb)
	buckets := []diff.RollupBucket{
		{Prefix: "database", Added: 2, Removed: 1, Modified: 3, Total: 6},
	}
	if err := f.Format(buckets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "database") {
		t.Errorf("expected 'database' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "6") {
		t.Errorf("expected total '6' in output, got:\n%s", out)
	}
}

func TestRollupFormatter_MultipleBuckets(t *testing.T) {
	var sb strings.Builder
	f := NewRollupFormatter(&sb)
	buckets := []diff.RollupBucket{
		{Prefix: "server", Added: 1, Total: 1},
		{Prefix: "logging", Removed: 2, Total: 2},
	}
	if err := f.Format(buckets); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "server") || !strings.Contains(out, "logging") {
		t.Errorf("expected both bucket prefixes in output:\n%s", out)
	}
}
