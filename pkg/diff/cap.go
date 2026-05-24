package diff

// CapOptions controls how Cap trims a change list.
type CapOptions struct {
	// MaxChanges is the hard upper bound on the number of changes returned.
	// A value of 0 means no cap is applied.
	MaxChanges int

	// Priority defines which change types are kept first when trimming.
	// Types listed earlier are preserved over later ones.
	// If empty, the original order is used.
	Priority []ChangeType
}

// DefaultCapOptions returns a CapOptions with no cap applied.
func DefaultCapOptions() CapOptions {
	return CapOptions{
		MaxChanges: 0,
		Priority:   nil,
	}
}

// Cap limits the number of changes to opts.MaxChanges.
// When Priority is set, changes whose type appears earlier in the priority
// slice are retained before lower-priority types. Within the same priority
// tier the original order is preserved.
// If MaxChanges is 0 the original slice is returned unchanged.
func Cap(changes []Change, opts CapOptions) []Change {
	if opts.MaxChanges <= 0 || len(changes) <= opts.MaxChanges {
		return changes
	}

	if len(opts.Priority) == 0 {
		return changes[:opts.MaxChanges]
	}

	// Build a rank map: lower rank == higher priority.
	rank := make(map[ChangeType]int, len(opts.Priority))
	for i, ct := range opts.Priority {
		rank[ct] = i
	}
	unranked := len(opts.Priority) // sentinel for types not in Priority

	// Stable bucket sort by priority tier.
	buckets := make([][]Change, unranked+1)
	for _, c := range changes {
		r, ok := rank[c.Type]
		if !ok {
			r = unranked
		}
		buckets[r] = append(buckets[r], c)
	}

	result := make([]Change, 0, opts.MaxChanges)
	for _, bucket := range buckets {
		for _, c := range bucket {
			if len(result) >= opts.MaxChanges {
				return result
			}
			result = append(result, c)
		}
	}
	return result
}
