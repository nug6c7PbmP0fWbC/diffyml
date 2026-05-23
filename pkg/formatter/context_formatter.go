package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// NewContextFormatter returns a Formatter that renders changes with
// surrounding context lines, similar to a unified diff view.
// opts controls how many neighbouring changes are shown around each change.
func NewContextFormatter(opts diff.ContextOptions) Formatter {
	return func(w io.Writer, changes []diff.Change) error {
		if len(changes) == 0 {
			_, err := fmt.Fprintln(w, "(no changes)")
			return err
		}

		contextual := diff.WithContext(changes, opts)

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("--- context diff (±%d lines) ---\n", opts.Lines))

		for _, c := range contextual {
			var symbol string
			switch c.Type {
			case diff.Added:
				symbol = "+"
			case diff.Removed:
				symbol = "-"
			case diff.Modified:
				symbol = "~"
			default:
				symbol = " "
			}

			switch c.Type {
			case diff.Added:
				sb.WriteString(fmt.Sprintf("%s %s: %v\n", symbol, c.Path, c.After))
			case diff.Removed:
				sb.WriteString(fmt.Sprintf("%s %s: %v\n", symbol, c.Path, c.Before))
			case diff.Modified:
				sb.WriteString(fmt.Sprintf("%s %s: %v -> %v\n", symbol, c.Path, c.Before, c.After))
			default:
				sb.WriteString(fmt.Sprintf("%s %s\n", symbol, c.Path))
			}
		}

		_, err := io.WriteString(w, sb.String())
		return err
	}
}
