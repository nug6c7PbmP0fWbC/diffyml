package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/diffyml/diffyml/pkg/diff"
)

type tableFormatter struct {
	w io.Writer
}

// NewTableFormatter creates a formatter that renders changes as an ASCII table.
func NewTableFormatter(w io.Writer) Formatter {
	return &tableFormatter{w: w}
}

func (f *tableFormatter) Format(changes []diff.Change) error {
	if len(changes) == 0 {
		_, err := fmt.Fprintln(f.w, "No changes detected.")
		return err
	}

	const colPath = 40
	const colType = 10
	const colVal = 20

	sep := strings.Repeat("-", colPath+colType+colVal*2+7)
	header := fmt.Sprintf("| %-*s | %-*s | %-*s | %-*s |",
		colPath-2, "PATH",
		colType-2, "TYPE",
		colVal-2, "OLD VALUE",
		colVal-2, "NEW VALUE",
	)

	lines := []string{sep, header, sep}

	for _, c := range changes {
		old := formatTableVal(c.OldValue)
		new := formatTableVal(c.NewValue)
		row := fmt.Sprintf("| %-*s | %-*s | %-*s | %-*s |",
			colPath-2, truncate(c.Path, colPath-2),
			colType-2, string(c.Type),
			colVal-2, truncate(old, colVal-2),
			colVal-2, truncate(new, colVal-2),
		)
		lines = append(lines, row)
	}
	lines = append(lines, sep)

	_, err := fmt.Fprintln(f.w, strings.Join(lines, "\n"))
	return err
}

func formatTableVal(v interface{}) string {
	if v == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", v)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
