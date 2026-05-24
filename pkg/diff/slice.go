package diff

// SliceOptions controls how Slice trims a change list to a sub-range.
type SliceOptions struct {
	// Start is the inclusive zero-based index of the first change to keep.
	Start int
	// End is the exclusive index. A value of 0 means "until the end".
	End int
}

// DefaultSliceOptions returns a SliceOptions that keeps all changes.
func DefaultSliceOptions() SliceOptions {
	return SliceOptions{Start: 0, End: 0}
}

// Slice returns the sub-slice of changes defined by [Start, End).
// Negative or out-of-bound indices are clamped to valid positions so the
// function never panics.
//
// If End is 0 (the zero value) it is treated as len(changes), i.e. the slice
// extends to the last element.
func Slice(changes []Change, opts SliceOptions) []Change {
	n := len(changes)
	if n == 0 {
		return []Change{}
	}

	start := opts.Start
	if start < 0 {
		start = 0
	}
	if start > n {
		start = n
	}

	end := opts.End
	if end == 0 || end > n {
		end = n
	}
	if end < start {
		end = start
	}

	result := make([]Change, end-start)
	copy(result, changes[start:end])
	return result
}
