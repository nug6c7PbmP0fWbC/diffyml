package diff

import (
	"strings"
	"testing"
)

func TestTransform_Empty(t *testing.T) {
	result := Transform(nil, DefaultTransformOptions())
	if len(result) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(result))
	}
}

func TestTransform_NoRename(t *testing.T) {
	changes := []Change{
		{Type: ChangeModified, Path: "app.host", OldValue: "a", NewValue: "b"},
	}
	out := Transform(changes, DefaultTransformOptions())
	if out[0].Path != "app.host" {
		t.Errorf("expected path unchanged, got %s", out[0].Path)
	}
}

func TestTransform_KeyRename(t *testing.T) {
	changes := []Change{
		{Type: ChangeModified, Path: "app.host", OldValue: "a", NewValue: "b"},
		{Type: ChangeAdded, Path: "db.port", NewValue: 5432},
	}
	opts := DefaultTransformOptions()
	opts.KeyRename = map[string]string{"host": "hostname", "port": "listen_port"}
	out := Transform(changes, opts)
	if out[0].Path != "app.hostname" {
		t.Errorf("expected app.hostname, got %s", out[0].Path)
	}
	if out[1].Path != "db.listen_port" {
		t.Errorf("expected db.listen_port, got %s", out[1].Path)
	}
}

func TestTransform_ValueMapper(t *testing.T) {
	changes := []Change{
		{Type: ChangeModified, Path: "app.secret", OldValue: "old-secret", NewValue: "new-secret"},
	}
	opts := DefaultTransformOptions()
	opts.ValueMapper = func(path string, val interface{}) interface{} {
		if strings.Contains(path, "secret") {
			return "***"
		}
		return val
	}
	out := Transform(changes, opts)
	if out[0].OldValue != "***" {
		t.Errorf("expected masked OldValue, got %v", out[0].OldValue)
	}
	if out[0].NewValue != "***" {
		t.Errorf("expected masked NewValue, got %v", out[0].NewValue)
	}
}

func TestTransform_DoesNotMutateOriginal(t *testing.T) {
	original := []Change{
		{Type: ChangeAdded, Path: "x.key", NewValue: "v"},
	}
	opts := DefaultTransformOptions()
	opts.KeyRename = map[string]string{"key": "renamed"}
	Transform(original, opts)
	if original[0].Path != "x.key" {
		t.Errorf("original mutated: got %s", original[0].Path)
	}
}

func TestTransform_RenameOnlyLastSegment(t *testing.T) {
	changes := []Change{
		{Type: ChangeModified, Path: "host.host", OldValue: "a", NewValue: "b"},
	}
	opts := DefaultTransformOptions()
	opts.KeyRename = map[string]string{"host": "hostname"}
	out := Transform(changes, opts)
	if out[0].Path != "host.hostname" {
		t.Errorf("expected host.hostname, got %s", out[0].Path)
	}
}
