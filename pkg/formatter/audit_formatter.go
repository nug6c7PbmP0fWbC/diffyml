package formatter

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// AuditFormatter renders an AuditLog as a human-readable report.
type AuditFormatter struct {
	w io.Writer
}

// NewAuditFormatter creates an AuditFormatter writing to w.
func NewAuditFormatter(w io.Writer) *AuditFormatter {
	return &AuditFormatter{w: w}
}

// Format writes the audit log to the underlying writer.
func (f *AuditFormatter) Format(log diff.AuditLog) error {
	if len(log.Entries) == 0 {
		_, err := fmt.Fprintln(f.w, "audit: no entries")
		return err
	}

	// Sort entries by timestamp then path for deterministic output.
	sorted := make([]diff.AuditEntry, len(log.Entries))
	copy(sorted, log.Entries)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Timestamp.Equal(sorted[j].Timestamp) {
			return sorted[i].Path < sorted[j].Path
		}
		return sorted[i].Timestamp.Before(sorted[j].Timestamp)
	})

	_, err := fmt.Fprintf(f.w, "audit log (%d entries):\n", len(sorted))
	if err != nil {
		return err
	}
	for _, e := range sorted {
		line := buildAuditLine(e)
		if _, err := fmt.Fprintln(f.w, line); err != nil {
			return err
		}
	}
	return nil
}

func buildAuditLine(e diff.AuditEntry) string {
	ts := e.Timestamp.Format("2006-01-02T15:04:05Z")
	op := strings.ToUpper(e.Operation)
	switch e.Operation {
	case "added":
		return fmt.Sprintf("  [%s] %s actor=%s path=%s new=%v", ts, op, e.Actor, e.Path, e.NewValue)
	case "removed":
		return fmt.Sprintf("  [%s] %s actor=%s path=%s old=%v", ts, op, e.Actor, e.Path, e.OldValue)
	default:
		return fmt.Sprintf("  [%s] %s actor=%s path=%s old=%v new=%v", ts, op, e.Actor, e.Path, e.OldValue, e.NewValue)
	}
}
