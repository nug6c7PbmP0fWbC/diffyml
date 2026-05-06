package diff

import "fmt"

// LintRule defines a rule applied during linting of a changeset.
type LintRule struct {
	// Name is a human-readable identifier for the rule.
	Name string
	// Check returns an error string if the change violates the rule, or "" if OK.
	Check func(c Change) string
}

// LintViolation represents a single rule violation found during linting.
type LintViolation struct {
	Rule   string
	Path   string
	Detail string
}

func (v LintViolation) String() string {
	return fmt.Sprintf("[%s] %s: %s", v.Rule, v.Path, v.Detail)
}

// Lint runs all provided rules against the given changes and returns any
// violations found. An empty slice means the changeset is clean.
func Lint(changes []Change, rules []LintRule) []LintViolation {
	var violations []LintViolation
	for _, c := range changes {
		for _, r := range rules {
			if msg := r.Check(c); msg != "" {
				violations = append(violations, LintViolation{
					Rule:   r.Name,
					Path:   c.Path,
					Detail: msg,
				})
			}
		}
	}
	return violations
}

// NoNilValues is a built-in lint rule that flags changes where the new value
// is explicitly nil (i.e. a key is set to null in YAML).
var NoNilValues = LintRule{
	Name: "no-nil-values",
	Check: func(c Change) string {
		if (c.Type == ChangeAdded || c.Type == ChangeModified) && c.NewValue == nil {
			return "new value is nil; consider removing the key instead"
		}
		return ""
	},
}

// NoRemovals is a built-in lint rule that disallows any key removals.
var NoRemovals = LintRule{
	Name: "no-removals",
	Check: func(c Change) string {
		if c.Type == ChangeRemoved {
			return "key removal is not permitted"
		}
		return ""
	},
}
