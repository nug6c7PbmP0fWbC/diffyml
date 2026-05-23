package diff

// SplitOptions controls how changes are split into buckets.
type SplitOptions struct {
	// MaxBuckets is the maximum number of buckets to produce.
	// Zero means no limit.
	MaxBuckets int
	// Predicate is a user-supplied function that returns the bucket key for a
	// given change. When nil, changes are split by ChangeType.
	Predicate func(c Change) string
}

// DefaultSplitOptions returns sensible defaults for SplitOptions.
func DefaultSplitOptions() SplitOptions {
	return SplitOptions{
		MaxBuckets: 0,
		Predicate:  nil,
	}
}

// Split partitions changes into named buckets according to opts.Predicate.
// If Predicate is nil, changes are bucketed by their ChangeType string.
// Bucket order follows first-seen insertion order. When MaxBuckets > 0 any
// change whose bucket would exceed the limit is placed into the last bucket.
func Split(changes []Change, opts SplitOptions) map[string][]Change {
	if opts.Predicate == nil {
		opts.Predicate = func(c Change) string {
			return string(c.Type)
		}
	}

	buckets := make(map[string][]Change)
	order := []string{}

	for _, c := range changes {
		key := opts.Predicate(c)

		if opts.MaxBuckets > 0 && len(order) >= opts.MaxBuckets {
			// Overflow into the last bucket.
			last := order[len(order)-1]
			buckets[last] = append(buckets[last], c)
			continue
		}

		if _, exists := buckets[key]; !exists {
			order = append(order, key)
		}
		buckets[key] = append(buckets[key], c)
	}

	return buckets
}
