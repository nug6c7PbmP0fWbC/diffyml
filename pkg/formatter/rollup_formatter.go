package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// RollupFormatter renders RollupBuckets as a human-readable table.
type RollupFormatter struct {
	w io.Writer
}

// NewRollupFormatter creates a RollupFormatter writing to w.
func NewRollupFormatter(w io.Writer) *RollupFormatter {
	return &RollupFormatter{w: w}
}

// Format writes a rollup summary table to the underlying writer.
func (f *RollupFormatter) Format(buckets []diff.RollupBucket) error {
	if len(buckets) == 0 {
		_, err := fmt.Fprintln(f.w, "No changes to roll up.")
		return err
	}

	// Header
	_, err := fmt.Fprintf(f.w, "%-40s  %6s  %6s  %6s  %6s\n",
		"PREFIX", "ADDED", "REMOVED", "MODIFIED", "TOTAL")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(f.w, strings.Repeat("-", 72))
	if err != nil {
		return err
	}

	for _, b := range buckets {
		_, err = fmt.Fprintf(f.w, "%-40s  %6d  %6d  %6d  %6d\n",
			b.Prefix, b.Added, b.Removed, b.Modified, b.Total)
		if err != nil {
			return err
		}
	}
	return nil
}
