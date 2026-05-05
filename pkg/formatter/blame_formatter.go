package formatter

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// BlameFormatter renders a diff.BlameMap as a human-readable table.
type BlameFormatter struct{}

// NewBlameFormatter returns a new BlameFormatter.
func NewBlameFormatter() *BlameFormatter {
	return &BlameFormatter{}
}

// Format writes the blame map to w as a tab-separated table.
func (f *BlameFormatter) Format(w io.Writer, bm diff.BlameMap) error {
	if len(bm) == 0 {
		_, err := fmt.Fprintln(w, "(no blame entries)")
		return err
	}

	// Collect and sort paths for deterministic output.
	paths := make([]string, 0, len(bm))
	for p := range bm {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PATH\tTYPE\tAUTHOR\tCHANGED AT")
	fmt.Fprintln(tw, "----\t----\t------\t----------")
	for _, p := range paths {
		e := bm[p]
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			e.Path,
			e.ChangeType,
			e.Author,
			e.ChangedAt.Format("2006-01-02 15:04:05 UTC"),
		)
	}
	return tw.Flush()
}
