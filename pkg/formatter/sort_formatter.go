package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// SortFormatter renders changes after sorting them according to the
// provided SortOptions, delegating actual rendering to an inner Formatter.
type SortFormatter struct {
	inner   Formatter
	opts    diff.SortOptions
	writer  io.Writer
}

// Formatter is the shared interface expected by all formatters in this package.
type Formatter interface {
	Format(changes []diff.Change) error
}

// NewSortFormatter wraps an existing Formatter and sorts changes before
// passing them through. If inner is nil, a plain text representation is used.
func NewSortFormatter(w io.Writer, opts diff.SortOptions) *SortFormatter {
	return &SortFormatter{
		opts:   opts,
		writer: w,
	}
}

// Format sorts changes and writes a plain sorted listing to the writer.
func (f *SortFormatter) Format(changes []diff.Change) error {
	sorted := diff.Sort(changes, f.opts)

	if len(sorted) == 0 {
		_, err := fmt.Fprintln(f.writer, "no changes")
		return err
	}

	var sb strings.Builder
	for _, c := range sorted {
		var symbol string
		switch c.Type {
		case diff.ChangeAdded:
			symbol = "+"
		case diff.ChangeRemoved:
			symbol = "-"
		case diff.ChangeModified:
			symbol = "~"
		default:
			symbol = "?"
		}
		sb.WriteString(fmt.Sprintf("%s %s\n", symbol, c.Path))
	}

	_, err := fmt.Fprint(f.writer, sb.String())
	return err
}
