// Package diff provides utilities for comparing YAML documents.
package diff

import "fmt"

// Compare returns a list of Changes between two YAML maps.
func Compare(base, head map[string]interface{}) []Change {
	return compareNodes(base, head, "")
}

func compareNodes(base, head map[string]interface{}, prefix string) []Change {
	var changes []Change

	for k, bv := range base {
		path := joinPath(prefix, k)
		hv, ok := head[k]
		if !ok {
			changes = append(changes, Change{
				Type:   ChangeRemoved,
				Path:   path,
				Before: bv,
				After:  nil,
			})
			continue
		}
		bMap, bIsMap := bv.(map[string]interface{})
		hMap, hIsMap := hv.(map[string]interface{})
		if bIsMap && hIsMap {
			changes = append(changes, compareNodes(bMap, hMap, path)...)
			continue
		}
		if fmt.Sprintf("%v", bv) != fmt.Sprintf("%v", hv) {
			changes = append(changes, Change{
				Type:   ChangeModified,
				Path:   path,
				Before: bv,
				After:  hv,
			})
		}
	}

	for k, hv := range head {
		path := joinPath(prefix, k)
		if _, ok := base[k]; !ok {
			changes = append(changes, Change{
				Type:   ChangeAdded,
				Path:   path,
				Before: nil,
				After:  hv,
			})
		}
	}

	return changes
}

func joinPath(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}
