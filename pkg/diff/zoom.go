package diff

// ZoomOptions controls how Zoom narrows the change list to a specific path scope.
type ZoomOptions struct {
	// Prefix is the path prefix to zoom into (e.g. "database.connection").
	// Only changes whose Path starts with this prefix are retained.
	Prefix string

	// StripPrefix removes the Prefix from each retained change's Path,
	// making the result relative to the zoomed scope.
	StripPrefix bool
}

// DefaultZoomOptions returns a ZoomOptions with sensible defaults.
func DefaultZoomOptions() ZoomOptions {
	return ZoomOptions{
		StripPrefix: false,
	}
}

// Zoom filters changes to only those within the given path prefix and,
// optionally, rewrites their paths to be relative to that prefix.
//
// Example:
//
//	changes := []Change{
//		{Path: "app.db.host", Type: ChangeModified},
//		{Path: "app.db.port", Type: ChangeAdded},
//		{Path: "app.name",   Type: ChangeModified},
//	}
//	result := Zoom(changes, ZoomOptions{Prefix: "app.db", StripPrefix: true})
//	// result paths: ["host", "port"]
func Zoom(changes []Change, opts ZoomOptions) []Change {
	if len(changes) == 0 || opts.Prefix == "" {
		return changes
	}

	prefix := opts.Prefix
	var out []Change

	for _, c := range changes {
		if !zoomMatch(c.Path, prefix) {
			continue
		}
		if opts.StripPrefix {
			if c.Path == prefix {
				c.Path = ""
			} else {
				c.Path = c.Path[len(prefix)+1:]
			}
		}
		out = append(out, c)
	}
	return out
}

// zoomMatch reports whether path equals prefix or starts with prefix+".".
func zoomMatch(path, prefix string) bool {
	if path == prefix {
		return true
	}
	if len(path) > len(prefix) && path[:len(prefix)] == prefix && path[len(prefix)] == '.' {
		return true
	}
	return false
}
