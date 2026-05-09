package diff

import (
	"testing"
)

func TestPin_NoRules(t *testing.T) {
	changes := []Change{
		{Path: "a.b", Type: ChangeModified, From: "x", To: "y"},
	}
	v := Pin(changes, nil)
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %d", len(v))
	}
}

func TestPin_NoChanges(t *testing.T) {
	rules := []PinRule{{Path: "a.b"}}
	v := Pin(nil, rules)
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %d", len(v))
	}
}

func TestPin_ExactMatch(t *testing.T) {
	changes := []Change{
		{Path: "a.b", Type: ChangeModified, From: "1", To: "2"},
	}
	rules := []PinRule{{Path: "a.b", Message: "locked"}}
	v := Pin(changes, rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Message != "locked" {
		t.Errorf("expected message 'locked', got %q", v[0].Message)
	}
	if v[0].Change.Path != "a.b" {
		t.Errorf("unexpected change path %q", v[0].Change.Path)
	}
}

func TestPin_PrefixMatch(t *testing.T) {
	changes := []Change{
		{Path: "infra.region", Type: ChangeModified, From: "us-east-1", To: "eu-west-1"},
		{Path: "infra.az", Type: ChangeAdded, To: "a"},
		{Path: "app.name", Type: ChangeModified, From: "old", To: "new"},
	}
	rules := []PinRule{{Path: "infra.", Prefix: true, Message: "infra is frozen"}}
	v := Pin(changes, rules)
	if len(v) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(v))
	}
}

func TestPin_NoFalsePositive(t *testing.T) {
	changes := []Change{
		{Path: "other.key", Type: ChangeModified, From: "a", To: "b"},
	}
	rules := []PinRule{{Path: "pinned.key"}}
	v := Pin(changes, rules)
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %d", len(v))
	}
}

func TestPin_DefaultMessage(t *testing.T) {
	changes := []Change{
		{Path: "x", Type: ChangeRemoved, From: "v"},
	}
	rules := []PinRule{{Path: "x"}}
	v := Pin(changes, rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Message == "" {
		t.Error("expected a default message but got empty string")
	}
}
