// Package diff provides utilities for comparing, merging, and summarising
// differences between two YAML documents represented as Go maps.
//
// # Watcher
//
// The Watcher type enables continuous monitoring of two YAML sources.
// It polls both sources on a configurable interval and emits a WatchEvent
// whenever the diff between them changes.
//
// Basic usage:
//
//		w := diff.NewWatcher(5 * time.Second)
//		go w.Watch(
//			"base.yml", "override.yml",
//			func() (map[string]interface{}, error) {
//				return loader.LoadFile("base.yml")
//			},
//			func() (map[string]interface{}, error) {
//				return loader.LoadFile("override.yml")
//			},
//		)
//
//		for evt := range w.Events() {
//			fmt.Printf("[%s] %d change(s) detected\n",
//				evt.DetectedAt.Format(time.RFC3339), len(evt.Changes))
//		}
//
// Call w.Stop() to terminate the polling goroutine.
//
// Duplicate events are suppressed: a new WatchEvent is only emitted when
// the set of changes differs from the previous poll cycle.
package diff
