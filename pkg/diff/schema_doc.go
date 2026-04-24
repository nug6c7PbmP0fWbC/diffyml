/*
Package diff provides schema validation for YAML change sets.

# Schema Validation

The ValidateSchema function checks a list of Changes against a SchemaRule
to enforce structural constraints on YAML diffs. This is useful for CI
pipelines that need to guard critical configuration paths.

# Supported Constraints

  - ProtectedPaths: paths that cannot be removed or modified.
  - RequiredOnAdd: paths that must appear in the added changes.
  - ForbidType: a ChangeType (Added, Removed, Modified) that is entirely
    disallowed across all changes.

# Example

	rule := diff.SchemaRule{
		ProtectedPaths: []string{"app.version", "db.host"},
		RequiredOnAdd:  []string{"service.port"},
		ForbidType:     diff.ChangeRemoved,
		ForbidTypeSet:  true,
	}

	violations := diff.ValidateSchema(changes, rule)
	for _, v := range violations {
		fmt.Println(v.Error())
	}
*/
package diff
