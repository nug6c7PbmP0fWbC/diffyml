package diff

import "sort"

// RankOptions controls how changes are ranked.
type RankOptions struct {
	// Weights for each change type (0.0 – 1.0 scale, default 1.0).
	AddedWeight    float64
	RemovedWeight  float64
	ModifiedWeight float64
	// TopN limits the result to the N highest-ranked changes.
	// 0 means no limit.
	TopN int
}

// DefaultRankOptions returns sensible defaults.
func DefaultRankOptions() RankOptions {
	return RankOptions{
		AddedWeight:    1.0,
		RemovedWeight:  1.0,
		ModifiedWeight: 0.5,
		TopN:           0,
	}
}

// RankedChange pairs a Change with its computed score.
type RankedChange struct {
	Change Change
	Score  float64
}

// Rank scores each change according to opts and returns them sorted
// from highest to lowest score. Ties are broken by path (ascending).
func Rank(changes []Change, opts RankOptions) []RankedChange {
	ranked := make([]RankedChange, 0, len(changes))

	for _, c := range changes {
		var w float64
		switch c.Type {
		case ChangeAdded:
			w = opts.AddedWeight
		case ChangeRemoved:
			w = opts.RemovedWeight
		case ChangeModified:
			w = opts.ModifiedWeight
		default:
			w = 1.0
		}
		ranked = append(ranked, RankedChange{Change: c, Score: w})
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].Score != ranked[j].Score {
			return ranked[i].Score > ranked[j].Score
		}
		return ranked[i].Change.Path < ranked[j].Change.Path
	})

	if opts.TopN > 0 && len(ranked) > opts.TopN {
		ranked = ranked[:opts.TopN]
	}
	return ranked
}
