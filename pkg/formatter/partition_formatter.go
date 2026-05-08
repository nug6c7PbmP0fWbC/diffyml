package formatter

import (
	"fmt"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// PartitionFormatter renders partitioned change sets with section headers.
type PartitionFormatter struct {
	opts diff.PartitionOptions
}

// NewPartitionFormatter creates a PartitionFormatter with the given options.
func NewPartitionFormatter(opts diff.PartitionOptions) *PartitionFormatter {
	return &PartitionFormatter{opts: opts}
}

// Format partitions the changes and returns a human-readable string with
// numbered sections, one per partition.
func (f *PartitionFormatter) Format(changes []diff.Change) string {
	if len(changes) == 0 {
		return "No changes.\n"
	}

	partitions := diff.Partition(changes, f.opts)
	var sb strings.Builder

	for i, partition := range partitions {
		fmt.Fprintf(&sb, "=== Partition %d (%d change(s)) ===\n", i+1, len(partition))
		for _, c := range partition {
			symbol := changeSymbol(c.Type)
			switch c.Type {
			case diff.Added:
				fmt.Fprintf(&sb, "  %s %s: %v\n", symbol, c.Path, c.After)
			case diff.Removed:
				fmt.Fprintf(&sb, "  %s %s: %v\n", symbol, c.Path, c.Before)
			case diff.Modified:
				fmt.Fprintf(&sb, "  %s %s: %v -> %v\n", symbol, c.Path, c.Before, c.After)
			}
		}
	}
	return sb.String()
}

func changeSymbol(ct diff.ChangeType) string {
	switch ct {
	case diff.Added:
		return "+"
	case diff.Removed:
		return "-"
	case diff.Modified:
		return "~"
	}
	return "?"
}
