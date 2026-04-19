package formatter

import (
	"fmt"
	"io"

	"github.com/diffyml/diffyml/pkg/diff"
)

// Formatter is the interface all output formatters must implement.
type Formatter interface {
	Format(w io.Writer, changes []diff.Change) error
}

// TextFormatter formats changes as human-readable text.
type TextFormatter struct{}

// NewTextFormatter returns a new TextFormatter.
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{}
}

// Format writes text-formatted changes to w.
func (f *TextFormatter) Format(w io.Writer, changes []diff.Change) error {
	if len(changes) == 0 {
		_, err := fmt.Fprintln(w, "No changes detected.")
		return err
	}
	for _, c := range changes {
		line := formatChange(c)
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func formatChange(c diff.Change) string {
	switch c.Type {
	case diff.Added:
		return fmt.Sprintf("+ %s: %v", c.Path, c.NewValue)
	case diff.Removed:
		return fmt.Sprintf("- %s: %v", c.Path, c.OldValue)
	case diff.Modified:
		// Show old -> new for modified fields; clearer than just listing new value
		return fmt.Sprintf("~ %s: %v -> %v", c.Path, c.OldValue, c.NewValue)
	default:
		return fmt.Sprintf("? %s (unknown change type)", c.Path)
	}
}
