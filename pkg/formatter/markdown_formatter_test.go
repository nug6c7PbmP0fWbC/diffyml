package formatter

import (
	"strings"
	"testing"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

func TestMarkdownFormatter_NoChanges(t *testing.T) {
	f := NewMarkdownFormatter()
	out, err := f.Format([]diff.Change{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No changes") {
		t.Errorf("expected no-changes message, got: %s", out)
	}
}

func TestMarkdownFormatter_Added(t *testing.T) {
	f := NewMarkdownFormatter()
	changes := []diff.Change{
		{Path: "service.port", Type: diff.Added, OldValue: nil, NewValue: 8080},
	}
	out, err := f.Format(changes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "service.port") {
		t.Errorf("expected path in output, got: %s", out)
	}
	if !strings.Contains(out, string(diff.Added)) {
		t.Errorf("expected type 'added' in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected new value in output, got: %s", out)
	}
}

func TestMarkdownFormatter_Removed(t *testing.T) {
	f := NewMarkdownFormatter()
	changes := []diff.Change{
		{Path: "service.name", Type: diff.Removed, OldValue: "myapp", NewValue: nil},
	}
	out, err := f.Format(changes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "myapp") {
		t.Errorf("expected old value in output, got: %s", out)
	}
}

func TestMarkdownFormatter_Modified(t *testing.T) {
	f := NewMarkdownFormatter()
	changes := []diff.Change{
		{Path: "app.version", Type: diff.Modified, OldValue: "1.0", NewValue: "2.0"},
	}
	out, err := f.Format(changes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "1.0") || !strings.Contains(out, "2.0") {
		t.Errorf("expected both old and new values in output, got: %s", out)
	}
	if !strings.Contains(out, "|") {
		t.Errorf("expected markdown table format, got: %s", out)
	}
}
