package diff

import (
	"testing"
)

func TestCluster_Empty(t *testing.T) {
	out := Cluster(nil, DefaultClusterOptions())
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %v", out)
	}
}

func TestCluster_SingleChange(t *testing.T) {
	changes := []Change{
		{Path: "app.name", Type: ChangeModified},
	}
	out := Cluster(changes, ClusterOptions{MaxDistance: 2, MinSize: 1})
	if len(out) != 1 {
		t.Fatalf("expected 1 cluster, got %d", len(out))
	}
}

func TestCluster_MinSizeFiltersSmallClusters(t *testing.T) {
	changes := []Change{
		{Path: "app.name", Type: ChangeModified},
		{Path: "app.port", Type: ChangeModified},
		{Path: "zzz.unrelated", Type: ChangeAdded},
	}
	out := Cluster(changes, ClusterOptions{MaxDistance: 5, MinSize: 2})
	// "zzz.unrelated" may be alone in its cluster and thus removed.
	for k, v := range out {
		if len(v) < 2 {
			t.Errorf("cluster %q has only %d change(s), expected >= 2", k, len(v))
		}
	}
}

func TestCluster_GroupsRelatedPaths(t *testing.T) {
	changes := []Change{
		{Path: "db.host", Type: ChangeModified},
		{Path: "db.port", Type: ChangeModified},
		{Path: "db.name", Type: ChangeAdded},
	}
	out := Cluster(changes, ClusterOptions{MaxDistance: 3, MinSize: 1})
	total := 0
	for _, v := range out {
		total += len(v)
	}
	if total != len(changes) {
		t.Errorf("expected %d total changes across clusters, got %d", len(changes), total)
	}
}

func TestEditDistance_Equal(t *testing.T) {
	if d := editDistance("abc", "abc"); d != 0 {
		t.Errorf("expected 0, got %d", d)
	}
}

func TestEditDistance_Insert(t *testing.T) {
	if d := editDistance("abc", "abcd"); d != 1 {
		t.Errorf("expected 1, got %d", d)
	}
}

func TestEditDistance_Replace(t *testing.T) {
	if d := editDistance("abc", "axc"); d != 1 {
		t.Errorf("expected 1, got %d", d)
	}
}

func TestEditDistance_EmptyStrings(t *testing.T) {
	if d := editDistance("", ""); d != 0 {
		t.Errorf("expected 0, got %d", d)
	}
	if d := editDistance("", "abc"); d != 3 {
		t.Errorf("expected 3, got %d", d)
	}
}
