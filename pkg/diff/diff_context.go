package diff

// ContextOptions configures the WithContext function.
type ContextOptions struct {
	// Lines is the number of surrounding (non-changed) neighbours to include
	// on each side of every changed entry.
	Lines int
}

// DefaultContextOptions returns sensible defaults.
func DefaultContextOptions() ContextOptions {
	return ContextOptions{Lines: 2}
}

// WithContext returns a new slice that contains every changed entry from
// changes plus up to opts.Lines unchanged neighbours on either side.
// Neighbours are taken from the ordered all slice which must contain the
// full, ordered set of changes (both changed and unchanged).
//
// If all is nil or empty the function simply returns the original changes
// slice unchanged.
func WithContext(changes []Change, all []Change, opts ContextOptions) []Change {
	if len(all) == 0 || opts.Lines <= 0 {
		return changes
	}

	// Build a set of indices that are already "changed".
	changedIdx := make(map[int]bool)
	for i, c := range all {
		for _, ch := range changes {
			if c.Path == ch.Path && c.Type == ch.Type {
				changedIdx[i] = true
				break
			}
		}
	}

	// Expand the set with neighbours.
	include := make(map[int]bool)
	for idx := range changedIdx {
		for d := -opts.Lines; d <= opts.Lines; d++ {
			n := idx + d
			if n >= 0 && n < len(all) {
				include[n] = true
			}
		}
	}

	// Collect in original order, deduplicating.
	out := make([]Change, 0, len(include))
	for i, c := range all {
		if include[i] {
			out = append(out, c)
		}
	}
	return out
}
