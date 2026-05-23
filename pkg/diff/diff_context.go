package diff

// ContextOptions controls how many surrounding unchanged lines
// (or sibling keys) are retained around each change when producing
// a context-aware view of the diff.
type ContextOptions struct {
	// Lines is the number of neighbouring changes to include on each
	// side of a real change. Defaults to 2.
	Lines int
}

// DefaultContextOptions returns a ContextOptions with sensible defaults.
func DefaultContextOptions() ContextOptions {
	return ContextOptions{
		Lines: 2,
	}
}

// WithContext returns a new slice that includes each change together with up
// to opts.Lines neighbouring changes on either side. Duplicates that would
// appear because two windows overlap are collapsed automatically.
//
// The order of the input slice is preserved.
func WithContext(changes []Change, opts ContextOptions) []Change {
	if len(changes) == 0 {
		return []Change{}
	}

	lines := opts.Lines
	if lines < 0 {
		lines = 0
	}

	seen := make(map[int]struct{})
	var result []Change

	for i := range changes {
		start := i - lines
		if start < 0 {
			start = 0
		}
		end := i + lines
		if end >= len(changes) {
			end = len(changes) - 1
		}
		for j := start; j <= end; j++ {
			if _, ok := seen[j]; !ok {
				seen[j] = struct{}{}
				result = append(result, changes[j])
			}
		}
	}

	return result
}
