package formatter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// EnrichFormatter renders changes with their enriched metadata displayed
// alongside each change line.
type EnrichFormatter struct{}

// NewEnrichFormatter returns a new EnrichFormatter.
func NewEnrichFormatter() *EnrichFormatter { return &EnrichFormatter{} }

// Format writes each change with its metadata key=value pairs appended.
func (f *EnrichFormatter) Format(w io.Writer, changes []diff.Change) error {
	if len(changes) == 0 {
		_, err := fmt.Fprintln(w, "no changes")
		return err
	}
	for _, c := range changes {
		sym := changeSymbolEnrich(c.Type)
		line := fmt.Sprintf("%s %s", sym, c.Path)
		if len(c.Metadata) > 0 {
			line += "  [" + buildMetaLabel(c.Metadata) + "]"
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func changeSymbolEnrich(t diff.ChangeType) string {
	switch t {
	case diff.ChangeAdded:
		return "+"
	case diff.ChangeRemoved:
		return "-"
	case diff.ChangeModified:
		return "~"
	default:
		return "?"
	}
}

func buildMetaLabel(meta map[string]string) string {
	keys := make([]string, 0, len(meta))
	for k := range meta {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+meta[k])
	}
	return strings.Join(parts, ", ")
}
