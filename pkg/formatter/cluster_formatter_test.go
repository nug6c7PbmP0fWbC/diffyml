package formatter

import (
	"strings"
	"testing"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

func TestClusterFormatter_NoClusters(t *testing.T) {
	var sb strings.Builder
	f := NewClusterFormatter(&sb)
	if err := f.Format(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "no clusters") {
		t.Errorf("expected 'no clusters', got: %q", sb.String())
	}
}

func TestClusterFormatter_SingleCluster(t *testing.T) {
	clusters := map[string][]diff.Change{
		"db.host": {
			{Path: "db.host", Type: diff.ChangeModified, OldValue: "localhost", NewValue: "prod-db"},
		},
	}
	var sb strings.Builder
	f := NewClusterFormatter(&sb)
	if err := f.Format(clusters); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "db.host") {
		t.Errorf("expected cluster key in output, got: %q", out)
	}
	if !strings.Contains(out, "prod-db") {
		t.Errorf("expected new value in output, got: %q", out)
	}
}

func TestClusterFormatter_MultipleClusters_Sorted(t *testing.T) {
	clusters := map[string][]diff.Change{
		"z.key": {{Path: "z.key", Type: diff.ChangeAdded, NewValue: 1}},
		"a.key": {{Path: "a.key", Type: diff.ChangeRemoved, OldValue: 2}},
	}
	var sb strings.Builder
	f := NewClusterFormatter(&sb)
	if err := f.Format(clusters); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	aIdx := strings.Index(out, "a.key")
	zIdx := strings.Index(out, "z.key")
	if aIdx > zIdx {
		t.Errorf("expected 'a.key' before 'z.key' in output")
	}
}

func TestClusterFormatter_ChangeTypes(t *testing.T) {
	clusters := map[string][]diff.Change{
		"root": {
			{Path: "x", Type: diff.ChangeAdded, NewValue: "v1"},
			{Path: "y", Type: diff.ChangeRemoved, OldValue: "v2"},
			{Path: "z", Type: diff.ChangeModified, OldValue: "old", NewValue: "new"},
		},
	}
	var sb strings.Builder
	f := NewClusterFormatter(&sb)
	if err := f.Format(clusters); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	for _, marker := range []string{"+", "-", "~"} {
		if !strings.Contains(out, marker) {
			t.Errorf("expected marker %q in output", marker)
		}
	}
}
