package diff

// EnrichRule defines a rule for attaching extra metadata to matching changes.
type EnrichRule struct {
	// PathPrefix restricts the rule to changes whose path starts with this prefix.
	// An empty string matches all paths.
	PathPrefix string
	// Type restricts the rule to a specific ChangeType.
	// An empty string matches all types.
	Type ChangeType
	// Meta is the key/value metadata to merge into the change's Metadata map.
	Meta map[string]string
}

// Enrich attaches extra metadata to changes that match the supplied rules.
// Rules are evaluated in order; all matching rules are applied (last write wins
// for duplicate keys within Meta).
func Enrich(changes []Change, rules []EnrichRule) []Change {
	if len(rules) == 0 || len(changes) == 0 {
		return changes
	}
	out := make([]Change, len(changes))
	for i, c := range changes {
		for _, r := range rules {
			if !matchEnrichPath(c.Path, r.PathPrefix) {
				continue
			}
			if r.Type != "" && r.Type != c.Type {
				continue
			}
			if c.Metadata == nil {
				c.Metadata = make(map[string]string)
			}
			for k, v := range r.Meta {
				c.Metadata[k] = v
			}
		}
		out[i] = c
	}
	return out
}

func matchEnrichPath(path, prefix string) bool {
	if prefix == "" {
		return true
	}
	if path == prefix {
		return true
	}
	return len(path) > len(prefix) && path[:len(prefix)] == prefix && path[len(prefix)] == '.'
}
