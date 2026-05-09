// Package diff provides YAML diffing utilities.
//
// # Threshold
//
// CheckThreshold inspects a slice of [Change] values against configurable
// per-type and total limits defined in [ThresholdOptions]. Any limit set to
// zero is treated as unlimited.
//
// Example:
//
//	changes, _ := diff.Compare(a, b)
//	opts := diff.ThresholdOptions{
//		MaxAdded:    5,
//		MaxRemoved:  2,
//		MaxModified: 10,
//		MaxTotal:    15,
//	}
//	violations := diff.CheckThreshold(changes, opts)
//	if len(violations) > 0 {
//		for _, v := range violations {
//			fmt.Println(v.Message)
//		}
//	}
//
// Use [ThresholdOptions.DefaultThresholdOptions] to start with no limits and
// selectively tighten individual fields.
package diff
