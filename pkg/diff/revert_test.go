package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRevert_Empty(t *testing.T) {
	result := Revert(nil)
	assert.Empty(t, result)
}

func TestRevert_Added_BecomesRemoved(t *testing.T) {
	changes := []Change{
		{Type: ChangeAdded, Path: "a.b", Before: nil, After: "hello"},
	}
	result := Revert(changes)
	assert.Len(t, result, 1)
	assert.Equal(t, ChangeRemoved, result[0].Type)
	assert.Equal(t, "a.b", result[0].Path)
	assert.Equal(t, "hello", result[0].Before)
	assert.Nil(t, result[0].After)
}

func TestRevert_Removed_BecomesAdded(t *testing.T) {
	changes := []Change{
		{Type: ChangeRemoved, Path: "x.y", Before: 42, After: nil},
	}
	result := Revert(changes)
	assert.Len(t, result, 1)
	assert.Equal(t, ChangeAdded, result[0].Type)
	assert.Equal(t, "x.y", result[0].Path)
	assert.Nil(t, result[0].Before)
	assert.Equal(t, 42, result[0].After)
}

func TestRevert_Modified_SwapsValues(t *testing.T) {
	changes := []Change{
		{Type: ChangeModified, Path: "cfg.timeout", Before: 30, After: 60},
	}
	result := Revert(changes)
	assert.Len(t, result, 1)
	assert.Equal(t, ChangeModified, result[0].Type)
	assert.Equal(t, 60, result[0].Before)
	assert.Equal(t, 30, result[0].After)
}

func TestRevert_DoesNotMutateOriginal(t *testing.T) {
	original := []Change{
		{Type: ChangeAdded, Path: "k", Before: nil, After: "v"},
	}
	Revert(original)
	assert.Equal(t, ChangeAdded, original[0].Type)
	assert.Nil(t, original[0].Before)
	assert.Equal(t, "v", original[0].After)
}

func TestRevert_MultipleChanges(t *testing.T) {
	changes := []Change{
		{Type: ChangeAdded, Path: "a", Before: nil, After: 1},
		{Type: ChangeRemoved, Path: "b", Before: 2, After: nil},
		{Type: ChangeModified, Path: "c", Before: "old", After: "new"},
	}
	result := Revert(changes)
	assert.Len(t, result, 3)
	assert.Equal(t, ChangeRemoved, result[0].Type)
	assert.Equal(t, ChangeAdded, result[1].Type)
	assert.Equal(t, ChangeModified, result[2].Type)
	assert.Equal(t, "new", result[2].Before)
	assert.Equal(t, "old", result[2].After)
}
