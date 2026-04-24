// Package diff provides utilities for comparing, watching, merging,
// summarizing, and snapshotting YAML documents.
//
// # Snapshot
//
// The snapshot feature allows you to capture the state of a YAML document
// at a specific point in time and later compare it against a newer version
// to detect drift or unintended changes.
//
// Basic usage:
//
//	// Save a snapshot
//	_ = diff.SaveSnapshot("baseline.snap.json", "v1-release", data)
//
//	// Later, load and compare
//	snap, _ := diff.LoadSnapshot("baseline.snap.json")
//	changes := diff.DiffSnapshot(snap, currentData)
//
// Snapshots are stored as JSON files containing the label, UTC timestamp,
// and the full document data.
package diff
