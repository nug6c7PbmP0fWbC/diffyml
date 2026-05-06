package diff

import (
	"strings"
	"testing"
)

func TestLint_NoRules(t *testing.T) {
	changes := []Change{
		{Path: "a.b", Type: ChangeAdded, NewValue: "x"},
	}
	violations := Lint(changes, nil)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestLint_NoChanges(t *testing.T) {
	rules := []LintRule{NoNilValues, NoRemovals}
	violations := Lint(nil, rules)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestLint_NoNilValues_Triggered(t *testing.T) {
	changes := []Change{
		{Path: "service.port", Type: ChangeAdded, NewValue: nil},
	}
	violations := Lint(changes, []LintRule{NoNilValues})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Rule != "no-nil-values" {
		t.Errorf("unexpected rule name: %s", violations[0].Rule)
	}
	if violations[0].Path != "service.port" {
		t.Errorf("unexpected path: %s", violations[0].Path)
	}
}

func TestLint_NoNilValues_NotTriggered(t *testing.T) {
	changes := []Change{
		{Path: "service.port", Type: ChangeModified, OldValue: 80, NewValue: 443},
	}
	violations := Lint(changes, []LintRule{NoNilValues})
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestLint_NoRemovals_Triggered(t *testing.T) {
	changes := []Change{
		{Path: "database.host", Type: ChangeRemoved, OldValue: "localhost"},
	}
	violations := Lint(changes, []LintRule{NoRemovals})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Rule != "no-removals" {
		t.Errorf("unexpected rule name: %s", violations[0].Rule)
	}
}

func TestLint_MultipleRulesMultipleViolations(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeRemoved, OldValue: "v"},
		{Path: "b", Type: ChangeAdded, NewValue: nil},
	}
	violations := Lint(changes, []LintRule{NoRemovals, NoNilValues})
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(violations))
	}
}

func TestLintViolation_String(t *testing.T) {
	v := LintViolation{Rule: "no-removals", Path: "x.y", Detail: "removal not allowed"}
	s := v.String()
	if !strings.Contains(s, "no-removals") || !strings.Contains(s, "x.y") {
		t.Errorf("unexpected string representation: %s", s)
	}
}
