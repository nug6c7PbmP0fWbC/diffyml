package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// LimitFormatter renders a paginated view of changes, showing offset/total metadata.
type LimitFormatter struct {
	Offset int
	Total  int
}

// NewLimitFormatter creates a LimitFormatter given the offset used and the
// total number of changes before limiting was applied.
func NewLimitFormatter(offset, total int) *LimitFormatter {
	return &LimitFormatter{Offset: offset, Total: total}
}

// Format writes the limited changes with a pagination header to w.
func (f *LimitFormatter) Format(w io.Writer, changes []diff.Change) error {
	showing := len(changes)
	fmt.Fprintf(w, "# Changes (showing %d-%d of %d)\n",
		f.Offset+1,
		f.Offset+showing,
		f.Total,
	)

	if showing == 0 {
		fmt.Fprintln(w, "(no changes in this range)")
		return nil
	}

	for i, c := range changes {
		var line strings.Builder
		line.WriteString(fmt.Sprintf("%3d. ", f.Offset+i+1))
		switch c.Type {
		case diff.ChangeAdded:
			line.WriteString(fmt.Sprintf("[+] %s = %v", c.Path, c.After))
		case diff.ChangeRemoved:
			line.WriteString(fmt.Sprintf("[-] %s (was %v)", c.Path, c.Before))
		case diff.ChangeModified:
			line.WriteString(fmt.Sprintf("[~] %s: %v -> %v", c.Path, c.Before, c.After))
		default:
			line.WriteString(fmt.Sprintf("[?] %s", c.Path))
		}
		fmt.Fprintln(w, line.String())
	}
	return nil
}
