package diff

// Score represents a numeric similarity score between two YAML documents.
type Score struct {
	// Total number of keys considered across both documents.
	Total int
	// Matching keys (unchanged).
	Matching int
	// Similarity is a value between 0.0 and 1.0.
	Similarity float64
}

// ScoreOptions controls how scoring is computed.
type ScoreOptions struct {
	// Weight for added keys (default 1.0).
	AddWeight float64
	// Weight for removed keys (default 1.0).
	RemoveWeight float64
	// Weight for modified keys (default 0.5 — partial credit).
	ModifyWeight float64
}

// DefaultScoreOptions returns sensible defaults.
func DefaultScoreOptions() ScoreOptions {
	return ScoreOptions{
		AddWeight:    1.0,
		RemoveWeight: 1.0,
		ModifyWeight: 0.5,
	}
}

// ComputeScore calculates a similarity score given a slice of changes and options.
// A score of 1.0 means the documents are identical; 0.0 means completely different.
func ComputeScore(changes []Change, opts ScoreOptions) Score {
	if len(changes) == 0 {
		return Score{Total: 0, Matching: 0, Similarity: 1.0}
	}

	var penalty float64
	total := float64(len(changes))

	for _, c := range changes {
		switch c.Type {
		case ChangeAdded:
			penalty += opts.AddWeight
		case ChangeRemoved:
			penalty += opts.RemoveWeight
		case ChangeModified:
			penalty += opts.ModifyWeight
		}
	}

	// Normalise: penalty is relative to total weighted changes.
	maxPenalty := total * max3(opts.AddWeight, opts.RemoveWeight, 1.0)
	if maxPenalty == 0 {
		return Score{Total: int(total), Matching: 0, Similarity: 1.0}
	}

	similarity := 1.0 - (penalty / maxPenalty)
	if similarity < 0 {
		similarity = 0
	}

	return Score{
		Total:      int(total),
		Matching:   0,
		Similarity: similarity,
	}
}

func max3(a, b, c float64) float64 {
	if a >= b && a >= c {
		return a
	}
	if b >= c {
		return b
	}
	return c
}
