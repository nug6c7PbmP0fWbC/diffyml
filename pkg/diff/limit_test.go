package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeNChanges(n int) []Change {
	out := make([]Change, n)
	for i := 0; i < n; i++ {
		out[i] = Change{Type: ChangeModified, Path: fmt.Sprintf("key%d", i)}
	}
	return out
}

func TestLimit_Empty(t *testing.T) {
	result, err := Limit([]Change{}, DefaultLimitOptions())
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestLimit_NoLimit(t *testing.T) {
	changes := []Change{
		{Type: ChangeAdded, Path: "a"},
		{Type: ChangeRemoved, Path: "b"},
	}
	result, err := Limit(changes, DefaultLimitOptions())
	require.NoError(t, err)
	assert.Equal(t, changes, result)
}

func TestLimit_MaxCaps(t *testing.T) {
	changes := []Change{
		{Type: ChangeAdded, Path: "a"},
		{Type: ChangeAdded, Path: "b"},
		{Type: ChangeAdded, Path: "c"},
	}
	opts := LimitOptions{Max: 2, Offset: 0}
	result, err := Limit(changes, opts)
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "a", result[0].Path)
	assert.Equal(t, "b", result[1].Path)
}

func TestLimit_OffsetSkips(t *testing.T) {
	changes := []Change{
		{Type: ChangeAdded, Path: "a"},
		{Type: ChangeAdded, Path: "b"},
		{Type: ChangeAdded, Path: "c"},
	}
	opts := LimitOptions{Max: 0, Offset: 1}
	result, err := Limit(changes, opts)
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "b", result[0].Path)
	assert.Equal(t, "c", result[1].Path)
}

func TestLimit_OffsetAndMax(t *testing.T) {
	changes := []Change{
		{Type: ChangeAdded, Path: "a"},
		{Type: ChangeAdded, Path: "b"},
		{Type: ChangeAdded, Path: "c"},
		{Type: ChangeAdded, Path: "d"},
	}
	opts := LimitOptions{Max: 2, Offset: 1}
	result, err := Limit(changes, opts)
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "b", result[0].Path)
	assert.Equal(t, "c", result[1].Path)
}

func TestLimit_OffsetBeyondLength(t *testing.T) {
	changes := []Change{
		{Type: ChangeAdded, Path: "a"},
	}
	opts := LimitOptions{Max: 5, Offset: 10}
	result, err := Limit(changes, opts)
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestLimit_NegativeOffset(t *testing.T) {
	_, err := Limit([]Change{}, LimitOptions{Offset: -1})
	assert.Error(t, err)
}

func TestLimit_NegativeMax(t *testing.T) {
	_, err := Limit([]Change{}, LimitOptions{Max: -1})
	assert.Error(t, err)
}

func TestLimit_DoesNotMutateOriginal(t *testing.T) {
	changes := []Change{
		{Type: ChangeAdded, Path: "a"},
		{Type: ChangeAdded, Path: "b"},
	}
	opts := LimitOptions{Max: 1, Offset: 0}
	result, err := Limit(changes, opts)
	require.NoError(t, err)
	result[0].Path = "mutated"
	assert.Equal(t, "a", changes[0].Path)
}
