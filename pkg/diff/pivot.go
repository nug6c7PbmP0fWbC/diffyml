package diff

// PivotOptions controls how changes are pivoted (transposed) into a
// path-keyed map for quick lookup.
type PivotOptions struct {
	// KeyFunc builds the map key from a Change. Defaults to using the Path field.
	KeyFunc func(c Change) string
}

// DefaultPivotOptions returns sensible defaults for Pivot.
func DefaultPivotOptions() PivotOptions {
	return PivotOptions{
		KeyFunc: func(c Change) string { return c.Path },
	}
}

// PivotResult holds the pivoted representation of a change set.
type PivotResult struct {
	// Index maps each key (typically the path) to the list of changes
	// that share that key. Multiple changes may exist under the same key
	// when duplicates are present.
	Index map[string][]Change
	// Keys preserves insertion order so iteration is deterministic.
	Keys []string
}

// Pivot converts a flat slice of Changes into a PivotResult keyed by path
// (or by a custom KeyFunc). This is useful for O(1) lookups by path and for
// grouping duplicate paths together.
func Pivot(changes []Change, opts PivotOptions) PivotResult {
	if opts.KeyFunc == nil {
		opts.KeyFunc = DefaultPivotOptions().KeyFunc
	}

	result := PivotResult{
		Index: make(map[string][]Change, len(changes)),
	}

	for _, c := range changes {
		key := opts.KeyFunc(c)
		if _, exists := result.Index[key]; !exists {
			result.Keys = append(result.Keys, key)
		}
		result.Index[key] = append(result.Index[key], c)
	}

	return result
}

// Lookup returns the changes associated with the given key, or nil if absent.
func (p PivotResult) Lookup(key string) []Change {
	return p.Index[key]
}

// Has reports whether key exists in the pivot index.
func (p PivotResult) Has(key string) bool {
	_, ok := p.Index[key]
	return ok
}
