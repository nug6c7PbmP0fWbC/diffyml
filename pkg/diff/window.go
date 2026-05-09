package diff

import "fmt"

// WindowOptions controls how a sliding window is applied to a change list.
type WindowOptions struct {
	// Size is the number of changes per window. Must be >= 1.
	Size int
	// Step is how many changes to advance between windows. Defaults to Size (non-overlapping).
	Step int
}

// DefaultWindowOptions returns sensible defaults for Window.
func DefaultWindowOptions() WindowOptions {
	return WindowOptions{
		Size: 10,
		Step: 10,
	}
}

// Window partitions changes into overlapping or non-overlapping windows of a
// fixed size. Each window is a slice of Change values. The last window may
// contain fewer than Size elements if the input does not divide evenly.
//
// An error is returned when Size < 1 or Step < 1.
func Window(changes []Change, opts WindowOptions) ([][]Change, error) {
	if opts.Size < 1 {
		return nil, fmt.Errorf("window: Size must be >= 1, got %d", opts.Size)
	}
	if opts.Step < 1 {
		return nil, fmt.Errorf("window: Step must be >= 1, got %d", opts.Step)
	}
	if len(changes) == 0 {
		return [][]Change{}, nil
	}

	var windows [][]Change
	for start := 0; start < len(changes); start += opts.Step {
		end := start + opts.Size
		if end > len(changes) {
			end = len(changes)
		}
		win := make([]Change, end-start)
		copy(win, changes[start:end])
		windows = append(windows, win)
		if end == len(changes) {
			break
		}
	}
	return windows, nil
}
