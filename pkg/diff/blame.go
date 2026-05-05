package diff

import "time"

// BlameEntry records who last changed a given path and when.
type BlameEntry struct {
	Path      string    `json:"path" yaml:"path"`
	Author    string    `json:"author" yaml:"author"`
	ChangedAt time.Time `json:"changed_at" yaml:"changed_at"`
	ChangeType string   `json:"change_type" yaml:"change_type"`
}

// BlameMap maps a path to its most recent BlameEntry.
type BlameMap map[string]BlameEntry

// Blame builds a BlameMap from a slice of Changes, tagging each path
// with the provided author and timestamp. If a path appears more than
// once, the last entry wins (most-recent-write semantics).
func Blame(changes []Change, author string, at time.Time) BlameMap {
	bm := make(BlameMap, len(changes))
	for _, c := range changes {
		bm[c.Path] = BlameEntry{
			Path:       c.Path,
			Author:     author,
			ChangedAt:  at,
			ChangeType: string(c.Type),
		}
	}
	return bm
}

// MergeBlame merges src into dst. Entries in src overwrite entries in dst
// when the same path is present in both maps.
func MergeBlame(dst, src BlameMap) BlameMap {
	out := make(BlameMap, len(dst)+len(src))
	for k, v := range dst {
		out[k] = v
	}
	for k, v := range src {
		out[k] = v
	}
	return out
}
