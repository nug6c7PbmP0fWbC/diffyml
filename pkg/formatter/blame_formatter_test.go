package formatter

import (
	"strings"
	"testing"
	"time"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

var bfTime = time.Date(2024, 6, 1, 9, 0, 0, 0, time.UTC)

func TestBlameFormatter_Empty(t *testing.T) {
	f := NewBlameFormatter()
	var sb strings.Builder
	if err := f.Format(&sb, diff.BlameMap{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "no blame entries") {
		t.Errorf("expected empty message, got: %q", sb.String())
	}
}

func TestBlameFormatter_SingleEntry(t *testing.T) {
	f := NewBlameFormatter()
	bm := diff.BlameMap{
		"server.port": {
			Path:       "server.port",
			Author:     "alice",
			ChangedAt:  bfTime,
			ChangeType: "modified",
		},
	}
	var sb strings.Builder
	if err := f.Format(&sb, bm); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "server.port") {
		t.Errorf("expected path in output, got: %q", out)
	}
	if !strings.Contains(out, "alice") {
		t.Errorf("expected author in output, got: %q", out)
	}
	if !strings.Contains(out, "modified") {
		t.Errorf("expected change type in output, got: %q", out)
	}
}

func TestBlameFormatter_MultipleEntries_Sorted(t *testing.T) {
	f := NewBlameFormatter()
	bm := diff.BlameMap{
		"z.key": {Path: "z.key", Author: "bob", ChangedAt: bfTime, ChangeType: "added"},
		"a.key": {Path: "a.key", Author: "carol", ChangedAt: bfTime, ChangeType: "removed"},
	}
	var sb strings.Builder
	if err := f.Format(&sb, bm); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	aIdx := strings.Index(out, "a.key")
	zIdx := strings.Index(out, "z.key")
	if aIdx == -1 || zIdx == -1 {
		t.Fatalf("missing entries in output: %q", out)
	}
	if aIdx > zIdx {
		t.Errorf("expected a.key before z.key in sorted output")
	}
}

func TestBlameFormatter_HeaderPresent(t *testing.T) {
	f := NewBlameFormatter()
	bm := diff.BlameMap{
		"k": {Path: "k", Author: "x", ChangedAt: bfTime, ChangeType: "added"},
	}
	var sb strings.Builder
	_ = f.Format(&sb, bm)
	out := sb.String()
	for _, col := range []string{"PATH", "TYPE", "AUTHOR", "CHANGED AT"} {
		if !strings.Contains(out, col) {
			t.Errorf("expected column header %q in output", col)
		}
	}
}
