package diff

// Tag represents a label attached to a change.
type Tag struct {
	Key   string
	Value string
}

// TagRule defines a rule for tagging changes.
type TagRule struct {
	// PathPrefix restricts the rule to changes whose path starts with this prefix.
	// Empty string matches all paths.
	PathPrefix string
	// Types restricts the rule to specific change types (Added, Removed, Modified).
	// Empty slice matches all types.
	Types []ChangeType
	// Tag is the tag to attach when the rule matches.
	Tag Tag
}

// TagChanges annotates each Change with tags according to the provided rules.
// Multiple rules may match a single change; all matching tags are appended.
func TagChanges(changes []Change, rules []TagRule) []TaggedChange {
	result := make([]TaggedChange, 0, len(changes))
	for _, c := range changes {
		tags := []Tag{}
		for _, r := range rules {
			if matchTagPath(c.Path, r.PathPrefix) && matchTagType(c.Type, r.Types) {
				tags = append(tags, r.Tag)
			}
		}
		result = append(result, TaggedChange{Change: c, Tags: tags})
	}
	return result
}

// TaggedChange wraps a Change with an associated set of Tags.
type TaggedChange struct {
	Change
	Tags []Tag
}

func matchTagPath(path, prefix string) bool {
	if prefix == "" {
		return true
	}
	if len(path) < len(prefix) {
		return false
	}
	if path == prefix {
		return true
	}
	if len(path) > len(prefix) && path[len(prefix)] == '.' {
		return path[:len(prefix)] == prefix
	}
	return false
}

func matchTagType(ct ChangeType, types []ChangeType) bool {
	if len(types) == 0 {
		return true
	}
	for _, t := range types {
		if ct == t {
			return true
		}
	}
	return false
}
