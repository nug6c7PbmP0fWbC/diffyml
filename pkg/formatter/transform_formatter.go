package formatter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// TransformFormatter renders a before/after view of transformed changes,
// highlighting any paths that were renamed during transformation.
type TransformFormatter struct {
	original []diff.Change
}

// NewTransformFormatter creates a TransformFormatter that retains the original
// changes for comparison against the transformed output.
func NewTransformFormatter(original []diff.Change) *TransformFormatter {
	return &TransformFormatter{original: original}
}

// Format writes a human-readable diff of original vs transformed changes.
func (f *TransformFormatter) Format(w io.Writer, transformed []diff.Change) error {
	fmt.Fprintln(w, "=== Transform Report ===")

	// index originals by position for side-by-side comparison
	origMap := make(map[int]diff.Change, len(f.original))
	for i, c := range f.original {
		origMap[i] = c
	}

	keys := make([]int, 0, len(origMap))
	for k := range origMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, i := range keys {
		var tc diff.Change
		if i < len(transformed) {
			tc = transformed[i]
		} else {
			tc = origMap[i]
		}
		oc := origMap[i]
		pathNote := ""
		if oc.Path != tc.Path {
			pathNote = fmt.Sprintf(" (renamed from %s)", oc.Path)
		}
		valNote := buildValNote(oc, tc)
		fmt.Fprintf(w, "  [%s] %s%s%s\n",
			strings.ToUpper(string(tc.Type)),
			tc.Path,
			pathNote,
			valNote,
		)
	}
	return nil
}

func buildValNote(orig, trans diff.Change) string {
	if fmt.Sprintf("%v", orig.NewValue) != fmt.Sprintf("%v", trans.NewValue) {
		return fmt.Sprintf(" [value: %v -> %v]", orig.NewValue, trans.NewValue)
	}
	if fmt.Sprintf("%v", orig.OldValue) != fmt.Sprintf("%v", trans.OldValue) {
		return fmt.Sprintf(" [old: %v -> %v]", orig.OldValue, trans.OldValue)
	}
	return ""
}
