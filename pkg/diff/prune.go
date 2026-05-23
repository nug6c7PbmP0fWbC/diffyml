package diff

// PruneOptions controls which changes are removed from a changeset
// based on their value characteristics.
type PruneOptions struct {
	// RemoveNilValues drops changes where the new value is nil.
	RemoveNilValues bool
	// RemoveZeroValues drops changes where the new value is a zero value
	// (empty string, 0, false).
	RemoveZeroValues bool
	// RemoveUnchangedKeys drops changes where old and new values are equal
	// (defensive guard for identity changes that slipped through).
	RemoveUnchangedKeys bool
	// Paths is an optional allowlist of path prefixes to prune.
	// When empty, all paths are eligible.
	Paths []string
}

// DefaultPruneOptions returns a PruneOptions with sensible defaults.
func DefaultPruneOptions() PruneOptions {
	return PruneOptions{
		RemoveNilValues:     true,
		RemoveZeroValues:    false,
		RemoveUnchangedKeys: true,
	}
}

// Prune removes changes from the slice that match the configured prune
// criteria. It never mutates the input slice.
func Prune(changes []Change, opts PruneOptions) []Change {
	out := make([]Change, 0, len(changes))
	for _, c := range changes {
		if !shouldPrune(c, opts) {
			out = append(out, c)
		}
	}
	return out
}

func shouldPrune(c Change, opts PruneOptions) bool {
	if len(opts.Paths) > 0 && !matchesPrunePath(c.Path, opts.Paths) {
		return false
	}
	if opts.RemoveUnchangedKeys && c.Type == Modified && c.Before == c.After {
		return true
	}
	if opts.RemoveNilValues && c.After == nil {
		return true
	}
	if opts.RemoveZeroValues && isZero(c.After) {
		return true
	}
	return false
}

func matchesPrunePath(path string, prefixes []string) bool {
	for _, p := range prefixes {
		if len(path) >= len(p) && path[:len(p)] == p {
			return true
		}
	}
	return false
}

func isZero(v interface{}) bool {
	if v == nil {
		return true
	}
	switch val := v.(type) {
	case string:
		return val == ""
	case int:
		return val == 0
	case float64:
		return val == 0
	case bool:
		return !val
	}
	return false
}
