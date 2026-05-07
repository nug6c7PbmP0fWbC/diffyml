package formatter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// ClusterFormatter renders grouped clusters of changes as human-readable text.
type ClusterFormatter struct {
	w io.Writer
}

// NewClusterFormatter creates a ClusterFormatter writing to w.
func NewClusterFormatter(w io.Writer) *ClusterFormatter {
	return &ClusterFormatter{w: w}
}

// Format writes the clustered diff output.
func (f *ClusterFormatter) Format(clusters map[string][]diff.Change) error {
	if len(clusters) == 0 {
		_, err := fmt.Fprintln(f.w, "no clusters")
		return err
	}

	// Sort keys for deterministic output.
	keys := make([]string, 0, len(clusters))
	for k := range clusters {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		changes := clusters[k]
		_, err := fmt.Fprintf(f.w, "cluster [%s] (%d change(s)):\n", k, len(changes))
		if err != nil {
			return err
		}
		for _, c := range changes {
			line := formatClusterChange(c)
			if _, err := fmt.Fprintf(f.w, "  %s\n", line); err != nil {
				return err
			}
		}
	}
	return nil
}

func formatClusterChange(c diff.Change) string {
	switch c.Type {
	case diff.ChangeAdded:
		return fmt.Sprintf("+  %s: %v", c.Path, c.NewValue)
	case diff.ChangeRemoved:
		return fmt.Sprintf("-  %s: %v", c.Path, c.OldValue)
	case diff.ChangeModified:
		return fmt.Sprintf("~  %s: %v -> %v", c.Path, c.OldValue, c.NewValue)
	default:
		return strings.TrimSpace(fmt.Sprintf("?  %s", c.Path))
	}
}
