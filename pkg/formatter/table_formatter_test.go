package formatter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/diffyml/diffyml/pkg/diff"
)

func TestTableFormatter_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	f := NewTableFormatter(&buf)
	if err := f.Format(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No changes detected.") {
		t.Errorf("expected no-changes message, got: %s", buf.String())
	}
}

func TestTableFormatter_Added(t *testing.T) {
	var buf bytes.Buffer
	f := NewTableFormatter(&buf)
	changes := []diff.Change{
		{Path: "server.host", Type: diff.Added, NewValue: "localhost"},
	}
	if err := f.Format(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "server.host") {
		t.Errorf("expected path in output, got: %s", out)
	}
	if !strings.Contains(out, "added") {
		t.Errorf("expected type 'added' in output, got: %s", out)
	}
	if !strings.Contains(out, "localhost") {
		t.Errorf("expected new value in output, got: %s", out)
	}
}

func TestTableFormatter_Removed(t *testing.T) {
	var buf bytes.Buffer
	f := NewTableFormatter(&buf)
	changes := []diff.Change{
		{Path: "db.port", Type: diff.Removed, OldValue: 5432},
	}
	if err := f.Format(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "db.port") {
		t.Errorf("expected path in output, got: %s", out)
	}
	if !strings.Contains(out, "removed") {
		t.Errorf("expected type 'removed' in output, got: %s", out)
	}
	// also verify the old value is shown for removed entries
	if !strings.Contains(out, "5432") {
		t.Errorf("expected old value '5432' in output for removed entry, got: %s", out)
	}
}

func TestTableFormatter_Modified(t *testing.T) {
	var buf bytes.Buffer
	f := NewTableFormatter(&buf)
	changes := []diff.Change{
		{Path: "app.timeout", Type: diff.Modified, OldValue: 30, NewValue: 60},
	}
	if err := f.Format(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "modified") {
		t.Errorf("expected type 'modified' in output, got: %s", out)
	}
	if !strings.Contains(out, "30") || !strings.Contains(out, "60") {
		t.Errorf("expected old and new values in output, got: %s", out)
	}
}
