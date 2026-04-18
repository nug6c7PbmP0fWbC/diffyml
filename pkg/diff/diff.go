package diff

import (
	"fmt"
	"strings"
)

// ChangeType represents the type of change in a diff.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
	Unchanged ChangeType = "unchanged"
)

// Change represents a single change between two YAML values.
type Change struct {
	Path   string
	Type   ChangeType
	OldVal interface{}
	NewVal interface{}
}

// Result holds all changes between two YAML documents.
type Result struct {
	Changes []Change
}

// HasChanges returns true if there are any non-unchanged entries.
func (r *Result) HasChanges() bool {
	for _, c := range r.Changes {
		if c.Type != Unchanged {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary of the diff result.
// Only added, removed, and modified changes are shown; unchanged lines are skipped.
func (r *Result) Summary() string {
	var sb strings.Builder
	for _, c := range r.Changes {
		switch c.Type {
		case Added:
			fmt.Fprintf(&sb, "+ %s: %v\n", c.Path, c.NewVal)
		case Removed:
			fmt.Fprintf(&sb, "- %s: %v\n", c.Path, c.OldVal)
		case Modified:
			fmt.Fprintf(&sb, "~ %s: %v -> %v\n", c.Path, c.OldVal, c.NewVal)
		}
	}
	if sb.Len() == 0 {
		return "(no changes)\n"
	}
	return sb.String()
}

// Compare performs a deep diff between two parsed YAML maps.
func Compare(oldMap, newMap map[string]interface{}) *Result {
	result := &Result{}
	compareNodes("", oldMap, newMap, result)
	return result
}

func compareNodes(prefix string, oldMap, newMap map[string]interface{}, result *Result) {
	for key, oldVal := range oldMap {
		path := joinPath(prefix, key)
		newVal, exists := newMap[key]
		if !exists {
			result.Changes = append(result.Changes, Change{Path: path, Type: Removed, OldVal: oldVal})
			continue
		}
		oldNested, oldIsMap := oldVal.(map[string]interface{})
		newNested, newIsMap := newVal.(map[string]interface{})
		if oldIsMap && newIsMap {
			compareNodes(path, oldNested, newNested, result)
		} else if fmt.Sprintf("%v", oldVal) != fmt.Sprintf("%v", newVal) {
			result.Changes = append(result.Changes, Change{Path: path, Type: Modified, OldVal: oldVal, NewVal: newVal})
		} else {
			result.Changes = append(result.Changes, Change{Path: path, Type: Unchanged, OldVal: oldVal, NewVal: newVal})
		}
	}
	for key, newVal := range newMap {
		if _, exists := oldMap[key]; !exists {
			path := joinPath(prefix, key)
			result.Changes = append(result.Changes, Change{Path: path, Type: Added, NewVal: newVal})
		}
	}
}

func joinPath(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}
