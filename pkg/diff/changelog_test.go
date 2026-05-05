package diff

import (
	"strings"
	"testing"
)

func TestChangelog_AddEntry_Empty(t *testing.T) {
	cl := &Changelog{}
	cl.AddEntry("v1.0.0", []Change{})
	if len(cl.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(cl.Entries))
	}
	if cl.Entries[0].Version != "v1.0.0" {
		t.Errorf("expected version v1.0.0, got %s", cl.Entries[0].Version)
	}
	if cl.Entries[0].Summary.Added != 0 || cl.Entries[0].Summary.Removed != 0 || cl.Entries[0].Summary.Modified != 0 {
		t.Errorf("expected zero summary for empty changes")
	}
}

func TestChangelog_AddEntry_WithChanges(t *testing.T) {
	cl := &Changelog{}
	changes := []Change{
		{Type: ChangeAdded, Path: "a.b", NewValue: "x"},
		{Type: ChangeRemoved, Path: "c", OldValue: "y"},
		{Type: ChangeModified, Path: "d", OldValue: 1, NewValue: 2},
	}
	cl.AddEntry("v1.1.0", changes)
	s := cl.Entries[0].Summary
	if s.Added != 1 || s.Removed != 1 || s.Modified != 1 {
		t.Errorf("unexpected summary: %+v", s)
	}
}

func TestChangelog_MultipleEntries(t *testing.T) {
	cl := &Changelog{}
	cl.AddEntry("v1.0.0", []Change{})
	cl.AddEntry("v1.1.0", []Change{
		{Type: ChangeAdded, Path: "key", NewValue: "val"},
	})
	if len(cl.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(cl.Entries))
	}
}

func TestChangelog_Render_ContainsVersion(t *testing.T) {
	cl := &Changelog{}
	cl.AddEntry("v2.0.0", []Change{
		{Type: ChangeAdded, Path: "foo", NewValue: "bar"},
	})
	out := cl.Render()
	if !strings.Contains(out, "v2.0.0") {
		t.Errorf("expected version in output, got:\n%s", out)
	}
	if !strings.Contains(out, "[ADDED]") {
		t.Errorf("expected ADDED label in output, got:\n%s", out)
	}
	if !strings.Contains(out, "foo") {
		t.Errorf("expected path 'foo' in output, got:\n%s", out)
	}
}

func TestChangelog_Render_Empty(t *testing.T) {
	cl := &Changelog{}
	out := cl.Render()
	if out != "" {
		t.Errorf("expected empty output for empty changelog, got: %q", out)
	}
}
