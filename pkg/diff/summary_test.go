package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/szhekpisov/diffyml/pkg/diff/change"
)

func TestSummarize_Empty(t *testing.T) {
	s := Summarize(nil)
	assert.Equal(t, 0, s.Total)
	assert.False(t, s.HasChanges())
}

func TestSummarize_AllTypes(t *testing.T) {
	changes := []change.Change{
		{Type: change.Added, Path: "a"},
		{Type: change.Added, Path: "b"},
		{Type: change.Removed, Path: "c"},
		{Type: change.Modified, Path: "d"},
	}
	s := Summarize(changes)
	assert.Equal(t, 2, s.Added)
	assert.Equal(t, 1, s.Removed)
	assert.Equal(t, 1, s.Modified)
	assert.Equal(t, 4, s.Total)
	assert.True(t, s.HasChanges())
}

func TestSummarize_OnlyAdded(t *testing.T) {
	changes := []change.Change{
		{Type: change.Added, Path: "x"},
	}
	s := Summarize(changes)
	assert.Equal(t, 1, s.Added)
	assert.Equal(t, 0, s.Removed)
	assert.Equal(t, 0, s.Modified)
	assert.Equal(t, 1, s.Total)
}

func TestSummarize_OnlyRemoved(t *testing.T) {
	changes := []change.Change{
		{Type: change.Removed, Path: "y"},
		{Type: change.Removed, Path: "z"},
	}
	s := Summarize(changes)
	assert.Equal(t, 0, s.Added)
	assert.Equal(t, 2, s.Removed)
	assert.Equal(t, 0, s.Modified)
	assert.Equal(t, 2, s.Total)
}

func TestSummarize_OnlyModified(t *testing.T) {
	changes := []change.Change{
		{Type: change.Modified, Path: "m"},
	}
	s := Summarize(changes)
	assert.Equal(t, 0, s.Added)
	assert.Equal(t, 0, s.Removed)
	assert.Equal(t, 1, s.Modified)
	assert.True(t, s.HasChanges())
}
