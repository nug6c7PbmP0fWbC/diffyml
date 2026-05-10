// Package diff provides utilities for comparing, transforming, and replaying
// YAML change sets.
//
// # Replay
//
// Replay applies a sequence of change slices to a base map, producing an
// intermediate state after each step. This is useful for auditing how a
// configuration evolved over time or for implementing undo/redo workflows.
//
// Basic usage:
//
//	base := map[string]interface{}{"version": "1.0"}
//	steps := [][]diff.Change{
//	    {{Path: "version", Type: diff.ChangeModified, OldValue: "1.0", Value: "1.1"}},
//	    {{Path: "debug",   Type: diff.ChangeAdded,    Value: true}},
//	}
//	results, final := diff.Replay(base, steps, diff.DefaultReplayOptions())
//
// Reverse replay:
//
//	opts := diff.DefaultReplayOptions()
//	opts.Reverse = true
//	_, rolledBack := diff.Replay(final, steps, opts)
//
// Each ReplayResult records the step index, the changes applied at that step,
// and any error that occurred. When StopOnError is true (the default) replay
// halts on the first error.
package diff
