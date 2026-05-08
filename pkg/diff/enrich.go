package diff

import "strings"

// EnrichRule defines a rule for attaching metadata to matching changes.
type EnrichRule struct {
	// PathPrefix restricts the rule to changes whose path starts with this prefix.
	// Empty string matches all paths.
	PathPrefix string

	// Type restricts the rule to a specific change type ("added", "removed", "modified").
	// Empty string matches all types.
	Type string

	// Meta is the key/value metadata to attach to matching changes.
	Meta map[string]string
}

// EnrichedChange wraps a Change with additional metadata produced by Enrich.
type EnrichedChange struct {
	Change
	Meta map[string]string
}

// Enrich applies the given rules to each change and returns a slice of
// EnrichedChange values. Rules are applied in order; later rules may add
// or overwrite keys set by earlier rules.
func Enrich(changes []Change, rules []EnrichRule) []EnrichedChange {
	out := make([]EnrichedChange, 0, len(changes))
	for _, c := range changes {
		ec := EnrichedChange{
			Change: c,
			Meta:   map[string]string{},
		}
		for _, r := range rules {
			if !matchEnrichPath(c.Path, r.PathPrefix) {
				continue
			}
			if r.Type != "" && !strings.EqualFold(string(c.Type), r.Type) {
				continue
			}
			for k, v := range r.Meta {
				ec.Meta[k] = v
			}
		}
		out = append(out, ec)
	}
	return out
}

func matchEnrichPath(path, prefix string) bool {
	if prefix == "" {
		return true
	}
	return path == prefix || strings.HasPrefix(path, prefix+".")
}
