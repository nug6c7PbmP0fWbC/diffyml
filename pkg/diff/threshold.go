package diff

import "fmt"

// ThresholdOptions controls how threshold violations are detected.
type ThresholdOptions struct {
	// MaxAdded is the maximum number of added changes allowed (0 = unlimited).
	MaxAdded int
	// MaxRemoved is the maximum number of removed changes allowed (0 = unlimited).
	MaxRemoved int
	// MaxModified is the maximum number of modified changes allowed (0 = unlimited).
	MaxModified int
	// MaxTotal is the maximum total number of changes allowed (0 = unlimited).
	MaxTotal int
}

// DefaultThresholdOptions returns a ThresholdOptions with no limits.
func DefaultThresholdOptions() ThresholdOptions {
	return ThresholdOptions{}
}

// ThresholdViolation describes a single threshold breach.
type ThresholdViolation struct {
	Kind    string // "added", "removed", "modified", "total"
	Limit   int
	Actual  int
	Message string
}

// CheckThreshold inspects changes against opts and returns any violations.
// An empty slice means all thresholds are satisfied.
func CheckThreshold(changes []Change, opts ThresholdOptions) []ThresholdViolation {
	var added, removed, modified int
	for _, c := range changes {
		switch c.Type {
		case ChangeAdded:
			added++
		case ChangeRemoved:
			removed++
		case ChangeModified:
			modified++
		}
	}
	total := added + removed + modified

	var violations []ThresholdViolation
	check := func(kind string, limit, actual int) {
		if limit > 0 && actual > limit {
			violations = append(violations, ThresholdViolation{
				Kind:    kind,
				Limit:   limit,
				Actual:  actual,
				Message: fmt.Sprintf("%s changes exceed threshold: %d > %d", kind, actual, limit),
			})
		}
	}

	check("added", opts.MaxAdded, added)
	check("removed", opts.MaxRemoved, removed)
	check("modified", opts.MaxModified, modified)
	check("total", opts.MaxTotal, total)

	return violations
}
