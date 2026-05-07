package diff

// Graph represents a dependency graph of changes, where each node
// is a change path and edges represent paths that share a common prefix.
type Graph struct {
	Nodes []string
	Edges map[string][]string
}

// BuildGraph constructs a change dependency graph from a slice of Changes.
// Two changes are considered connected if one path is a prefix of the other.
func BuildGraph(changes []Change) *Graph {
	g := &Graph{
		Edges: make(map[string][]string),
	}

	paths := make([]string, 0, len(changes))
	for _, c := range changes {
		paths = append(paths, c.Path)
		g.Nodes = append(g.Nodes, c.Path)
		if _, ok := g.Edges[c.Path]; !ok {
			g.Edges[c.Path] = []string{}
		}
	}

	for i := 0; i < len(paths); i++ {
		for j := i + 1; j < len(paths); j++ {
			a, b := paths[i], paths[j]
			if isGraphPrefix(a, b) || isGraphPrefix(b, a) {
				g.Edges[a] = append(g.Edges[a], b)
				g.Edges[b] = append(g.Edges[b], a)
			}
		}
	}

	return g
}

// Roots returns all nodes that have no incoming edges from shorter paths.
func (g *Graph) Roots() []string {
	child := make(map[string]bool)
	for src, dsts := range g.Edges {
		for _, dst := range dsts {
			if len(dst) > len(src) {
				child[dst] = true
			}
		}
	}
	var roots []string
	for _, n := range g.Nodes {
		if !child[n] {
			roots = append(roots, n)
		}
	}
	return roots
}

func isGraphPrefix(prefix, path string) bool {
	if prefix == path {
		return false
	}
	if len(prefix) >= len(path) {
		return false
	}
	return path[:len(prefix)] == prefix && path[len(prefix)] == '.'
}
