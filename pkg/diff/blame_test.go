package diff

import (
	"testing"
	"time"
)

var blameTime = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func TestBlame_Empty(t *testing.T) {
	bm := Blame(nil, "alice", blameTime)
	if len(bm) != 0 {
		t.Fatalf("expected empty BlameMap, got %d entries", len(bm))
	}
}

func TestBlame_SingleChange(t *testing.T) {
	changes := []Change{
		{Path: "server.port", Type: ChangeModified, OldValue: 8080, NewValue: 9090},
	}
	bm := Blame(changes, "bob", blameTime)
	if len(bm) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(bm))
	}
	e := bm["server.port"]
	if e.Author != "bob" {
		t.Errorf("expected author bob, got %s", e.Author)
	}
	if e.ChangeType != string(ChangeModified) {
		t.Errorf("unexpected change type: %s", e.ChangeType)
	}
	if !e.ChangedAt.Equal(blameTime) {
		t.Errorf("unexpected timestamp: %v", e.ChangedAt)
	}
}

func TestBlame_MultipleChanges(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeRemoved},
		{Path: "c", Type: ChangeModified},
	}
	bm := Blame(changes, "carol", blameTime)
	if len(bm) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(bm))
	}
	for _, path := range []string{"a", "b", "c"} {
		if _, ok := bm[path]; !ok {
			t.Errorf("missing blame entry for path %q", path)
		}
	}
}

func TestBlame_DuplicatePath_LastWins(t *testing.T) {
	changes := []Change{
		{Path: "x", Type: ChangeAdded},
		{Path: "x", Type: ChangeModified},
	}
	bm := Blame(changes, "dave", blameTime)
	if len(bm) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(bm))
	}
	if bm["x"].ChangeType != string(ChangeModified) {
		t.Errorf("expected last-write to win, got %s", bm["x"].ChangeType)
	}
}

func TestMergeBlame_DisjointMaps(t *testing.T) {
	a := BlameMap{"p1": {Path: "p1", Author: "alice"}}
	b := BlameMap{"p2": {Path: "p2", Author: "bob"}}
	out := MergeBlame(a, b)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
}

func TestMergeBlame_SrcOverridesDst(t *testing.T) {
	a := BlameMap{"p": {Path: "p", Author: "alice"}}
	b := BlameMap{"p": {Path: "p", Author: "bob"}}
	out := MergeBlame(a, b)
	if out["p"].Author != "bob" {
		t.Errorf("expected src to override dst, got %s", out["p"].Author)
	}
}
