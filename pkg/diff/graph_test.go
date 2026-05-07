package diff

import (
	"sort"
	"testing"
)

func TestBuildGraph_Empty(t *testing.T) {
	g := BuildGraph(nil)
	if len(g.Nodes) != 0 {
		t.Fatalf("expected 0 nodes, got %d", len(g.Nodes))
	}
}

func TestBuildGraph_NoEdges(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeAdded},
	}
	g := BuildGraph(changes)
	if len(g.Nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(g.Nodes))
	}
	if len(g.Edges["a"]) != 0 || len(g.Edges["b"]) != 0 {
		t.Fatal("expected no edges between unrelated paths")
	}
}

func TestBuildGraph_WithEdges(t *testing.T) {
	changes := []Change{
		{Path: "server", Type: ChangeModified},
		{Path: "server.port", Type: ChangeModified},
		{Path: "server.host", Type: ChangeAdded},
	}
	g := BuildGraph(changes)
	if len(g.Nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(g.Nodes))
	}
	if len(g.Edges["server"]) != 2 {
		t.Fatalf("expected 2 edges from 'server', got %d", len(g.Edges["server"]))
	}
}

func TestBuildGraph_Roots(t *testing.T) {
	changes := []Change{
		{Path: "server", Type: ChangeModified},
		{Path: "server.port", Type: ChangeModified},
		{Path: "db", Type: ChangeAdded},
	}
	g := BuildGraph(changes)
	roots := g.Roots()
	sort.Strings(roots)
	if len(roots) != 2 {
		t.Fatalf("expected 2 roots, got %d: %v", len(roots), roots)
	}
	if roots[0] != "db" || roots[1] != "server" {
		t.Errorf("unexpected roots: %v", roots)
	}
}

func TestIsGraphPrefix(t *testing.T) {
	if !isGraphPrefix("server", "server.port") {
		t.Error("expected server to be prefix of server.port")
	}
	if isGraphPrefix("server", "server") {
		t.Error("identical paths should not be prefix")
	}
	if isGraphPrefix("serv", "server.port") {
		t.Error("partial word match should not count as prefix")
	}
}
