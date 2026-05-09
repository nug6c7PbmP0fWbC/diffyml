package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeWindowChanges(n int) []Change {
	out := make([]Change, n)
	for i := range out {
		out[i] = Change{Path: fmt.Sprintf("key%d", i), Type: ChangeModified}
	}
	return out
}

func TestWindow_Empty(t *testing.T) {
	windows, err := Window([]Change{}, DefaultWindowOptions())
	require.NoError(t, err)
	assert.Empty(t, windows)
}

func TestWindow_SmallerThanSize(t *testing.T) {
	changes := makeWindowChanges(3)
	windows, err := Window(changes, WindowOptions{Size: 10, Step: 10})
	require.NoError(t, err)
	require.Len(t, windows, 1)
	assert.Len(t, windows[0], 3)
}

func TestWindow_ExactMultiple(t *testing.T) {
	changes := makeWindowChanges(6)
	windows, err := Window(changes, WindowOptions{Size: 3, Step: 3})
	require.NoError(t, err)
	require.Len(t, windows, 2)
	assert.Len(t, windows[0], 3)
	assert.Len(t, windows[1], 3)
}

func TestWindow_NonDivisible(t *testing.T) {
	changes := makeWindowChanges(7)
	windows, err := Window(changes, WindowOptions{Size: 3, Step: 3})
	require.NoError(t, err)
	require.Len(t, windows, 3)
	assert.Len(t, windows[2], 1)
}

func TestWindow_Overlapping(t *testing.T) {
	changes := makeWindowChanges(5)
	windows, err := Window(changes, WindowOptions{Size: 3, Step: 1})
	require.NoError(t, err)
	// starts: 0,1,2,3,4 but end caps at 5
	assert.Len(t, windows, 5)
	assert.Len(t, windows[0], 3)
	assert.Len(t, windows[4], 1)
}

func TestWindow_InvalidSize(t *testing.T) {
	_, err := Window(makeWindowChanges(5), WindowOptions{Size: 0, Step: 1})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Size")
}

func TestWindow_InvalidStep(t *testing.T) {
	_, err := Window(makeWindowChanges(5), WindowOptions{Size: 3, Step: 0})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Step")
}

func TestWindow_DoesNotMutateOriginal(t *testing.T) {
	changes := makeWindowChanges(4)
	orig := make([]Change, len(changes))
	copy(orig, changes)
	_, err := Window(changes, WindowOptions{Size: 2, Step: 2})
	require.NoError(t, err)
	assert.Equal(t, orig, changes)
}
