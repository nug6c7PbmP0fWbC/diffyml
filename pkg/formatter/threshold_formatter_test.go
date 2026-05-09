package formatter

import (
	"bytes"
	"testing"

	"github.com/szhekpisov/diffyml/pkg/diff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThresholdFormatter_NoViolations(t *testing.T) {
	f := NewThresholdFormatter()
	var buf bytes.Buffer
	err := f.Format(&buf, nil)
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "all limits satisfied")
}

func TestThresholdFormatter_SingleViolation(t *testing.T) {
	f := NewThresholdFormatter()
	violations := []diff.ThresholdViolation{
		{Kind: "added", Limit: 2, Actual: 5, Message: "added changes exceed threshold: 5 > 2"},
	}
	var buf bytes.Buffer
	err := f.Format(&buf, violations)
	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "1 violation(s)")
	assert.Contains(t, out, "[ADDED]")
	assert.Contains(t, out, "5 > 2")
}

func TestThresholdFormatter_MultipleViolations(t *testing.T) {
	f := NewThresholdFormatter()
	violations := []diff.ThresholdViolation{
		{Kind: "removed", Limit: 1, Actual: 3, Message: "removed changes exceed threshold: 3 > 1"},
		{Kind: "total", Limit: 4, Actual: 7, Message: "total changes exceed threshold: 7 > 4"},
	}
	var buf bytes.Buffer
	err := f.Format(&buf, violations)
	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "2 violation(s)")
	assert.Contains(t, out, "[REMOVED]")
	assert.Contains(t, out, "[TOTAL]")
}

func TestThresholdFormatter_EmptyViolations(t *testing.T) {
	f := NewThresholdFormatter()
	var buf bytes.Buffer
	err := f.Format(&buf, []diff.ThresholdViolation{})
	require.NoError(t, err)
	assert.Contains(t, buf.String(), "all limits satisfied")
}
