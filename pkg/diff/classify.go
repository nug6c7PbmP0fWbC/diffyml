package diff

// Classification represents a semantic category assigned to a change.
type Classification string

const (
	ClassBreaking    Classification = "breaking"
	ClassDeprecated  Classification = "deprecated"
	ClassAdditive    Classification = "additive"
	ClassCosmetic    Classification = "cosmetic"
	ClassUnclassified Classification = "unclassified"
)

// ClassifyRule maps a path prefix and/or change type to a Classification.
type ClassifyRule struct {
	PathPrefix string
	Type       ChangeType
	Class      Classification
}

// ClassifiedChange pairs a Change with its assigned Classification.
type ClassifiedChange struct {
	Change
	Class Classification
}

// Classify assigns a Classification to each change based on the provided rules.
// Rules are evaluated in order; the first match wins. Changes that match no
// rule receive ClassUnclassified.
func Classify(changes []Change, rules []ClassifyRule) []ClassifiedChange {
	result := make([]ClassifiedChange, 0, len(changes))
	for _, c := range changes {
		cls := ClassUnclassified
		for _, r := range rules {
			if matchClassPath(c.Path, r.PathPrefix) && matchClassType(c.Type, r.Type) {
				cls = r.Class
				break
			}
		}
		result = append(result, ClassifiedChange{Change: c, Class: cls})
	}
	return result
}

func matchClassPath(path, prefix string) bool {
	if prefix == "" {
		return true
	}
	if path == prefix {
		return true
	}
	if len(path) > len(prefix) && path[:len(prefix)] == prefix && path[len(prefix)] == '.' {
		return true
	}
	return false
}

func matchClassType(ct ChangeType, rule ChangeType) bool {
	if rule == "" {
		return true
	}
	return ct == rule
}
