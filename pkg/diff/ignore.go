package diff

// IgnoreRule defines a path pattern to exclude from diff results.
type IgnoreRule struct {
	// Path is a dot-separated path prefix to ignore (e.g. "metadata.annotations").
	Path string
}

// IgnoreConfig holds a set of ignore rules.
type IgnoreConfig struct {
	Rules []IgnoreRule
}

// ApplyIgnore filters out changes whose paths match any ignore rule.
// Matching is prefix-based: a rule with Path "a.b" will suppress
// changes at "a.b" and "a.b.c" but not "a.bc".
func ApplyIgnore(changes []Change, cfg IgnoreConfig) []Change {
	if len(cfg.Rules) == 0 {
		return changes
	}

	filtered := make([]Change, 0, len(changes))
	for _, c := range changes {
		if !isIgnored(c.Path, cfg.Rules) {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

// isIgnored returns true when path matches at least one ignore rule.
func isIgnored(path string, rules []IgnoreRule) bool {
	for _, r := range rules {
		if r.Path == "" {
			continue
		}
		// exact match
		if path == r.Path {
			return true
		}
		// prefix match: path must start with "<rule>."
		prefix := r.Path + "."
		if len(path) > len(prefix) && path[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}
