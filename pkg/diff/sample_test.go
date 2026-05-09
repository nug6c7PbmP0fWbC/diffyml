package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeSampleChanges(n int) []Change {
	out := make([]Change, n)
	for i := range out {
		out[i] = Change{Path: fmt.Sprintf("key.%d", i), Type: Modified}
	}
	return out
}

func TestSample_Empty(t *testing.T) {
	result := Sample(nil, DefaultSampleOptions())
	assert.Nil(t, result)
}

func TestSample_NoLimit(t *testing.T) {
	changes := makeSampleChanges(5)
	opts := DefaultSampleOptions() // N == 0 means no limit
	result := Sample(changes, opts)
	assert.Equal(t, changes, result)
}

func TestSample_NGreaterThanLen(t *testing.T) {
	changes := makeSampleChanges(3)
	opts := SampleOptions{N: 10, Deterministic: true, Seed: 42}
	result := Sample(changes, opts)
	assert.Equal(t, changes, result)
}

func TestSample_ReducesToN(t *testing.T) {
	changes := makeSampleChanges(20)
	opts := SampleOptions{N: 5, Deterministic: true, Seed: 1}
	result := Sample(changes, opts)
	require.Len(t, result, 5)
}

func TestSample_Deterministic(t *testing.T) {
	changes := makeSampleChanges(50)
	opts := SampleOptions{N: 10, Deterministic: true, Seed: 99}
	first := Sample(changes, opts)
	second := Sample(changes, opts)
	assert.Equal(t, first, second, "same seed must produce same sample")
}

func TestSample_DifferentSeeds(t *testing.T) {
	changes := makeSampleChanges(50)
	a := Sample(changes, SampleOptions{N: 10, Deterministic: true, Seed: 1})
	b := Sample(changes, SampleOptions{N: 10, Deterministic: true, Seed: 2})
	// With 50 items and picking 10 it is astronomically unlikely both seeds
	// produce the identical ordered subset.
	assert.NotEqual(t, a, b)
}

func TestSample_PreservesOriginalOrder(t *testing.T) {
	changes := makeSampleChanges(20)
	opts := SampleOptions{N: 10, Deterministic: true, Seed: 7}
	result := Sample(changes, opts)

	// Verify each returned change appears in the original slice in order.
	prev := -1
	for _, c := range result {
		for i, orig := range changes {
			if orig.Path == c.Path {
				assert.Greater(t, i, prev, "order must be preserved")
				prev = i
				break
			}
		}
	}
}
