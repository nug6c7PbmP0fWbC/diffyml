package diff

// DedupeOptions controls deduplication behaviour.
type DedupeOptions struct {
	// PreferLast keeps the last occurrence of a duplicate path instead of the first.
	PreferLast bool
}

// DefaultDedupeOptions returns sensible defaults: keep the first occurrence.
func DefaultDedupeOptions() DedupeOptions {
	return DedupeOptions{PreferLast: false}
}

// Dedupe removes duplicate Change entries that share the same Path.
// When multiple changes exist for the same path the first (or last, if
// PreferLast is set) entry is retained and the rest are discarded.
//
// Duplicates can arise when changes are produced by multiple diff passes
// (e.g. after a pipeline merge) or when the same key appears in both a
// base and override document.
func Dedupe(changes []Change, opts DedupeOptions) []Change {
	if len(changes) == 0 {
		return changes
	}

	seen := make(map[string]int, len(changes)) // path -> index in result
	result := make([]Change, 0, len(changes))

	for _, c := range changes {
		if idx, exists := seen[c.Path]; exists {
			if opts.PreferLast {
				result[idx] = c
			}
			// PreferFirst: simply skip
			continue
		}
		seen[c.Path] = len(result)
		result = append(result, c)
	}

	return result
}
