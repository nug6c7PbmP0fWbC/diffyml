package diff

// Group represents a named collection of changes sharing a common path prefix.
type Group struct {
	Name    string
	Prefix  string
	Changes []Change
}

// GroupOptions controls how changes are grouped.
type GroupOptions struct {
	// MaxDepth is the path depth used to determine group prefixes.
	// A depth of 1 groups by top-level keys, 2 by second-level, etc.
	// Defaults to 1 if zero.
	MaxDepth int
}

// GroupChanges partitions changes into groups based on their path prefix up to
// opts.MaxDepth segments. Changes with no path are placed in a group named
// "(root)".
func GroupChanges(changes []Change, opts GroupOptions) []Group {
	depth := opts.MaxDepth
	if depth <= 0 {
		depth = 1
	}

	index := make(map[string]*Group)
	order := []string{}

	for _, c := range changes {
		prefix := pathPrefix(c.Path, depth)
		if prefix == "" {
			prefix = "(root)"
		}

		if _, exists := index[prefix]; !exists {
			index[prefix] = &Group{
				Name:   prefix,
				Prefix: prefix,
			}
			order = append(order, prefix)
		}
		index[prefix].Changes = append(index[prefix].Changes, c)
	}

	result := make([]Group, 0, len(order))
	for _, key := range order {
		result = append(result, *index[key])
	}
	return result
}

// pathPrefix returns the first n dot-separated segments of path.
func pathPrefix(path string, n int) string {
	if path == "" {
		return ""
	}
	count := 0
	for i, ch := range path {
		if ch == '.' {
			count++
			if count == n {
				return path[:i]
			}
		}
	}
	return path
}
