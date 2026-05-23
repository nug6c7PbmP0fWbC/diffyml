package formatter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// SplitFormatter renders a bucketed split result produced by diff.Split.
type SplitFormatter struct {
	w io.Writer
}

// NewSplitFormatter returns a SplitFormatter that writes to w.
func NewSplitFormatter(w io.Writer) *SplitFormatter {
	return &SplitFormatter{w: w}
}

// Format writes the split buckets to the underlying writer.
// Buckets are printed in sorted key order for deterministic output.
func (f *SplitFormatter) Format(buckets map[string][]diff.Change) error {
	if len(buckets) == 0 {
		_, err := fmt.Fprintln(f.w, "no changes")
		return err
	}

	keys := make([]string, 0, len(buckets))
	for k := range buckets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		changes := buckets[k]
		sb.WriteString(fmt.Sprintf("[%s] (%d)\n", k, len(changes)))
		for _, c := range changes {
			switch c.Type {
			case diff.ChangeAdded:
				sb.WriteString(fmt.Sprintf("  + %s: %v\n", c.Path, c.Value))
			case diff.ChangeRemoved:
				sb.WriteString(fmt.Sprintf("  - %s: %v\n", c.Path, c.Before))
			case diff.ChangeModified:
				sb.WriteString(fmt.Sprintf("  ~ %s: %v -> %v\n", c.Path, c.Before, c.Value))
			default:
				sb.WriteString(fmt.Sprintf("  ? %s\n", c.Path))
			}
		}
	}

	_, err := fmt.Fprint(f.w, sb.String())
	return err
}
