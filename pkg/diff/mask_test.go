package diff

import (
	"testing"
)

func TestMask_NoRules(t *testing.T) {
	changes := []Change{
		{Path: "db.password", Type: ChangeModified, OldValue: "secret", NewValue: "newsecret"},
	}
	opts := DefaultMaskOptions()
	out := Mask(changes, opts)
	if out[0].OldValue != "secret" {
		t.Errorf("expected original value, got %v", out[0].OldValue)
	}
}

func TestMask_ExactMatch(t *testing.T) {
	changes := []Change{
		{Path: "db.password", Type: ChangeModified, OldValue: "secret", NewValue: "newsecret"},
		{Path: "db.host", Type: ChangeModified, OldValue: "localhost", NewValue: "remotehost"},
	}
	opts := DefaultMaskOptions()
	opts.Paths = []string{"db.password"}
	out := Mask(changes, opts)
	if out[0].OldValue != "***" || out[0].NewValue != "***" {
		t.Errorf("expected masked value, got old=%v new=%v", out[0].OldValue, out[0].NewValue)
	}
	if out[1].OldValue != "localhost" {
		t.Errorf("expected unmasked host, got %v", out[1].OldValue)
	}
}

func TestMask_PrefixMatch(t *testing.T) {
	changes := []Change{
		{Path: "secrets.api_key", Type: ChangeAdded, OldValue: nil, NewValue: "abc123"},
		{Path: "secrets.token", Type: ChangeAdded, OldValue: nil, NewValue: "tok456"},
		{Path: "public.name", Type: ChangeModified, OldValue: "old", NewValue: "new"},
	}
	opts := DefaultMaskOptions()
	opts.Paths = []string{"secrets."}
	opts.PrefixMatch = true
	out := Mask(changes, opts)
	if out[0].NewValue != "***" {
		t.Errorf("expected masked api_key, got %v", out[0].NewValue)
	}
	if out[1].NewValue != "***" {
		t.Errorf("expected masked token, got %v", out[1].NewValue)
	}
	if out[2].NewValue != "new" {
		t.Errorf("expected unmasked public.name, got %v", out[2].NewValue)
	}
}

func TestMask_CustomPlaceholder(t *testing.T) {
	changes := []Change{
		{Path: "key", Type: ChangeModified, OldValue: "v1", NewValue: "v2"},
	}
	opts := DefaultMaskOptions()
	opts.Paths = []string{"key"}
	opts.Placeholder = "[REDACTED]"
	out := Mask(changes, opts)
	if out[0].OldValue != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %v", out[0].OldValue)
	}
}

func TestMask_PreservesTypeAndPath(t *testing.T) {
	changes := []Change{
		{Path: "auth.secret", Type: ChangeAdded, OldValue: nil, NewValue: "mytoken"},
	}
	opts := DefaultMaskOptions()
	opts.Paths = []string{"auth.secret"}
	out := Mask(changes, opts)
	if out[0].Path != "auth.secret" {
		t.Errorf("path should be preserved, got %v", out[0].Path)
	}
	if out[0].Type != ChangeAdded {
		t.Errorf("type should be preserved, got %v", out[0].Type)
	}
}
