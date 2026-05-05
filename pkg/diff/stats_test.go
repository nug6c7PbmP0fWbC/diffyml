package diff

import (
	"testing"
)

func makeChanges() []Change {
	return []Change{
		{Path: "a", Type: ChangeAdded, NewValue: "1"},
		{Path: "b", Type: ChangeAdded, NewValue: "2"},
		{Path: "c", Type: ChangeRemoved, OldValue: "3"},
		{Path: "d", Type: ChangeModified, OldValue: "4", NewValue: "5"},
	}
}

func TestCollectStats_Empty(t *testing.T) {
	s := CollectStats(nil)
	if s.Total != 0 || s.Added != 0 || s.Removed != 0 || s.Modified != 0 {
		t.Errorf("expected all zeros, got %s", s)
	}
}

func TestCollectStats_Mixed(t *testing.T) {
	s := CollectStats(makeChanges())
	if s.Total != 4 {
		t.Errorf("expected total=4, got %d", s.Total)
	}
	if s.Added != 2 {
		t.Errorf("expected added=2, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("expected removed=1, got %d", s.Removed)
	}
	if s.Modified != 1 {
		t.Errorf("expected modified=1, got %d", s.Modified)
	}
}

func TestStats_Percent(t *testing.T) {
	s := CollectStats(makeChanges())

	got := s.Percent(ChangeAdded)
	if got != 50.0 {
		t.Errorf("expected 50.0%% for added, got %.2f", got)
	}

	got = s.Percent(ChangeRemoved)
	if got != 25.0 {
		t.Errorf("expected 25.0%% for removed, got %.2f", got)
	}

	got = s.Percent(ChangeModified)
	if got != 25.0 {
		t.Errorf("expected 25.0%% for modified, got %.2f", got)
	}
}

func TestStats_Percent_ZeroTotal(t *testing.T) {
	s := Stats{}
	if s.Percent(ChangeAdded) != 0 {
		t.Error("expected 0 percent when total is 0")
	}
}

func TestStats_String(t *testing.T) {
	s := Stats{Total: 4, Added: 2, Removed: 1, Modified: 1}
	expected := "total=4 added=2 removed=1 modified=1"
	if s.String() != expected {
		t.Errorf("expected %q, got %q", expected, s.String())
	}
}
