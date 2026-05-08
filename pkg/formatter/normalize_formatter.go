package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// NormalizeFormatter formats changes after normalization, showing which
// changes were retained and summarising the normalization applied.
type NormalizeFormatter struct {
	opts diff.NormalizeOptions
}

// NewNormalizeFormatter creates a formatter that normalizes changes using
// opts before rendering them as human-readable text.
func NewNormalizeFormatter(opts diff.NormalizeOptions) *NormalizeFormatter {
	return &NormalizeFormatter{opts: opts}
}

// Format normalizes changes and writes a text report to w.
func (f *NormalizeFormatter) Format(w io.Writer, changes []diff.Change) error {
	normalized := diff.Normalize(changes, f.opts)

	var sb strings.Builder
	sb.WriteString("# Normalized Diff\n")
	fmt.Fprintf(&sb, "# original: %d  after normalization: %d\n", len(changes), len(normalized))

	if len(normalized) == 0 {
		sb.WriteString("no changes\n")
		_, err := io.WriteString(w, sb.String())
		return err
	}

	for _, c := range normalized {
		switch c.Type {
		case diff.ChangeAdded:
			fmt.Fprintf(&sb, "+ [%s] = %v\n", c.Path, c.NewValue)
		case diff.ChangeRemoved:
			fmt.Fprintf(&sb, "- [%s] = %v\n", c.Path, c.OldValue)
		case diff.ChangeModified:
			fmt.Fprintf(&sb, "~ [%s]: %v -> %v\n", c.Path, c.OldValue, c.NewValue)
		default:
			fmt.Fprintf(&sb, "? [%s]\n", c.Path)
		}
	}

	_, err := io.WriteString(w, sb.String())
	return err
}
