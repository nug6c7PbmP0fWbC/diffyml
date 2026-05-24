package formatter

import (
	"fmt"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// NewContextFormatter returns a Formatter that renders a context-aware diff,
// visually distinguishing real changes from surrounding context lines.
func NewContextFormatter() Formatter {
	return FormatterFunc(func(changes []diff.Change) (string, error) {
		if len(changes) == 0 {
			return "(no changes)\n", nil
		}

		var sb strings.Builder
		sb.WriteString("--- context diff ---\n")

		prevWasContext := false
		for _, c := range changes {
			isCtx := c.Metadata != nil && c.Metadata["context"] == true

			if isCtx && !prevWasContext {
				// no separator needed at start
			} else if !isCtx && prevWasContext {
				// transitioning from context back to real change — no separator
			}

			if isCtx {
				sb.WriteString(fmt.Sprintf("  %s\n", c.Path))
			} else {
				switch c.Type {
				case diff.Added:
					sb.WriteString(fmt.Sprintf("+ %s: %v\n", c.Path, c.Value))
				case diff.Removed:
					sb.WriteString(fmt.Sprintf("- %s: %v\n", c.Path, c.Before))
				case diff.Modified:
					sb.WriteString(fmt.Sprintf("~ %s: %v -> %v\n", c.Path, c.Before, c.Value))
				default:
					sb.WriteString(fmt.Sprintf("  %s\n", c.Path))
				}
			}
			prevWasContext = isCtx
		}

		return sb.String(), nil
	})
}
