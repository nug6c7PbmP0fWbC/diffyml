package diff

// CompactOptions controls how Compact reduces a change list.
type CompactOptions struct {
	// MaxPerPath is the maximum number of changes to retain per unique path.
	// Zero means no limit.
	MaxPerPath int

	// KeepLast, when true, keeps the last N changes per path instead of the first N.
	KeepLast bool
}

// DefaultCompactOptions returns sensible defaults for CompactOptions.
func DefaultCompactOptions() CompactOptions {
	return CompactOptions{
		MaxPerPath: 1,
		KeepLast:   true,
	}
}

// Compact reduces a slice of Changes by collapsing multiple changes for the
// same path into at most MaxPerPath entries. When KeepLast is true the last
// occurrence(s) win; otherwise the first occurrence(s) are kept.
//
// The relative order of paths is preserved based on the first time each path
// is encountered.
func Compact(changes []Change, opts CompactOptions) []Change {
	if len(changes) == 0 {
		return changes
	}

	if opts.MaxPerPath <= 0 {
		return changes
	}

	// Collect changes per path, preserving insertion order of paths.
	order := make([]string, 0, len(changes))
	buckets := make(map[string][]Change, len(changes))

	for _, c := range changes {
		if _, seen := buckets[c.Path]; !seen {
			order = append(order, c.Path)
		}
		buckets[c.Path] = append(buckets[c.Path], c)
	}

	result := make([]Change, 0, len(order))
	for _, path := range order {
		slice := buckets[path]
		if len(slice) <= opts.MaxPerPath {
			result = append(result, slice...)
			continue
		}
		if opts.KeepLast {
			result = append(result, slice[len(slice)-opts.MaxPerPath:]...)
		} else {
			result = append(result, slice[:opts.MaxPerPath]...)
		}
	}
	return result
}
