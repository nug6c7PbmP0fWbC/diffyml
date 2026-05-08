package diff

import "errors"

// LimitOptions controls how many changes are returned by Limit.
type LimitOptions struct {
	// Max is the maximum number of changes to return.
	// A value of 0 means no limit.
	Max int

	// Offset skips the first N changes before applying the limit.
	Offset int
}

// DefaultLimitOptions returns a LimitOptions with no restrictions.
func DefaultLimitOptions() LimitOptions {
	return LimitOptions{
		Max:    0,
		Offset: 0,
	}
}

// Limit returns a slice of changes respecting the given offset and max count.
// If opts.Max is 0, all changes after the offset are returned.
// Returns an error if Offset or Max is negative.
func Limit(changes []Change, opts LimitOptions) ([]Change, error) {
	if opts.Offset < 0 {
		return nil, errors.New("limit: offset must be non-negative")
	}
	if opts.Max < 0 {
		return nil, errors.New("limit: max must be non-negative")
	}

	if opts.Offset >= len(changes) {
		return []Change{}, nil
	}

	sliced := changes[opts.Offset:]

	if opts.Max == 0 || opts.Max >= len(sliced) {
		result := make([]Change, len(sliced))
		copy(result, sliced)
		return result, nil
	}

	result := make([]Change, opts.Max)
	copy(result, sliced[:opts.Max])
	return result, nil
}
