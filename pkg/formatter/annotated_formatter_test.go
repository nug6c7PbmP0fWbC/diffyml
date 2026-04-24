package formatter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

func TestAnnotatedFormatter_NoAnnotations(t *testing.T) {
	var buf bytes.Buffer
	inner := NewTextFormatter(&buf)
	af := NewAnnotatedFormatter(inner, nil, &buf)

	changes := []diff.Change{
		{Path: "a", Type: diff.ChangeTypeAdded, NewValue: "1"},
	}
	if err := af.Format(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "Annotations") {
		t.Fatal("expected no annotation section")
	}
}

func TestAnnotatedFormatter_WithAnnotations(t *testing.T) {
	var buf bytes.Buffer
	inner := NewTextFormatter(&buf)
	rules := []diff.AnnotationRule{
		{PathPrefix: "db", Severity: "critical", Message: "database change detected"},
	}
	af := NewAnnotatedFormatter(inner, rules, &buf)

	changes := []diff.Change{
		{Path: "db.host", Type: diff.ChangeTypeModified, OldValue: "localhost", NewValue: "prod-db"},
	}
	if err := af.Format(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Annotations") {
		t.Fatal("expected annotation section")
	}
	if !strings.Contains(out, "CRITICAL") {
		t.Fatal("expected CRITICAL severity")
	}
	if !strings.Contains(out, "database change detected") {
		t.Fatal("expected annotation message")
	}
}

func TestAnnotatedFormatter_SeverityUppercase(t *testing.T) {
	var buf bytes.Buffer
	inner := NewTextFormatter(&buf)
	rules := []diff.AnnotationRule{
		{PathPrefix: "", Severity: "warning", Message: "generic"},
	}
	af := NewAnnotatedFormatter(inner, rules, &buf)

	changes := []diff.Change{
		{Path: "x", Type: diff.ChangeTypeRemoved, OldValue: "v"},
	}
	_ = af.Format(changes)
	out := buf.String()
	if !strings.Contains(out, "WARNING") {
		t.Fatalf("expected WARNING in output, got: %s", out)
	}
}
