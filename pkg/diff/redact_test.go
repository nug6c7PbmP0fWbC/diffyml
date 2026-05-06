package diff

import (
	"testing"
)

func TestRedact_NoRules(t *testing.T) {
	changes := []Change{
		{Path: "db.password", Type: ChangeModified, OldValue: "secret", NewValue: "newsecret"},
	}
	out := Redact(changes, nil)
	if out[0].OldValue != "secret" || out[0].NewValue != "newsecret" {
		t.Errorf("expected values unchanged, got old=%v new=%v", out[0].OldValue, out[0].NewValue)
	}
}

func TestRedact_ExactMatch(t *testing.T) {
	changes := []Change{
		{Path: "db.password", Type: ChangeModified, OldValue: "secret", NewValue: "newsecret"},
	}
	rules := []RedactRule{{PathPrefix: "db.password"}}
	out := Redact(changes, rules)
	if out[0].OldValue != defaultPlaceholder || out[0].NewValue != defaultPlaceholder {
		t.Errorf("expected redacted values, got old=%v new=%v", out[0].OldValue, out[0].NewValue)
	}
}

func TestRedact_PrefixMatch(t *testing.T) {
	changes := []Change{
		{Path: "credentials.api_key", Type: ChangeAdded, OldValue: nil, NewValue: "abc123"},
	}
	rules := []RedactRule{{PathPrefix: "credentials", Placeholder: "***"}}
	out := Redact(changes, rules)
	if out[0].NewValue != "***" {
		t.Errorf("expected *** got %v", out[0].NewValue)
	}
}

func TestRedact_NoFalsePositive(t *testing.T) {
	changes := []Change{
		{Path: "credentialsExtra", Type: ChangeModified, OldValue: "x", NewValue: "y"},
	}
	rules := []RedactRule{{PathPrefix: "credentials"}}
	out := Redact(changes, rules)
	if out[0].OldValue != "x" || out[0].NewValue != "y" {
		t.Errorf("false positive redaction on path %s", out[0].Path)
	}
}

func TestRedact_CustomPlaceholder(t *testing.T) {
	changes := []Change{
		{Path: "auth.token", Type: ChangeModified, OldValue: "old", NewValue: "new"},
	}
	rules := []RedactRule{{PathPrefix: "auth.token", Placeholder: "<hidden>"}}
	out := Redact(changes, rules)
	if out[0].OldValue != "<hidden>" || out[0].NewValue != "<hidden>" {
		t.Errorf("unexpected placeholder: old=%v new=%v", out[0].OldValue, out[0].NewValue)
	}
}

func TestRedact_OriginalUnchanged(t *testing.T) {
	original := []Change{
		{Path: "secret.key", Type: ChangeModified, OldValue: "a", NewValue: "b"},
	}
	rules := []RedactRule{{PathPrefix: "secret.key"}}
	Redact(original, rules)
	if original[0].OldValue != "a" {
		t.Error("Redact must not mutate the original slice")
	}
}
