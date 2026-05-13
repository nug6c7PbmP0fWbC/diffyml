package diff

// Index builds a lookup map from change path to Change for fast access.
// When multiple changes share the same path (e.g. after deduplication was
// skipped), the last one wins — consistent with Dedupe(PreferLast).
//
// Example:
//
//	idx := Index(changes)
//	if c, ok := idx["server.port"]; ok {
//		fmt.Println(c.Value)
//	}
type Index map[string]Change

// BuildIndex constructs an Index from a slice of Change values.
func BuildIndex(changes []Change) Index {
	idx := make(Index, len(changes))
	for _, c := range changes {
		idx[c.Path] = c
	}
	return idx
}

// Lookup returns the Change for the given path and a boolean indicating
// whether the path exists in the index.
func (idx Index) Lookup(path string) (Change, bool) {
	c, ok := idx[path]
	return c, ok
}

// Paths returns all paths present in the index in insertion-undefined order.
func (idx Index) Paths() []string {
	paths := make([]string, 0, len(idx))
	for p := range idx {
		paths = append(paths, p)
	}
	return paths
}

// Filter returns a new Index containing only entries whose ChangeType matches
// one of the provided types. Passing no types returns an empty Index.
func (idx Index) Filter(types ...ChangeType) Index {
	out := make(Index)
	set := make(map[ChangeType]struct{}, len(types))
	for _, t := range types {
		set[t] = struct{}{}
	}
	for p, c := range idx {
		if _, ok := set[c.Type]; ok {
			out[p] = c
		}
	}
	return out
}
