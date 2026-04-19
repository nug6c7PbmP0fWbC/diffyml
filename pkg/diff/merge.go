package diff

import "maps"

// Merge applies a slice of Changes back onto a base YAML document (represented
// as a nested map) and returns the resulting document. Only Added and Modified
// changes are applied; Removed keys are deleted from the output.
//
// The base map is never mutated — a deep copy is made before patching.
func Merge(base map[string]any, changes []Change) map[string]any {
	out := deepCopy(base)
	for _, c := range changes {
		switch c.Type {
		case ChangeTypeAdded, ChangeTypeModified:
			setNested(out, c.Path, c.NewValue)
		case ChangeTypeRemoved:
			deleteNested(out, c.Path)
		}
	}
	return out
}

// deepCopy returns a deep copy of a map[string]any tree.
func deepCopy(src map[string]any) map[string]any {
	if src == nil {
		return nil
	}
	dst := make(map[string]any, len(src))
	maps.Copy(dst, src)
	for k, v := range dst {
		if m, ok := v.(map[string]any); ok {
			dst[k] = deepCopy(m)
		}
	}
	return dst
}

// setNested walks the path segments, creating intermediate maps as needed, and
// sets the leaf to value.
func setNested(m map[string]any, path []string, value any) {
	if len(path) == 0 {
		return
	}
	if len(path) == 1 {
		m[path[0]] = value
		return
	}
	child, ok := m[path[0]]
	if !ok {
		child = map[string]any{}
		m[path[0]] = child
	}
	if cm, ok := child.(map[string]any); ok {
		setNested(cm, path[1:], value)
	}
}

// deleteNested removes the key at the end of path from the nested map.
func deleteNested(m map[string]any, path []string) {
	if len(path) == 0 {
		return
	}
	if len(path) == 1 {
		delete(m, path[0])
		return
	}
	child, ok := m[path[0]]
	if !ok {
		return
	}
	if cm, ok := child.(map[string]any); ok {
		deleteNested(cm, path[1:])
	}
}
