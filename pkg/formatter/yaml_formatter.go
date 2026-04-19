package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/diffyml/diffyml/pkg/diff"
)

// YAMLFormatter writes a YAML-style diff report.
type YAMLFormatter struct {
	w io.Writer
}

// NewYAMLFormatter creates a YAMLFormatter that writes to w.
func NewYAMLFormatter(w io.Writer) *YAMLFormatter {
	return &YAMLFormatter{w: w}
}

// Format writes the changes as a YAML document.
// Note: values are not quoted, so complex strings with colons or special
// characters may produce invalid YAML. Good enough for simple diffs.
func (f *YAMLFormatter) Format(changes []diff.Change) error {
	if len(changes) == 0 {
		_, err := fmt.Fprintln(f.w, "changes: []")
		return err
	}
	lines := []string{"changes:"}
	for _, c := range changes {
		lines = append(lines, fmt.Sprintf("  - path: %s", c.Path))
		lines = append(lines, fmt.Sprintf("    type: %s", c.Type))
		switch c.Type {
		case diff.Added:
			lines = append(lines, fmt.Sprintf("    new_value: %v", c.NewValue))
		case diff.Removed:
			lines = append(lines, fmt.Sprintf("    old_value: %v", c.OldValue))
		case diff.Modified:
			lines = append(lines, fmt.Sprintf("    old_value: %v", c.OldValue))
			lines = append(lines, fmt.Sprintf("    new_value: %v", c.NewValue))
		}
	}
	_, err := fmt.Fprintln(f.w, strings.Join(lines, "\n"))
	return err
}
