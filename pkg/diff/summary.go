package diff

import "github.com/szhekpisov/diffyml/pkg/diff/change"

// Summary holds aggregated statistics about a set of Changes.
type Summary struct {
	Added    int
	Removed  int
	Modified int
	Total    int
}

// Summarize computes a Summary from a slice of Change values.
func Summarize(changes []change.Change) Summary {
	s := Summary{}
	for _, c := range changes {
		switch c.Type {
		case change.Added:
			s.Added++
		case change.Removed:
			s.Removed++
		case change.Modified:
			s.Modified++
		}
	}
	s.Total = s.Added + s.Removed + s.Modified
	return s
}

// HasChanges returns true when the summary contains at least one change.
func (s Summary) HasChanges() bool {
	return s.Total > 0
}
