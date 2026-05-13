package diff

import (
	"fmt"
	"strings"
	"time"
)

// TimelineEntry records a labelled snapshot of changes at a point in time.
type TimelineEntry struct {
	Label     string
	Timestamp time.Time
	Changes   []Change
}

// Timeline is an ordered sequence of TimelineEntry values.
type Timeline struct {
	Entries []TimelineEntry
}

// AddEntry appends a new entry to the timeline.
func (t *Timeline) AddEntry(label string, changes []Change) {
	t.Entries = append(t.Entries, TimelineEntry{
		Label:     label,
		Timestamp: time.Now().UTC(),
		Changes:   changes,
	})
}

// Render returns a human-readable summary of the timeline.
func (t *Timeline) Render() string {
	if len(t.Entries) == 0 {
		return "(empty timeline)\n"
	}
	var sb strings.Builder
	for _, e := range t.Entries {
		sb.WriteString(fmt.Sprintf("[%s] %s — %d change(s)\n",
			e.Timestamp.Format(time.RFC3339),
			e.Label,
			len(e.Changes),
		))
		for _, c := range e.Changes {
			sb.WriteString(fmt.Sprintf("  %s %s\n", c.Type, c.Path))
		}
	}
	return sb.String()
}

// TotalChanges returns the total number of changes across all entries.
func (t *Timeline) TotalChanges() int {
	n := 0
	for _, e := range t.Entries {
		n += len(e.Changes)
	}
	return n
}
