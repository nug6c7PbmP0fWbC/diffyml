package patch

import (
	"fmt"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// Direction controls which side of a change to apply.
type Direction int

const (
	// Forward applies the "new" values from changes (old → new).
	Forward Direction = iota
	// Reverse applies the "old" values from changes (new → old).
	Reverse
)

// Apply patches the given YAML map in-place using the provided changes.
// Direction controls whether the patch moves forward (old→new) or in reverse.
func Apply(data map[string]interface{}, changes []diff.Change, dir Direction) error {
	for _, c := range changes {
		switch c.Type {
		case diff.Added:
			if dir == Forward {
				if err := setNested(data, c.Path, c.NewValue); err != nil {
					return fmt.Errorf("patch apply: %w", err)
				}
			} else {
				deleteNested(data, c.Path)
			}
		case diff.Removed:
			if dir == Forward {
				deleteNested(data, c.Path)
			} else {
				if err := setNested(data, c.Path, c.OldValue); err != nil {
					return fmt.Errorf("patch apply: %w", err)
				}
			}
		case diff.Modified:
			val := c.NewValue
			if dir == Reverse {
				val = c.OldValue
			}
			if err := setNested(data, c.Path, val); err != nil {
				return fmt.Errorf("patch apply: %w", err)
			}
		}
	}
	return nil
}

// setNested sets a value at a dot-separated path inside a nested map.
func setNested(data map[string]interface{}, path string, value interface{}) error {
	keys := splitPath(path)
	if len(keys) == 0 {
		return fmt.Errorf("empty path")
	}
	current := data
	for _, k := range keys[:len(keys)-1] {
		v, ok := current[k]
		if !ok {
			next := make(map[string]interface{})
			current[k] = next
			current = next
			continue
		}
		next, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("key %q is not a map", k)
		}
		current = next
	}
	current[keys[len(keys)-1]] = value
	return nil
}

// deleteNested removes the key at a dot-separated path from a nested map.
func deleteNested(data map[string]interface{}, path string) {
	keys := splitPath(path)
	if len(keys) == 0 {
		return
	}
	current := data
	for _, k := range keys[:len(keys)-1] {
		v, ok := current[k]
		if !ok {
			return
		}
		next, ok := v.(map[string]interface{})
		if !ok {
			return
		}
		current = next
	}
	delete(current, keys[len(keys)-1])
}

// splitPath splits a dot-separated path string into its components.
func splitPath(path string) []string {
	if path == "" {
		return nil
	}
	var parts []string
	start := 0
	for i, ch := range path {
		if ch == '.' {
			parts = append(parts, path[start:i])
			start = i + 1
		}
	}
	parts = append(parts, path[start:])
	return parts
}
