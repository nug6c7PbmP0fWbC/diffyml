package diff

import (
	"testing"
)

func TestReplay_Empty(t *testing.T) {
	base := map[string]interface{}{"a": "1"}
	results, final := Replay(base, nil, DefaultReplayOptions())
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
	if final["a"] != "1" {
		t.Fatalf("expected base unchanged")
	}
}

func TestReplay_ForwardAppliesChanges(t *testing.T) {
	base := map[string]interface{}{"a": "old"}
	steps := [][]Change{
		{{Path: "a", Type: ChangeModified, OldValue: "old", Value: "new"}},
		{{Path: "b", Type: ChangeAdded, Value: "added"}},
	}
	_, final := Replay(base, steps, DefaultReplayOptions())
	if final["a"] != "new" {
		t.Errorf("expected a=new, got %v", final["a"])
	}
	if final["b"] != "added" {
		t.Errorf("expected b=added, got %v", final["b"])
	}
}

func TestReplay_ReverseOrder(t *testing.T) {
	base := map[string]interface{}{}
	steps := [][]Change{
		{{Path: "x", Type: ChangeAdded, Value: "step0"}},
		{{Path: "x", Type: ChangeModified, OldValue: "step0", Value: "step1"}},
	}
	opts := DefaultReplayOptions()
	opts.Reverse = true
	_, final := Replay(base, steps, opts)
	// In reverse: step1 applied first (reverse=true means swap old/new), then step0
	if final["x"] != "step0" {
		t.Errorf("expected x=step0 after reverse replay, got %v", final["x"])
	}
}

func TestReplay_ResultsCount(t *testing.T) {
	base := map[string]interface{}{}
	steps := [][]Change{
		{{Path: "a", Type: ChangeAdded, Value: "1"}},
		{{Path: "b", Type: ChangeAdded, Value: "2"}},
		{{Path: "c", Type: ChangeAdded, Value: "3"}},
	}
	results, _ := Replay(base, steps, DefaultReplayOptions())
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestReplay_RemoveChange(t *testing.T) {
	base := map[string]interface{}{"key": "val"}
	steps := [][]Change{
		{{Path: "key", Type: ChangeRemoved, OldValue: "val"}},
	}
	_, final := Replay(base, steps, DefaultReplayOptions())
	if _, ok := final["key"]; ok {
		t.Error("expected key to be removed")
	}
}

func TestReplay_DoesNotMutateBase(t *testing.T) {
	base := map[string]interface{}{"a": "original"}
	steps := [][]Change{
		{{Path: "a", Type: ChangeModified, OldValue: "original", Value: "changed"}},
	}
	Replay(base, steps, DefaultReplayOptions())
	if base["a"] != "original" {
		t.Error("Replay must not mutate base map")
	}
}
