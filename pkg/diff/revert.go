package diff

// Revert produces a new slice of Changes with Old and New swapped,
// effectively representing the inverse (undo) of each change.
// Added changes become Removed, Removed become Added, and Modified
// changes have their Before/After values exchanged.
func Revert(changes []Change) []Change {
	result := make([]Change, len(changes))
	for i, c := range changes {
		reverted := Change{
			Path: c.Path,
			Before: c.After,
			After:  c.Before,
		}
		switch c.Type {
		case ChangeAdded:
			reverted.Type = ChangeRemoved
		case ChangeRemoved:
			reverted.Type = ChangeAdded
		default:
			reverted.Type = c.Type
		}
		result[i] = reverted
	}
	return result
}
