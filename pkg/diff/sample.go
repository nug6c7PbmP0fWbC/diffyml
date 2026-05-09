package diff

import "math/rand"

// SampleOptions controls how changes are sampled.
type SampleOptions struct {
	// N is the maximum number of changes to return.
	// If N <= 0, all changes are returned.
	N int

	// Seed is used to initialise the random source for reproducibility.
	// A zero value means a non-deterministic order.
	Seed int64

	// Deterministic forces the use of Seed even when Seed == 0.
	Deterministic bool
}

// DefaultSampleOptions returns conservative defaults that return all changes.
func DefaultSampleOptions() SampleOptions {
	return SampleOptions{N: 0}
}

// Sample returns a random subset of at most opts.N changes.
// When N <= 0 or N >= len(changes) the original slice is returned unchanged.
// The original slice is never mutated.
func Sample(changes []Change, opts SampleOptions) []Change {
	if len(changes) == 0 {
		return changes
	}
	if opts.N <= 0 || opts.N >= len(changes) {
		return changes
	}

	var rng *rand.Rand
	if opts.Deterministic || opts.Seed != 0 {
		//nolint:gosec // weak rand is intentional for sampling
		rng = rand.New(rand.NewSource(opts.Seed))
	} else {
		//nolint:gosec
		rng = rand.New(rand.NewSource(rand.Int63()))
	}

	// Copy indices and shuffle.
	idx := make([]int, len(changes))
	for i := range idx {
		idx[i] = i
	}
	rng.Shuffle(len(idx), func(i, j int) { idx[i], idx[j] = idx[j], idx[i] })

	// Pick the first N indices and preserve original order.
	chosen := make(map[int]struct{}, opts.N)
	for _, i := range idx[:opts.N] {
		chosen[i] = struct{}{}
	}

	out := make([]Change, 0, opts.N)
	for i, c := range changes {
		if _, ok := chosen[i]; ok {
			out = append(out, c)
		}
	}
	return out
}
