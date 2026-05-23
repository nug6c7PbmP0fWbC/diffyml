package formatter

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/szhekpisov/diffyml/pkg/diff"
)

func TestRankFormatter_NoRanked(t *testing.T) {
	f := NewRankFormatter(nil)
	out, err := f.Format(nil)
	assert.NoError(t, err)
	assert.Contains(t, out, "no ranked changes")
}

func TestRankFormatter_HeaderPresent(t *testing.T) {
	ranked := []diff.RankedChange{
		{Change: diff.Change{Path: "foo", Type: diff.ChangeAdded, NewValue: "bar"}, Score: 1.0},
	}
	f := NewRankFormatter(ranked)
	out, err := f.Format(nil)
	assert.NoError(t, err)
	assert.Contains(t, out, "Ranked Changes")
	assert.Contains(t, out, "---")
}

func TestRankFormatter_AddedLine(t *testing.T) {
	ranked := []diff.RankedChange{
		{Change: diff.Change{Path: "a.b", Type: diff.ChangeAdded, NewValue: "v1"}, Score: 1.0},
	}
	f := NewRankFormatter(ranked)
	out, err := f.Format(nil)
	assert.NoError(t, err)
	assert.Contains(t, out, "#1")
	assert.Contains(t, out, "+ a.b")
	assert.Contains(t, out, "v1")
}

func TestRankFormatter_ModifiedLine(t *testing.T) {
	ranked := []diff.RankedChange{
		{Change: diff.Change{Path: "x", Type: diff.ChangeModified, OldValue: "old", NewValue: "new"}, Score: 0.5},
	}
	f := NewRankFormatter(ranked)
	out, err := f.Format(nil)
	assert.NoError(t, err)
	assert.Contains(t, out, "~ x")
	assert.Contains(t, out, "old -> new")
	assert.Contains(t, out, "0.50")
}

func TestRankFormatter_MultipleEntries_Ordered(t *testing.T) {
	ranked := []diff.RankedChange{
		{Change: diff.Change{Path: "first", Type: diff.ChangeAdded}, Score: 1.0},
		{Change: diff.Change{Path: "second", Type: diff.ChangeRemoved}, Score: 0.8},
	}
	f := NewRankFormatter(ranked)
	out, err := f.Format(nil)
	assert.NoError(t, err)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	// header + separator + 2 entries
	assert.GreaterOrEqual(t, len(lines), 4)
	assert.True(t, strings.Contains(lines[2], "#1") && strings.Contains(lines[2], "first"))
	assert.True(t, strings.Contains(lines[3], "#2") && strings.Contains(lines[3], "second"))
}
