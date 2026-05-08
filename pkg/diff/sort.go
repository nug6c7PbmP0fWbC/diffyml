package diff

import (
	"sort"
	"strings"
)

// SortOrder defines the ordering direction for changes.
type SortOrder int

const (
	// SortAscending orders changes from A to Z by path.
	SortAscending SortOrder = iota
	// SortDescending orders changes from Z to A by path.
	SortDescending
	// SortByType groups changes by their type (added, removed, modified).
	SortByType
)

// SortOptions controls how changes are sorted.
type SortOptions struct {
	// Order specifies the sort direction or strategy.
	Order SortOrder
	// Stable preserves the relative order of equal elements.
	Stable bool
}

// DefaultSortOptions returns options that sort changes ascending by path.
func DefaultSortOptions() SortOptions {
	return SortOptions{
		Order:  SortAscending,
		Stable: true,
	}
}

// Sort returns a new slice of changes ordered according to opts.
// The original slice is not modified.
func Sort(changes []Change, opts SortOptions) []Change {
	if len(changes) == 0 {
		return changes
	}

	out := make([]Change, len(changes))
	copy(out, changes)

	less := buildLess(out, opts.Order)

	if opts.Stable {
		sort.SliceStable(out, less)
	} else {
		sort.Slice(out, less)
	}

	return out
}

func buildLess(changes []Change, order SortOrder) func(i, j int) bool {
	switch order {
	case SortDescending:
		return func(i, j int) bool {
			return strings.Compare(changes[i].Path, changes[j].Path) > 0
		}
	case SortByType:
		return func(i, j int) bool {
			ti := typeRank(changes[i].Type)
			tj := typeRank(changes[j].Type)
			if ti != tj {
				return ti < tj
			}
			return strings.Compare(changes[i].Path, changes[j].Path) < 0
		}
	default: // SortAscending
		return func(i, j int) bool {
			return strings.Compare(changes[i].Path, changes[j].Path) < 0
		}
	}
}

func typeRank(t ChangeType) int {
	switch t {
	case ChangeAdded:
		return 0
	case ChangeModified:
		return 1
	case ChangeRemoved:
		return 2
	default:
		return 3
	}
}
