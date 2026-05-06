package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// LintFormatter formats lint violations as human-readable text.
type LintFormatter struct {
	w io.Writer
}

// NewLintFormatter creates a LintFormatter that writes to w.
func NewLintFormatter(w io.Writer) *LintFormatter {
	return &LintFormatter{w: w}
}

// Format writes lint violations to the underlying writer.
// If there are no violations it prints a clean-pass message.
func (f *LintFormatter) Format(violations []diff.LintViolation) error {
	if len(violations) == 0 {
		_, err := fmt.Fprintln(f.w, "lint: no violations found")
		return err
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("lint: %d violation(s) found\n", len(violations)))
	for _, v := range violations {
		sb.WriteString(fmt.Sprintf("  %-20s %s  %s\n", "["+v.Rule+"]", v.Path, v.Detail))
	}
	_, err := fmt.Fprint(f.w, sb.String())
	return err
}
