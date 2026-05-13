package formatter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// IndexFormatter renders a diff.Index as a sorted, human-readable table
// mapping each path to its change type and current value.
type IndexFormatter struct {
	w io.Writer
}

// NewIndexFormatter returns an IndexFormatter that writes to w.
func NewIndexFormatter(w io.Writer) *IndexFormatter {
	return &IndexFormatter{w: w}
}

// Format writes the index contents to the underlying writer.
func (f *IndexFormatter) Format(idx diff.Index) error {
	if len(idx) == 0 {
		_, err := fmt.Fprintln(f.w, "(empty index)")
		return err
	}

	paths := idx.Paths()
	sort.Strings(paths)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-40s %-10s %s\n", "PATH", "TYPE", "VALUE"))
	sb.WriteString(strings.Repeat("-", 70) + "\n")

	for _, p := range paths {
		c, _ := idx.Lookup(p)
		typeLabel := changeTypeLabel(c.Type)
		value := fmt.Sprintf("%v", c.Value)
		if value == "<nil>" || value == "" {
			value = "(none)"
		}
		sb.WriteString(fmt.Sprintf("%-40s %-10s %s\n", p, typeLabel, value))
	}

	_, err := fmt.Fprint(f.w, sb.String())
	return err
}

func changeTypeLabel(t diff.ChangeType) string {
	switch t {
	case diff.ChangeAdded:
		return "added"
	case diff.ChangeRemoved:
		return "removed"
	case diff.ChangeModified:
		return "modified"
	default:
		return "unknown"
	}
}
