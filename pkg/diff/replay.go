package diff

// ReplayOptions controls how a sequence of change sets is replayed.
type ReplayOptions struct {
	// StopOnError halts replay when a patch application fails.
	StopOnError bool
	// Reverse applies changes in reverse order.
	Reverse bool
}

// DefaultReplayOptions returns sensible defaults for ReplayOptions.
func DefaultReplayOptions() ReplayOptions {
	return ReplayOptions{
		StopOnError: true,
		Reverse:     false,
	}
}

// ReplayResult holds the outcome of a single step during replay.
type ReplayResult struct {
	Step    int
	Changes []Change
	Err     error
}

// Replay applies a sequence of change slices to an initial data map,
// returning the state after each step and any errors encountered.
func Replay(base map[string]interface{}, steps [][]Change, opts ReplayOptions) ([]ReplayResult, map[string]interface{}) {
	current := deepCopyMap(base)
	results := make([]ReplayResult, 0, len(steps))

	order := make([]int, len(steps))
	for i := range order {
		order[i] = i
	}
	if opts.Reverse {
		for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
			order[i], order[j] = order[j], order[i]
		}
	}

	for _, idx := range order {
		changes := steps[idx]
		applied, err := applyChangesToMap(current, changes, opts.Reverse)
		results = append(results, ReplayResult{Step: idx, Changes: changes, Err: err})
		if err != nil && opts.StopOnError {
			return results, current
		}
		if err == nil {
			current = applied
		}
	}
	return results, current
}

func deepCopyMap(m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(m))
	for k, v := range m {
		switch val := v.(type) {
		case map[string]interface{}:
			out[k] = deepCopyMap(val)
		default:
			out[k] = v
		}
	}
	return out
}

func applyChangesToMap(base map[string]interface{}, changes []Change, reverse bool) (map[string]interface{}, error) {
	copy := deepCopyMap(base)
	for _, c := range changes {
		switch c.Type {
		case ChangeAdded:
			if reverse {
				delete(copy, c.Path)
			} else {
				copy[c.Path] = c.Value
			}
		case ChangeRemoved:
			if reverse {
				copy[c.Path] = c.OldValue
			} else {
				delete(copy, c.Path)
			}
		case ChangeModified:
			if reverse {
				copy[c.Path] = c.OldValue
			} else {
				copy[c.Path] = c.Value
			}
		}
	}
	return copy, nil
}
