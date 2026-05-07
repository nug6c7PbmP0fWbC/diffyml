package diff

import (
	"testing"
)

func TestGroupChanges_Empty(t *testing.T) {
	groups := GroupChanges(nil, GroupOptions{})
	if len(groups) != 0 {
		t.Fatalf("expected 0 groups, got %d", len(groups))
	}
}

func TestGroupChanges_DefaultDepth(t *testing.T) {
	changes := []Change{
		{Path: "server.host", Type: ChangeModified},
		{Path: "server.port", Type: ChangeAdded},
		{Path: "database.url", Type: ChangeModified},
	}
	groups := GroupChanges(changes, GroupOptions{})
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Name != "server" {
		t.Errorf("expected first group 'server', got %q", groups[0].Name)
	}
	if len(groups[0].Changes) != 2 {
		t.Errorf("expected 2 changes in 'server', got %d", len(groups[0].Changes))
	}
	if groups[1].Name != "database" {
		t.Errorf("expected second group 'database', got %q", groups[1].Name)
	}
}

func TestGroupChanges_Depth2(t *testing.T) {
	changes := []Change{
		{Path: "a.b.c", Type: ChangeAdded},
		{Path: "a.b.d", Type: ChangeRemoved},
		{Path: "a.x.y", Type: ChangeModified},
	}
	groups := GroupChanges(changes, GroupOptions{MaxDepth: 2})
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Name != "a.b" {
		t.Errorf("expected group 'a.b', got %q", groups[0].Name)
	}
	if groups[1].Name != "a.x" {
		t.Errorf("expected group 'a.x', got %q", groups[1].Name)
	}
}

func TestGroupChanges_RootPath(t *testing.T) {
	changes := []Change{
		{Path: "", Type: ChangeAdded},
	}
	groups := GroupChanges(changes, GroupOptions{})
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Name != "(root)" {
		t.Errorf("expected group '(root)', got %q", groups[0].Name)
	}
}

func TestGroupChanges_PreservesOrder(t *testing.T) {
	changes := []Change{
		{Path: "z.key", Type: ChangeAdded},
		{Path: "a.key", Type: ChangeModified},
		{Path: "m.key", Type: ChangeRemoved},
	}
	groups := GroupChanges(changes, GroupOptions{})
	names := []string{groups[0].Name, groups[1].Name, groups[2].Name}
	expected := []string{"z", "a", "m"}
	for i, n := range names {
		if n != expected[i] {
			t.Errorf("position %d: expected %q, got %q", i, expected[i], n)
		}
	}
}
