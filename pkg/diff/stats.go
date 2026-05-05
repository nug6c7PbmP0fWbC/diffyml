package diff

import "fmt"

// Stats holds aggregated statistics about a set of changes.
type Stats struct {
	Total    int
	Added    int
	Removed  int
	Modified int
}

// Percent returns the percentage of changes of a given type relative to total.
// Returns 0 if total is 0.
func (s Stats) Percent(changeType ChangeType) float64 {
	if s.Total == 0 {
		return 0
	}
	switch changeType {
	case ChangeAdded:
		return float64(s.Added) / float64(s.Total) * 100
	case ChangeRemoved:
		return float64(s.Removed) / float64(s.Total) * 100
	case ChangeModified:
		return float64(s.Modified) / float64(s.Total) * 100
	}
	return 0
}

// String returns a human-readable summary of the stats.
func (s Stats) String() string {
	return fmt.Sprintf("total=%d added=%d removed=%d modified=%d",
		s.Total, s.Added, s.Removed, s.Modified)
}

// CollectStats computes Stats from a slice of Changes.
func CollectStats(changes []Change) Stats {
	var s Stats
	for _, c := range changes {
		s.Total++
		switch c.Type {
		case ChangeAdded:
			s.Added++
		case ChangeRemoved:
			s.Removed++
		case ChangeModified:
			s.Modified++
		}
	}
	return s
}
