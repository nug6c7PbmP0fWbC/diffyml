package diff

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// WatchEvent represents a detected change event between two polling cycles.
type WatchEvent struct {
	Changes   []Change
	DetectedAt time.Time
	FileA     string
	FileB     string
}

// Watcher polls two YAML sources at a given interval and emits events when
// differences are detected.
type Watcher struct {
	interval time.Duration
	events   chan WatchEvent
	stop     chan struct{}
	lastHash string
}

// NewWatcher creates a Watcher that checks for changes every interval.
func NewWatcher(interval time.Duration) *Watcher {
	return &Watcher{
		interval: interval,
		events:   make(chan WatchEvent, 8),
		stop:     make(chan struct{}),
	}
}

// Events returns the read-only channel of WatchEvents.
func (w *Watcher) Events() <-chan WatchEvent {
	return w.events
}

// Stop signals the watcher to cease polling.
func (w *Watcher) Stop() {
	close(w.stop)
}

// Watch begins polling dataA and dataB using the provided loader functions.
// Each function should return the latest parsed YAML map.
func (w *Watcher) Watch(fileA, fileB string, loadA, loadB func() (map[string]interface{}, error)) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-w.stop:
			return
		case <-ticker.C:
			a, err := loadA()
			if err != nil {
				continue
			}
			b, err := loadB()
			if err != nil {
				continue
			}

			changes := Compare(a, b)
			hash := hashChanges(changes)
			if hash == w.lastHash {
				continue
			}
			w.lastHash = hash

			if len(changes) > 0 {
				w.events <- WatchEvent{
					Changes:    changes,
					DetectedAt: time.Now(),
					FileA:      fileA,
					FileB:      fileB,
				}
			}
		}
	}
}

func hashChanges(changes []Change) string {
	h := sha256.New()
	for _, c := range changes {
		fmt.Fprintf(h, "%s|%s|%v|%v;", c.Path, c.Type, c.OldValue, c.NewValue)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
