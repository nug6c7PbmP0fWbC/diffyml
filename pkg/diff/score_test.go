package diff

import (
	"testing"
)

func TestComputeScore_NoChanges(t *testing.T) {
	s := ComputeScore([]Change{}, DefaultScoreOptions())
	if s.Similarity != 1.0 {
		t.Errorf("expected similarity 1.0, got %f", s.Similarity)
	}
	if s.Total != 0 {
		t.Errorf("expected total 0, got %d", s.Total)
	}
}

func TestComputeScore_AllAdded(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeAdded},
	}
	s := ComputeScore(changes, DefaultScoreOptions())
	if s.Similarity >= 1.0 {
		t.Errorf("expected similarity < 1.0, got %f", s.Similarity)
	}
	if s.Total != 2 {
		t.Errorf("expected total 2, got %d", s.Total)
	}
}

func TestComputeScore_AllRemoved(t *testing.T) {
	changes := []Change{
		{Path: "x", Type: ChangeRemoved},
	}
	s := ComputeScore(changes, DefaultScoreOptions())
	if s.Similarity >= 1.0 {
		t.Errorf("expected similarity < 1.0, got %f", s.Similarity)
	}
}

func TestComputeScore_ModifiedGetsPartialCredit(t *testing.T) {
	modOpts := DefaultScoreOptions()
	modOpts.ModifyWeight = 0.5

	changesModified := []Change{{Path: "a", Type: ChangeModified}}
	changesAdded := []Change{{Path: "a", Type: ChangeAdded}}

	sModified := ComputeScore(changesModified, modOpts)
	sAdded := ComputeScore(changesAdded, modOpts)

	if sModified.Similarity <= sAdded.Similarity {
		t.Errorf("modified should score higher than added: modified=%f added=%f",
			sModified.Similarity, sAdded.Similarity)
	}
}

func TestComputeScore_SimilarityBounds(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeRemoved},
		{Path: "c", Type: ChangeModified},
	}
	s := ComputeScore(changes, DefaultScoreOptions())
	if s.Similarity < 0 || s.Similarity > 1.0 {
		t.Errorf("similarity out of bounds: %f", s.Similarity)
	}
}

func TestDefaultScoreOptions(t *testing.T) {
	opts := DefaultScoreOptions()
	if opts.AddWeight != 1.0 || opts.RemoveWeight != 1.0 || opts.ModifyWeight != 0.5 {
		t.Errorf("unexpected default options: %+v", opts)
	}
}
