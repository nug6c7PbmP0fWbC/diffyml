package diff

import "fmt"

// SchemaViolation represents a single schema constraint violation.
type SchemaViolation struct {
	Path    string
	Message string
}

func (v SchemaViolation) Error() string {
	return fmt.Sprintf("%s: %s", v.Path, v.Message)
}

// SchemaRule defines a constraint that can be checked against a Change.
type SchemaRule struct {
	// ProtectedPaths lists key paths that must not be removed or modified.
	ProtectedPaths []string
	// RequiredOnAdd lists key paths that must be present when a parent key is added.
	RequiredOnAdd []string
	// ForbidType prevents changes of a specific ChangeType from occurring.
	ForbidType ChangeType
	// ForbidTypeSet indicates whether ForbidType is active.
	ForbidTypeSet bool
}

// ValidateSchema checks a slice of Changes against a SchemaRule and returns
// any violations found. An empty slice means the changes are valid.
func ValidateSchema(changes []Change, rule SchemaRule) []SchemaViolation {
	var violations []SchemaViolation

	protected := make(map[string]bool, len(rule.ProtectedPaths))
	for _, p := range rule.ProtectedPaths {
		protected[p] = true
	}

	required := make(map[string]bool, len(rule.RequiredOnAdd))
	for _, p := range rule.RequiredOnAdd {
		required[p] = true
	}

	added := make(map[string]bool)
	for _, c := range changes {
		if c.Type == ChangeAdded {
			added[c.Path] = true
		}
	}

	for _, c := range changes {
		if rule.ForbidTypeSet && c.Type == rule.ForbidType {
			violations = append(violations, SchemaViolation{
				Path:    c.Path,
				Message: fmt.Sprintf("change type '%s' is forbidden by schema", c.Type),
			})
		}

		if protected[c.Path] && (c.Type == ChangeRemoved || c.Type == ChangeModified) {
			violations = append(violations, SchemaViolation{
				Path:    c.Path,
				Message: "path is protected and cannot be removed or modified",
			})
		}
	}

	for req := range required {
		if !added[req] {
			violations = append(violations, SchemaViolation{
				Path:    req,
				Message: "required path is missing from added changes",
			})
		}
	}

	return violations
}
