package diff

// DefaultContextOptions returns a ContextOptions with sensible defaults.
func DefaultContextOptions() ContextOptions {
	return ContextOptions{
		Before: 2,
		After:  2,
	}
}

// ContextOptions controls how many surrounding (unchanged) changes are
// included around each real change.
type ContextOptions struct {
	// Before is the number of preceding changes to include.
	Before int
	// After is the number of following changes to include.
	After int
}

// WithContext returns a new slice that contains every change that is within
// opts.Before or opts.After positions of a "real" (non-context) change.
// Changes that were already present are tagged with metadata key "context"
// set to true so formatters can render them differently.
func WithContext(changes []Change, opts ContextOptions) []Change {
	if len(changes) == 0 {
		return nil
	}

	include := make([]bool, len(changes))

	for i, c := range changes {
		if c.Type == ChangeAdded || c.Type == ChangeRemoved || c.Type == ChangeModified {
			include[i] = true
			for b := 1; b <= opts.Before; b++ {
				if i-b >= 0 {
					include[i-b] = true
				}
			}
			for a := 1; a <= opts.After; a++ {
				if i+a < len(changes) {
					include[i+a] = true
				}
			}
		}
	}

	out := make([]Change, 0, len(changes))
	for i, c := range changes {
		if !include[i] {
			continue
		}
		if c.Type != ChangeAdded && c.Type != ChangeRemoved && c.Type != ChangeModified {
			if c.Metadata == nil {
				c.Metadata = map[string]interface{}{}
			}
			c.Metadata["context"] = true
		}
		out = append(out, c)
	}
	return out
}
