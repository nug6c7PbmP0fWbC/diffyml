package formatter

import (
	"fmt"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// NewContextFormatter returns a Formatter that renders changes produced by
// diff.WithContext, visually distinguishing context lines from real changes.
func NewContextFormatter() Formatter {
	return &contextFormatter{}
}

type contextFormatter struct{}

func (f *contextFormatter) Format(changes []diff.Change) (string, error) {
	if len(changes) == 0 {
		return "(no changes in context window)\n", nil
	}

	var sb strings.Builder
	sb.WriteString("=== context diff ===\n")

	for _, c := range changes {
		isCtx := c.Metadata != nil && c.Metadata["context"] == true
		if isCtx {
			sb.WriteString(fmt.Sprintf("  ~ %s\n", c.Path))
			continue
		}
		switch c.Type {
		case diff.ChangeAdded:
			sb.WriteString(fmt.Sprintf("  + %s: %v\n", c.Path, c.After))
		case diff.ChangeRemoved:
			sb.WriteString(fmt.Sprintf("  - %s: %v\n", c.Path, c.Before))
		case diff.ChangeModified:
			sb.WriteString(fmt.Sprintf("  ~ %s: %v -> %v\n", c.Path, c.Before, c.After))
		default:
			sb.WriteString(fmt.Sprintf("    %s\n", c.Path))
		}
	}

	return sb.String(), nil
}
