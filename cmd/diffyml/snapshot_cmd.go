package main

import (
	"fmt"
	"os"

	"github.com/szhekpisov/diffyml/pkg/diff"
	"github.com/szhekpisov/diffyml/pkg/loader"
)

// runSnapshot handles the 'snapshot' subcommand.
// Usage: diffyml snapshot <yaml-file> <output-snap.json> [label]
func runSnapshot(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: diffyml snapshot <yaml-file> <snap.json> [label]")
	}
	yamlPath := args[0]
	snapPath := args[1]
	label := "snapshot"
	if len(args) >= 3 {
		label = args[2]
	}

	data, err := loader.LoadFile(yamlPath)
	if err != nil {
		return fmt.Errorf("snapshot: load yaml: %w", err)
	}

	if err := diff.SaveSnapshot(snapPath, label, data); err != nil {
		return fmt.Errorf("snapshot: save: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Snapshot saved to %s (label: %s)\n", snapPath, label)
	return nil
}

// runSnapshotDiff handles the 'snapshot-diff' subcommand.
// Usage: diffyml snapshot-diff <snap.json> <yaml-file>
func runSnapshotDiff(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: diffyml snapshot-diff <snap.json> <yaml-file>")
	}
	snapPath := args[0]
	yamlPath := args[1]

	snap, err := diff.LoadSnapshot(snapPath)
	if err != nil {
		return fmt.Errorf("snapshot-diff: load snapshot: %w", err)
	}

	current, err := loader.LoadFile(yamlPath)
	if err != nil {
		return fmt.Errorf("snapshot-diff: load yaml: %w", err)
	}

	changes := diff.DiffSnapshot(snap, current)
	if len(changes) == 0 {
		fmt.Fprintln(os.Stdout, "No changes since snapshot:", snap.Label)
		return nil
	}

	fmt.Fprintf(os.Stdout, "%d change(s) since snapshot '%s' (%s):\n",
		len(changes), snap.Label, snap.Timestamp.Format("2006-01-02T15:04:05Z"))
	for _, c := range changes {
		fmt.Fprintf(os.Stdout, "  [%s] %s\n", c.Type, c.Path)
	}
	return nil
}
