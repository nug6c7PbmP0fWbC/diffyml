package diff

// RollupOptions controls how changes are aggregated into rollup summaries.
type RollupOptions struct {
	// Depth is the path depth at which to group changes (1 = top-level key).
	Depth int
	// IncludeCount adds the count of changes per bucket.
	IncludeCount bool
}

// DefaultRollupOptions returns sensible defaults for Rollup.
func DefaultRollupOptions() RollupOptions {
	return RollupOptions{
		Depth:        1,
		IncludeCount: true,
	}
}

// RollupBucket represents an aggregated group of changes sharing a path prefix.
type RollupBucket struct {
	Prefix  string
	Added   int
	Removed int
	Modified int
	Total   int
}

// Rollup aggregates a slice of Changes into RollupBuckets grouped by path prefix
// up to the given depth.
func Rollup(changes []Change, opts RollupOptions) []RollupBucket {
	if opts.Depth < 1 {
		opts.Depth = 1
	}

	index := map[string]*RollupBucket{}
	order := []string{}

	for _, c := range changes {
		prefix := rollupPrefix(c.Path, opts.Depth)
		if _, ok := index[prefix]; !ok {
			index[prefix] = &RollupBucket{Prefix: prefix}
			order = append(order, prefix)
		}
		b := index[prefix]
		switch c.Type {
		case ChangeAdded:
			b.Added++
		case ChangeRemoved:
			b.Removed++
		case ChangeModified:
			b.Modified++
		}
		b.Total++
	}

	result := make([]RollupBucket, 0, len(order))
	for _, k := range order {
		result = append(result, *index[k])
	}
	return result
}

// rollupPrefix returns the first `depth` segments of a dot-separated path.
func rollupPrefix(path string, depth int) string {
	if path == "" {
		return "(root)"
	}
	count := 0
	for i, ch := range path {
		if ch == '.' {
			count++
			if count == depth {
				return path[:i]
			}
		}
	}
	return path
}
