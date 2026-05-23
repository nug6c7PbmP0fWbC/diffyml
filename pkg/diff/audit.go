package diff

import "time"

// AuditEntry records a single audited diff operation.
type AuditEntry struct {
	Timestamp time.Time
	Operation string
	Path      string
	OldValue  interface{}
	NewValue  interface{}
	Actor     string
}

// AuditLog holds a sequence of audit entries.
type AuditLog struct {
	Entries []AuditEntry
}

// AuditOptions controls how changes are audited.
type AuditOptions struct {
	// Actor is the identifier of the entity performing the operation.
	Actor string
	// Operations is a list of change types to audit ("added", "removed", "modified").
	// Empty means audit all.
	Operations []string
}

// DefaultAuditOptions returns sensible defaults.
func DefaultAuditOptions() AuditOptions {
	return AuditOptions{
		Actor:      "unknown",
		Operations: nil,
	}
}

// Audit converts a slice of Changes into an AuditLog, recording each change
// as an AuditEntry stamped with the provided time.
func Audit(changes []Change, ts time.Time, opts AuditOptions) AuditLog {
	log := AuditLog{}
	allowed := make(map[string]bool)
	for _, op := range opts.Operations {
		allowed[op] = true
	}
	for _, c := range changes {
		op := string(c.Type)
		if len(allowed) > 0 && !allowed[op] {
			continue
		}
		log.Entries = append(log.Entries, AuditEntry{
			Timestamp: ts,
			Operation: op,
			Path:      c.Path,
			OldValue:  c.OldValue,
			NewValue:  c.NewValue,
			Actor:     opts.Actor,
		})
	}
	return log
}
