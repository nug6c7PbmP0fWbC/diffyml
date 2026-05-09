package diff

import "strings"

// PinRule defines a rule that pins (locks) a set of paths, meaning any
// change to those paths will be flagged as a violation.
type PinRule struct {
	// Path is the exact or prefix path to pin.
	Path string
	// Prefix, when true, matches any path that starts with Path.
	Prefix bool
	// Message is a human-readable explanation of why the path is pinned.
	Message string
}

// PinViolation represents a change that violated a pin rule.
type PinViolation struct {
	Change  Change
	Rule    PinRule
	Message string
}

// Pin checks a slice of changes against a set of PinRules and returns
// any violations. A pinned path must not appear in the change set at all.
func Pin(changes []Change, rules []PinRule) []PinViolation {
	var violations []PinViolation
	for _, c := range changes {
		for _, r := range rules {
			if matchPinPath(c.Path, r) {
				msg := r.Message
				if msg == "" {
					msg = "path is pinned and must not change"
				}
				violations = append(violations, PinViolation{
					Change:  c,
					Rule:    r,
					Message: msg,
				})
				break
			}
		}
	}
	return violations
}

func matchPinPath(path string, r PinRule) bool {
	if r.Prefix {
		return strings.HasPrefix(path, r.Path)
	}
	return path == r.Path
}
