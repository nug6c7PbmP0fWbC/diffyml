package formatter

import (
	"strings"
	"testing"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

func TestReplayFormatter_NoSteps(t *testing.T) {
	f := NewReplayFormatter()
	out := f.Format(nil)
	if !strings.Contains(out, "no steps") {
		t.Errorf("expected 'no steps', got: %s", out)
	}
}

func TestReplayFormatter_HeaderPresent(t *testing.T) {
	f := NewReplayFormatter()
	results := []diff.ReplayResult{
		{Step: 0, Changes: []diff.Change{}},
	}
	out := f.Format(results)
	if !strings.Contains(out, "Replay Report") {
		t.Errorf("expected header, got: %s", out)
	}
}

func TestReplayFormatter_AddedChange(t *testing.T) {
	f := NewReplayFormatter()
	results := []diff.ReplayResult{
		{Step: 0, Changes: []diff.Change{
			{Path: "key", Type: diff.ChangeAdded, Value: "val"},
		}},
	}
	out := f.Format(results)
	if !strings.Contains(out, "+ key = val") {
		t.Errorf("expected added line, got: %s", out)
	}
}

func TestReplayFormatter_RemovedChange(t *testing.T) {
	f := NewReplayFormatter()
	results := []diff.ReplayResult{
		{Step: 1, Changes: []diff.Change{
			{Path: "old", Type: diff.ChangeRemoved, OldValue: "gone"},
		}},
	}
	out := f.Format(results)
	if !strings.Contains(out, "- old (was gone)") {
		t.Errorf("expected removed line, got: %s", out)
	}
}

func TestReplayFormatter_ModifiedChange(t *testing.T) {
	f := NewReplayFormatter()
	results := []diff.ReplayResult{
		{Step: 0, Changes: []diff.Change{
			{Path: "ver", Type: diff.ChangeModified, OldValue: "1", Value: "2"},
		}},
	}
	out := f.Format(results)
	if !strings.Contains(out, "~ ver: 1 -> 2") {
		t.Errorf("expected modified line, got: %s", out)
	}
}

func TestReplayFormatter_StepIndex(t *testing.T) {
	f := NewReplayFormatter()
	results := []diff.ReplayResult{
		{Step: 3, Changes: []diff.Change{}},
	}
	out := f.Format(results)
	if !strings.Contains(out, "Step 3") {
		t.Errorf("expected step index 3, got: %s", out)
	}
}
