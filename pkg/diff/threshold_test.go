package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckThreshold_NoChanges(t *testing.T) {
	opts := DefaultThresholdOptions()
	violations := CheckThreshold(nil, opts)
	assert.Empty(t, violations)
}

func TestCheckThreshold_NoLimits(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeRemoved},
		{Path: "c", Type: ChangeModified},
	}
	violations := CheckThreshold(changes, DefaultThresholdOptions())
	assert.Empty(t, violations)
}

func TestCheckThreshold_AddedExceeds(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeAdded},
		{Path: "c", Type: ChangeAdded},
	}
	opts := ThresholdOptions{MaxAdded: 2}
	violations := CheckThreshold(changes, opts)
	assert.Len(t, violations, 1)
	assert.Equal(t, "added", violations[0].Kind)
	assert.Equal(t, 2, violations[0].Limit)
	assert.Equal(t, 3, violations[0].Actual)
}

func TestCheckThreshold_RemovedExceeds(t *testing.T) {
	changes := []Change{
		{Path: "x", Type: ChangeRemoved},
		{Path: "y", Type: ChangeRemoved},
	}
	opts := ThresholdOptions{MaxRemoved: 1}
	violations := CheckThreshold(changes, opts)
	assert.Len(t, violations, 1)
	assert.Equal(t, "removed", violations[0].Kind)
}

func TestCheckThreshold_ModifiedExceeds(t *testing.T) {
	changes := []Change{
		{Path: "m", Type: ChangeModified},
	}
	opts := ThresholdOptions{MaxModified: 0} // 0 = unlimited, should not trigger
	violations := CheckThreshold(changes, opts)
	assert.Empty(t, violations)
}

func TestCheckThreshold_TotalExceeds(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeRemoved},
		{Path: "c", Type: ChangeModified},
	}
	opts := ThresholdOptions{MaxTotal: 2}
	violations := CheckThreshold(changes, opts)
	assert.Len(t, violations, 1)
	assert.Equal(t, "total", violations[0].Kind)
	assert.Equal(t, 3, violations[0].Actual)
}

func TestCheckThreshold_MultipleViolations(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeAdded},
		{Path: "c", Type: ChangeRemoved},
		{Path: "d", Type: ChangeRemoved},
	}
	opts := ThresholdOptions{MaxAdded: 1, MaxRemoved: 1, MaxTotal: 3}
	violations := CheckThreshold(changes, opts)
	assert.Len(t, violations, 3)
}

func TestCheckThreshold_ExactlyAtLimit(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeAdded},
	}
	opts := ThresholdOptions{MaxAdded: 2}
	violations := CheckThreshold(changes, opts)
	assert.Empty(t, violations)
}
