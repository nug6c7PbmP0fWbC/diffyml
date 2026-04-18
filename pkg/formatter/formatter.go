package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/diffyml/diffyml/pkg/diff"
)

// Format represents the output format type.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// TextFormatter writes a human-readable diff to w.
type TextFormatter struct {
	w io.Writer
}

// NewTextFormatter creates a TextFormatter writing to w.
func NewTextFormatter(w io.Writer) *TextFormatter {
	return &TextFormatter{w: w}
}

// Write outputs the list of changes in plain-text format.
func (f *TextFormatter) Write(changes []diff.Change) error {
	if len(changes) == 0 {
		_, err := fmt.Fprintln(f.w, "No differences found.")
		return err
	}
	for _, c := range changes {
		line, err := formatChange(c)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(f.w, line); err != nil {
			return err
		}
	}
	return nil
}

func formatChange(c diff.Change) (string, error) {
	path := strings.Join(c.Path, ".")
	switch c.Type {
	case diff.Added:
		return fmt.Sprintf("+ %s: %v", path, c.To), nil
	case diff.Removed:
		return fmt.Sprintf("- %s: %v", path, c.From), nil
	case diff.Modified:
		return fmt.Sprintf("~ %s: %v -> %v", path, c.From, c.To), nil
	}
	return "", fmt.Errorf("unknown change type: %v", c.Type)
}
