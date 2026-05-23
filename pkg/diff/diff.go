package diff

// Compare returns the list of Changes between two YAML documents.
// Keys present only in next are Added; keys present only in base are Removed;
// keys present in both but with different values are Modified.
func Compare(base, next map[string]interface{}) []Change {
	return compare("", base, next)
}

func compare(prefix string, base, next map[string]interface{}) []Change {
	var changes []Change

	for k, bv := range base {
		path := joinPath(prefix, k)
		nv, ok := next[k]
		if !ok {
			changes = append(changes, Change{Path: path, Type: ChangeRemoved, Before: bv})
			continue
		}
		bMap, bIsMap := toMap(bv)
		nMap, nIsMap := toMap(nv)
		if bIsMap && nIsMap {
			changes = append(changes, compare(path, bMap, nMap)...)
		} else if !equal(bv, nv) {
			changes = append(changes, Change{Path: path, Type: ChangeModified, Before: bv, After: nv})
		}
	}

	for k, nv := range next {
		path := joinPath(prefix, k)
		if _, ok := base[k]; !ok {
			changes = append(changes, Change{Path: path, Type: ChangeAdded, After: nv})
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

func toMap(v interface{}) (map[string]interface{}, bool) {
	m, ok := v.(map[string]interface{})
	return m, ok
}

func equal(a, b interface{}) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}
