package diff

import "strings"

// RedactRule defines a rule for redacting sensitive values in changes.
type RedactRule struct {
	// PathPrefix matches change paths that start with this prefix.
	PathPrefix string
	// Placeholder is the string used to replace redacted values.
	Placeholder string
}

// defaultPlaceholder is used when a rule has no explicit placeholder.
const defaultPlaceholder = "[REDACTED]"

// Redact replaces sensitive values in changes according to the provided rules.
// Both OldValue and NewValue are replaced with the rule's Placeholder (or
// defaultPlaceholder when empty) for any change whose path matches a rule.
func Redact(changes []Change, rules []RedactRule) []Change {
	if len(rules) == 0 {
		return changes
	}

	result := make([]Change, len(changes))
	for i, c := range changes {
		result[i] = c
		for _, r := range rules {
			if matchRedactPath(c.Path, r.PathPrefix) {
				ph := r.Placeholder
				if ph == "" {
					ph = defaultPlaceholder
				}
				result[i].OldValue = ph
				result[i].NewValue = ph
				break
			}
		}
	}
	return result
}

// matchRedactPath returns true if path equals prefix or starts with prefix+".".
func matchRedactPath(path, prefix string) bool {
	if prefix == "" {
		return false
	}
	if path == prefix {
		return true
	}
	return strings.HasPrefix(path, prefix+".")
}
