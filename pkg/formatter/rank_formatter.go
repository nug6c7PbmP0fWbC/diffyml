package formatter

import (
	"fmt"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// NewRankFormatter returns a Formatter that renders a ranked list of changes.
func NewRankFormatter(ranked []diff.RankedChange) Formatter {
	return &rankFormatter{ranked: ranked}
}

type rankFormatter struct {
	ranked []diff.RankedChange
}

func (f *rankFormatter) Format(changes []diff.Change) (string, error) {
	if len(f.ranked) == 0 {
		return "(no ranked changes)\n", nil
	}

	var sb strings.Builder
	sb.WriteString("Ranked Changes\n")
	sb.WriteString(strings.Repeat("-", 40) + "\n")

	for i, rc := range f.ranked {
		symbol := changeSymbolRank(rc.Change.Type)
		line := fmt.Sprintf("#%d [%.2f] %s %s",
			i+1, rc.Score, symbol, rc.Change.Path)
		if rc.Change.Type == diff.ChangeModified {
			line += fmt.Sprintf(" (%v -> %v)", rc.Change.OldValue, rc.Change.NewValue)
		} else if rc.Change.Type == diff.ChangeAdded {
			line += fmt.Sprintf(" (+ %v)", rc.Change.NewValue)
		} else if rc.Change.Type == diff.ChangeRemoved {
			line += fmt.Sprintf(" (- %v)", rc.Change.OldValue)
		}
		sb.WriteString(line + "\n")
	}
	return sb.String(), nil
}

func changeSymbolRank(ct diff.ChangeType) string {
	switch ct {
	case diff.ChangeAdded:
		return "+"
	case diff.ChangeRemoved:
		return "-"
	case diff.ChangeModified:
		return "~"
	default:
		return "?"
	}
}
