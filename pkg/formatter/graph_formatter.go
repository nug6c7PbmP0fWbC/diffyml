package formatter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// GraphFormatter renders a change dependency graph as an ASCII tree.
type GraphFormatter struct{}

// NewGraphFormatter returns a new GraphFormatter.
func NewGraphFormatter() *GraphFormatter {
	return &GraphFormatter{}
}

// Format renders the graph built from changes as an indented tree.
func (f *GraphFormatter) Format(changes []diff.Change) string {
	if len(changes) == 0 {
		return "(no changes)\n"
	}

	g := diff.BuildGraph(changes)
	roots := g.Roots()
	sort.Strings(roots)

	var sb strings.Builder
	sb.WriteString("Change Graph:\n")

	visited := make(map[string]bool)
	for _, root := range roots {
		writeGraphNode(&sb, g, root, "", visited)
	}
	return sb.String()
}

func writeGraphNode(sb *strings.Builder, g *diff.Graph, node, indent string, visited map[string]bool) {
	if visited[node] {
		return
	}
	visited[node] = true

	fmt.Fprintf(sb, "%s- %s\n", indent, node)

	children := make([]string, 0)
	for _, neighbor := range g.Edges[node] {
		if len(neighbor) > len(node) && !visited[neighbor] {
			children = append(children, neighbor)
		}
	}
	sort.Strings(children)
	for _, child := range children {
		writeGraphNode(sb, g, child, indent+"  ", visited)
	}
}
