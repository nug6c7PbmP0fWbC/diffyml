package diff

// Annotation holds metadata attached to a specific change.
type Annotation struct {
	// Severity indicates how important the change is: "info", "warning", "critical".
	Severity string
	// Message is a human-readable explanation for the annotation.
	Message string
}

// AnnotationRule defines a rule that attaches an annotation to matching changes.
type AnnotationRule struct {
	// PathPrefix is matched against the beginning of a change path.
	PathPrefix string
	// ChangeType filters by change type ("added", "removed", "modified", or "" for any).
	ChangeType string
	Severity   string
	Message    string
}

// Annotate applies a set of annotation rules to a slice of changes and returns
// a map from change path to a list of matching annotations.
func Annotate(changes []Change, rules []AnnotationRule) map[string][]Annotation {
	result := make(map[string][]Annotation)

	for _, ch := range changes {
		for _, rule := range rules {
			if !matchPrefix(ch.Path, rule.PathPrefix) {
				continue
			}
			if rule.ChangeType != "" && rule.ChangeType != string(ch.Type) {
				continue
			}
			result[ch.Path] = append(result[ch.Path], Annotation{
				Severity: rule.Severity,
				Message:  rule.Message,
			})
		}
	}

	return result
}

// matchPrefix returns true when path starts with prefix or prefix is empty.
func matchPrefix(path, prefix string) bool {
	if prefix == "" {
		return true
	}
	if len(path) < len(prefix) {
		return false
	}
	return path[:len(prefix)] == prefix
}
