package diff

// DefaultContextOptions returns a ContextOptions with sensible defaults.
func DefaultContextOptions() ContextOptions {
	return ContextOptions{
		Before: 2,
		After:  2,
	}
}

// ContextOptions controls how many surrounding (unchanged) changes are included
// around each real change when calling WithContext.
type ContextOptions struct {
	// Before is the number of preceding changes to include as context.
	Before int
	// After is the number of following changes to include as context.
	After int
}

// WithContext returns a new slice that includes the original changes plus
// neighbouring context entries marked with metadata key "context": true.
// Changes that are already in the input are preserved as-is; only the
// surrounding entries receive the context marker.
func WithContext(changes []Change, opts ContextOptions) []Change {
	if len(changes) == 0 {
		return changes
	}

	included := make([]bool, len(changes))

	for i, c := range changes {
		if c.Type == Added || c.Type == Removed || c.Type == Modified {
			included[i] = true
			for b := 1; b <= opts.Before; b++ {
				if i-b >= 0 {
					included[i-b] = true
				}
			}
			for a := 1; a <= opts.After; a++ {
				if i+a < len(changes) {
					included[i+a] = true
				}
			}
		}
	}

	out := make([]Change, 0, len(changes))
	for i, c := range changes {
		if !included[i] {
			continue
		}
		if c.Type != Added && c.Type != Removed && c.Type != Modified {
			if c.Metadata == nil {
				c.Metadata = map[string]interface{}{}
			}
			c.Metadata["context"] = true
		}
		out = append(out, c)
	}
	return out
}
