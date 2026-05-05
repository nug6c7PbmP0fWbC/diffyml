package diff

import (
	"fmt"
	"strings"
	"time"
)

// ChangelogEntry represents a single entry in a human-readable changelog.
type ChangelogEntry struct {
	Timestamp time.Time
	Version   string
	Changes   []Change
	Summary   Summary
}

// Changelog holds a list of changelog entries ordered by time.
type Changelog struct {
	Entries []ChangelogEntry
}

// AddEntry appends a new entry to the changelog.
func (c *Changelog) AddEntry(version string, changes []Change) {
	entry := ChangelogEntry{
		Timestamp: time.Now().UTC(),
		Version:   version,
		Changes:   changes,
		Summary:   Summarize(changes),
	}
	c.Entries = append(c.Entries, entry)
}

// Render returns a human-readable text representation of the changelog.
func (c *Changelog) Render() string {
	var sb strings.Builder
	for _, entry := range c.Entries {
		sb.WriteString(fmt.Sprintf("## [%s] - %s\n", entry.Version, entry.Timestamp.Format("2006-01-02")))
		sb.WriteString(fmt.Sprintf("  Added: %d | Removed: %d | Modified: %d\n",
			entry.Summary.Added, entry.Summary.Removed, entry.Summary.Modified))
		for _, ch := range entry.Changes {
			sb.WriteString(fmt.Sprintf("  - [%s] %s", strings.ToUpper(string(ch.Type)), ch.Path))
			switch ch.Type {
			case ChangeAdded:
				sb.WriteString(fmt.Sprintf(": %v", ch.NewValue))
			case ChangeRemoved:
				sb.WriteString(fmt.Sprintf(": %v", ch.OldValue))
			case ChangeModified:
				sb.WriteString(fmt.Sprintf(": %v -> %v", ch.OldValue, ch.NewValue))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
