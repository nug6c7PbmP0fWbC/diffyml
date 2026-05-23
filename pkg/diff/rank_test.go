package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRank_Empty(t *testing.T) {
	result := Rank(nil, DefaultRankOptions())
	assert.Empty(t, result)
}

func TestRank_SortsByWeightDescending(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeModified},
		{Path: "b", Type: ChangeAdded},
		{Path: "c", Type: ChangeRemoved},
	}
	opts := DefaultRankOptions() // Added=1.0, Removed=1.0, Modified=0.5
	result := Rank(changes, opts)

	assert.Len(t, result, 3)
	// Added and Removed both score 1.0; Modified scores 0.5
	assert.Equal(t, 0.5, result[len(result)-1].Score)
	assert.Equal(t, ChangeModified, result[len(result)-1].Change.Type)
}

func TestRank_TiesBrokenByPath(t *testing.T) {
	changes := []Change{
		{Path: "z", Type: ChangeAdded},
		{Path: "a", Type: ChangeAdded},
		{Path: "m", Type: ChangeAdded},
	}
	result := Rank(changes, DefaultRankOptions())

	assert.Equal(t, "a", result[0].Change.Path)
	assert.Equal(t, "m", result[1].Change.Path)
	assert.Equal(t, "z", result[2].Change.Path)
}

func TestRank_TopNLimitsResults(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeRemoved},
		{Path: "c", Type: ChangeModified},
		{Path: "d", Type: ChangeAdded},
	}
	opts := DefaultRankOptions()
	opts.TopN = 2
	result := Rank(changes, opts)
	assert.Len(t, result, 2)
}

func TestRank_TopNZeroMeansNoLimit(t *testing.T) {
	changes := make([]Change, 10)
	for i := range changes {
		changes[i] = Change{Path: "p", Type: ChangeAdded}
	}
	opts := DefaultRankOptions()
	opts.TopN = 0
	result := Rank(changes, opts)
	assert.Len(t, result, 10)
}

func TestRank_CustomWeights(t *testing.T) {
	changes := []Change{
		{Path: "x", Type: ChangeModified},
		{Path: "y", Type: ChangeAdded},
	}
	opts := DefaultRankOptions()
	opts.ModifiedWeight = 2.0
	opts.AddedWeight = 0.1
	result := Rank(changes, opts)

	assert.Equal(t, ChangeModified, result[0].Change.Type)
	assert.Equal(t, 2.0, result[0].Score)
}
