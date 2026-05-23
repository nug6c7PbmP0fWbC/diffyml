package diff

// Bucket represents a named collection of changes.
type Bucket struct {
	Name    string
	Changes []Change
}

// BucketOptions controls how changes are bucketed.
type BucketOptions struct {
	// KeyFunc assigns a bucket name to each change.
	// If nil, changes are bucketed by their ChangeType string.
	KeyFunc func(c Change) string
	// PreserveOrder keeps changes in their original order within each bucket.
	PreserveOrder bool
}

// DefaultBucketOptions returns sensible defaults.
func DefaultBucketOptions() BucketOptions {
	return BucketOptions{
		KeyFunc:       nil,
		PreserveOrder: true,
	}
}

// BucketChanges groups changes into named buckets according to opts.KeyFunc.
// If KeyFunc is nil, changes are grouped by ChangeType (added/removed/modified).
// The returned slice preserves bucket insertion order.
func BucketChanges(changes []Change, opts BucketOptions) []Bucket {
	if opts.KeyFunc == nil {
		opts.KeyFunc = func(c Change) string {
			return string(c.Type)
		}
	}

	indexMap := make(map[string]int)
	var buckets []Bucket

	for _, c := range changes {
		key := opts.KeyFunc(c)
		idx, exists := indexMap[key]
		if !exists {
			idx = len(buckets)
			indexMap[key] = idx
			buckets = append(buckets, Bucket{Name: key})
		}
		buckets[idx].Changes = append(buckets[idx].Changes, c)
	}

	return buckets
}
