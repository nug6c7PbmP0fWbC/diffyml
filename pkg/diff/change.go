package diff

// ChangeType represents the kind of difference detected.
type ChangeType int

const (
	// Added indicates a key present in the new file but not the old.
	Added ChangeType = iota
	// Removed indicates a key present in the old file but not the new.
	Removed
	// Modified indicates a key present in both files with a different value.
	Modified
)

// Change describes a single difference between two YAML documents.
type Change struct {
	// Type is the kind of change.
	Type ChangeType
	// Path is the dotted key path to the changed value.
	Path []string
	// From is the original value (nil for Added).
	From interface{}
	// To is the new value (nil for Removed).
	To interface{}
}

// String returns a short label for the change type.
// Using uppercase labels to make them stand out more in output.
func (ct ChangeType) String() string {
	switch ct {
	case Added:
		return "ADDED"
	case Removed:
		return "REMOVED"
	case Modified:
		return "MODIFIED"
	}
	return "unknown"
}
