package formatter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

func TestTransformFormatter_NoChanges(t *testing.T) {
	f := NewTransformFormatter(nil)
	var buf bytes.Buffer
	if err := f.Format(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "Transform Report") {
		t.Error("expected header in output")
	}
}

func TestTransformFormatter_PathRenamed(t *testing.T) {
	orig := []diff.Change{
		{Type: diff.ChangeModified, Path: "app.host", OldValue: "a", NewValue: "b"},
	}
	transformed := []diff.Change{
		{Type: diff.ChangeModified, Path: "app.hostname", OldValue: "a", NewValue: "b"},
	}
	f := NewTransformFormatter(orig)
	var buf bytes.Buffer
	if err := f.Format(&buf, transformed); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "renamed from app.host") {
		t.Errorf("expected rename note, got:\n%s", out)
	}
	if !strings.Contains(out, "app.hostname") {
		t.Errorf("expected new path in output, got:\n%s", out)
	}
}

func TestTransformFormatter_ValueMasked(t *testing.T) {
	orig := []diff.Change{
		{Type: diff.ChangeModified, Path: "app.secret", OldValue: "plain", NewValue: "newplain"},
	}
	transformed := []diff.Change{
		{Type: diff.ChangeModified, Path: "app.secret", OldValue: "***", NewValue: "***"},
	}
	f := NewTransformFormatter(orig)
	var buf bytes.Buffer
	if err := f.Format(&buf, transformed); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "old:") {
		t.Errorf("expected old value note, got:\n%s", out)
	}
}

func TestTransformFormatter_NoRename_NoNote(t *testing.T) {
	orig := []diff.Change{
		{Type: diff.ChangeAdded, Path: "db.port", NewValue: 5432},
	}
	transformed := []diff.Change{
		{Type: diff.ChangeAdded, Path: "db.port", NewValue: 5432},
	}
	f := NewTransformFormatter(orig)
	var buf bytes.Buffer
	if err := f.Format(&buf, transformed); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "renamed") {
		t.Errorf("unexpected rename note when path unchanged:\n%s", out)
	}
}
