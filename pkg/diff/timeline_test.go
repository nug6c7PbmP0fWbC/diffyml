package diff

import (
	"strings"
	"testing"
)

func TestTimeline_Empty(t *testing.T) {
	var tl Timeline
	out := tl.Render()
	if !strings.Contains(out, "empty") {
		t.Errorf("expected empty message, got %q", out)
	}
	if tl.TotalChanges() != 0 {
		t.Errorf("expected 0 total changes")
	}
}

func TestTimeline_AddEntry_NoChanges(t *testing.T) {
	var tl Timeline
	tl.AddEntry("v1.0", nil)
	if len(tl.Entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if tl.Entries[0].Label != "v1.0" {
		t.Errorf("unexpected label %q", tl.Entries[0].Label)
	}
	if tl.TotalChanges() != 0 {
		t.Errorf("expected 0 total changes")
	}
}

func TestTimeline_AddEntry_WithChanges(t *testing.T) {
	var tl Timeline
	changes := []Change{
		{Type: ChangeAdded, Path: "a.b", After: "new"},
		{Type: ChangeRemoved, Path: "c", Before: "old"},
	}
	tl.AddEntry("v1.1", changes)
	if tl.TotalChanges() != 2 {
		t.Errorf("expected 2 total changes, got %d", tl.TotalChanges())
	}
}

func TestTimeline_MultipleEntries(t *testing.T) {
	var tl Timeline
	tl.AddEntry("v1.0", []Change{{Type: ChangeAdded, Path: "x"}})
	tl.AddEntry("v1.1", []Change{
		{Type: ChangeModified, Path: "y"},
		{Type: ChangeRemoved, Path: "z"},
	})
	if tl.TotalChanges() != 3 {
		t.Errorf("expected 3, got %d", tl.TotalChanges())
	}
}

func TestTimeline_Render_ContainsLabel(t *testing.T) {
	var tl Timeline
	tl.AddEntry("release-42", []Change{
		{Type: ChangeAdded, Path: "feature.flag"},
	})
	out := tl.Render()
	if !strings.Contains(out, "release-42") {
		t.Errorf("expected label in render output")
	}
	if !strings.Contains(out, "feature.flag") {
		t.Errorf("expected path in render output")
	}
}

func TestTimeline_Render_OrderPreserved(t *testing.T) {
	var tl Timeline
	tl.AddEntry("first", nil)
	tl.AddEntry("second", nil)
	out := tl.Render()
	idxFirst := strings.Index(out, "first")
	idxSecond := strings.Index(out, "second")
	if idxFirst > idxSecond {
		t.Errorf("expected 'first' before 'second' in output")
	}
}
