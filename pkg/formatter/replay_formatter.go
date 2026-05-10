package formatter

import (
	"fmt"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// ReplayFormatter renders the results of a diff.Replay call as human-readable
// text, showing each step index, its changes, and any error.
type ReplayFormatter struct{}

// NewReplayFormatter returns a new ReplayFormatter.
func NewReplayFormatter() *ReplayFormatter {
	return &ReplayFormatter{}
}

// Format writes a replay report to a strings.Builder and returns the string.
func (f *ReplayFormatter) Format(results []diff.ReplayResult) string {
	if len(results) == 0 {
		return "replay: no steps\n"
	}

	var sb strings.Builder
	sb.WriteString("=== Replay Report ===\n")

	for _, r := range results {
		if r.Err != nil {
			fmt.Fprintf(&sb, "Step %d [ERROR]: %v\n", r.Step, r.Err)
			continue
		}
		fmt.Fprintf(&sb, "Step %d [%d change(s)]:\n", r.Step, len(r.Changes))
		for _, c := range r.Changes {
			sb.WriteString("  ")
			sb.WriteString(formatReplayChange(c))
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func formatReplayChange(c diff.Change) string {
	switch c.Type {
	case diff.ChangeAdded:
		return fmt.Sprintf("+ %s = %v", c.Path, c.Value)
	case diff.ChangeRemoved:
		return fmt.Sprintf("- %s (was %v)", c.Path, c.OldValue)
	case diff.ChangeModified:
		return fmt.Sprintf("~ %s: %v -> %v", c.Path, c.OldValue, c.Value)
	default:
		return fmt.Sprintf("? %s", c.Path)
	}
}
