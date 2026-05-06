package formatter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

func TestLintFormatter_NoViolations(t *testing.T) {
	var buf bytes.Buffer
	f := NewLintFormatter(&buf)
	if err := f.Format(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no violations") {
		t.Errorf("expected clean-pass message, got: %s", buf.String())
	}
}

func TestLintFormatter_SingleViolation(t *testing.T) {
	var buf bytes.Buffer
	f := NewLintFormatter(&buf)
	violations := []diff.LintViolation{
		{Rule: "no-removals", Path: "db.host", Detail: "removal not permitted"},
	}
	if err := f.Format(violations); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "1 violation") {
		t.Errorf("expected violation count in output, got: %s", out)
	}
	if !strings.Contains(out, "no-removals") {
		t.Errorf("expected rule name in output, got: %s", out)
	}
	if !strings.Contains(out, "db.host") {
		t.Errorf("expected path in output, got: %s", out)
	}
}

func TestLintFormatter_MultipleViolations(t *testing.T) {
	var buf bytes.Buffer
	f := NewLintFormatter(&buf)
	violations := []diff.LintViolation{
		{Rule: "no-nil-values", Path: "service.port", Detail: "new value is nil"},
		{Rule: "no-removals", Path: "service.name", Detail: "removal not permitted"},
	}
	if err := f.Format(violations); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "2 violation") {
		t.Errorf("expected 2 violations in output, got: %s", out)
	}
	if !strings.Contains(out, "no-nil-values") {
		t.Errorf("expected first rule in output")
	}
	if !strings.Contains(out, "no-removals") {
		t.Errorf("expected second rule in output")
	}
}
