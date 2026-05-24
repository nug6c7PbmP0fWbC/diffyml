package formatter

import (
	"fmt"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// NewSliceFormatter returns a Formatter that renders a sliced window of
// changes with a header showing the displayed range.
func NewSliceFormatter(opts diff.SliceOptions) Formatter {
	return func(changes []diff.Change) (string, error) {
		var sb strings.Builder

		sliced := diff.Slice(changes, opts)

		end := opts.End
		if end == 0 {
			end = len(changes)
		}
		if end > len(changes) {
			end = len(changes)
		}

		fmt.Fprintf(&sb, "Slice [%d:%d] — %d change(s)\n", opts.Start, end, len(sliced))
		sb.WriteString(strings.Repeat("-", 40) + "\n")

		if len(sliced) == 0 {
			sb.WriteString("(no changes in range)\n")
			return sb.String(), nil
		}

		for i, c := range sliced {
			symbol := changeSymbolSlice(c.Type)
			fmt.Fprintf(&sb, "%3d. %s %s", opts.Start+i, symbol, c.Path)
			switch c.Type {
			case diff.ChangeAdded:
				fmt.Fprintf(&sb, " = %v", c.After)
			case diff.ChangeRemoved:
				fmt.Fprintf(&sb, " (was %v)", c.Before)
			case diff.ChangeModified:
				fmt.Fprintf(&sb, ": %v → %v", c.Before, c.After)
			}
			sb.WriteByte('\n')
		}

		return sb.String(), nil
	}
}

func changeSymbolSlice(ct diff.ChangeType) string {
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
