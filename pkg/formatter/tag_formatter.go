package formatter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// TagFormatter renders TaggedChange slices as human-readable text,
// grouping changes by their tag key=value pairs.
type TagFormatter struct{}

// NewTagFormatter returns a new TagFormatter.
func NewTagFormatter() *TagFormatter {
	return &TagFormatter{}
}

// Format writes tagged changes to w, grouped by tags.
func (f *TagFormatter) Format(w io.Writer, tagged []diff.TaggedChange) error {
	if len(tagged) == 0 {
		_, err := fmt.Fprintln(w, "No tagged changes.")
		return err
	}

	// Group by tag signature for display.
	type group struct {
		label   string
		changes []diff.TaggedChange
	}
	index := map[string]*group{}
	order := []string{}

	for _, tc := range tagged {
		label := buildTagLabel(tc.Tags)
		if _, ok := index[label]; !ok {
			index[label] = &group{label: label}
			order = append(order, label)
		}
		index[label].changes = append(index[label].changes, tc)
	}

	sort.Strings(order)

	for _, label := range order {
		g := index[label]
		if label == "" {
			fmt.Fprintln(w, "[untagged]")
		} else {
			fmt.Fprintf(w, "[%s]\n", label)
		}
		for _, tc := range g.changes {
			fmt.Fprintf(w, "  %s %s\n", tc.Type, tc.Path)
		}
	}
	return nil
}

func buildTagLabel(tags []diff.Tag) string {
	if len(tags) == 0 {
		return ""
	}
	parts := make([]string, len(tags))
	for i, t := range tags {
		parts[i] = t.Key + "=" + t.Value
	}
	sort.Strings(parts)
	return strings.Join(parts, ", ")
}
