package filter_test

import (
	"testing"

	"github.com/diffyml/diffyml/pkg/diff"
	"github.com/diffyml/diffyml/pkg/filter"
)

var sampleChanges = []diff.Change{
	{Path: "server.host", Type: diff.Added, NewValue: "localhost"},
	{Path: "server.port", Type: diff.Modified, OldValue: "80", NewValue: "8080"},
	{Path: "db.name", Type: diff.Removed, OldValue: "mydb"},
	{Path: "db.port", Type: diff.Added, NewValue: "5432"},
}

func TestApply_NoFilter(t *testing.T) {
	result := filter.Apply(sampleChanges, filter.Options{})
	if len(result) != len(sampleChanges) {
		t.Fatalf("expected %d changes, got %d", len(sampleChanges), len(result))
	}
}

func TestApply_FilterByPath(t *testing.T) {
	result := filter.Apply(sampleChanges, filter.Options{Paths: []string{"db."}})
	if len(result) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(result))
	}
}

func TestApply_FilterByType(t *testing.T) {
	result := filter.Apply(sampleChanges, filter.Options{Types: []diff.ChangeType{diff.Added}})
	if len(result) != 2 {
		t.Fatalf("expected 2 added changes, got %d", len(result))
	}
}

func TestApply_FilterByPathAndType(t *testing.T) {
	result := filter.Apply(sampleChanges, filter.Options{
		Paths: []string{"server."},
		Types: []diff.ChangeType{diff.Modified},
	})
	if len(result) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result))
	}
	if result[0].Path != "server.port" {
		t.Errorf("unexpected path %s", result[0].Path)
	}
}

func TestApply_NoMatch(t *testing.T) {
	result := filter.Apply(sampleChanges, filter.Options{Paths: []string{"nonexistent."}})
	if len(result) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(result))
	}
}
