package formatter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

func TestMaskFormatter_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	f := NewMaskFormatter(&buf)
	if err := f.Format(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes', got %q", buf.String())
	}
}

func TestMaskFormatter_MaskedValue(t *testing.T) {
	var buf bytes.Buffer
	f := NewMaskFormatter(&buf)
	changes := []diff.Change{
		{Path: "db.password", Type: diff.ChangeModified, OldValue: "***", NewValue: "***"},
	}
	if err := f.Format(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[masked]") {
		t.Errorf("expected [masked] tag, got %q", out)
	}
	if !strings.Contains(out, "db.password") {
		t.Errorf("expected path in output, got %q", out)
	}
}

func TestMaskFormatter_UnmaskedValue(t *testing.T) {
	var buf bytes.Buffer
	f := NewMaskFormatter(&buf)
	changes := []diff.Change{
		{Path: "app.name", Type: diff.ChangeModified, OldValue: "foo", NewValue: "bar"},
	}
	if err := f.Format(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "[masked]") {
		t.Errorf("did not expect [masked] tag, got %q", out)
	}
	if !strings.Contains(out, "foo") || !strings.Contains(out, "bar") {
		t.Errorf("expected original values in output, got %q", out)
	}
}

func TestMaskFormatter_AddedMasked(t *testing.T) {
	var buf bytes.Buffer
	f := NewMaskFormatter(&buf)
	changes := []diff.Change{
		{Path: "secrets.token", Type: diff.ChangeAdded, OldValue: nil, NewValue: "***"},
	}
	if err := f.Format(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[masked]") {
		t.Errorf("expected [masked] for added masked value, got %q", out)
	}
	if !strings.Contains(out, "+") {
		t.Errorf("expected '+' prefix for added change, got %q", out)
	}
}
