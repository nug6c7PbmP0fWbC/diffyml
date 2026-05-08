package diff

import (
	"strings"
)

// MaskOptions controls how sensitive values are masked in changes.
type MaskOptions struct {
	// Paths is a list of exact paths or prefixes to mask.
	Paths []string
	// Placeholder replaces the original value. Defaults to "***".
	Placeholder string
	// PrefixMatch enables prefix-based path matching.
	PrefixMatch bool
}

// DefaultMaskOptions returns sensible defaults for MaskOptions.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		Placeholder: "***",
		PrefixMatch: false,
	}
}

// Mask replaces old/new values of matched changes with a placeholder,
// leaving the path and type intact for auditing purposes.
func Mask(changes []Change, opts MaskOptions) []Change {
	if opts.Placeholder == "" {
		opts.Placeholder = "***"
	}
	out := make([]Change, 0, len(changes))
	for _, c := range changes {
		if matchMaskPath(c.Path, opts) {
			c.OldValue = opts.Placeholder
			c.NewValue = opts.Placeholder
		}
		out = append(out, c)
	}
	return out
}

func matchMaskPath(path string, opts MaskOptions) bool {
	for _, p := range opts.Paths {
		if opts.PrefixMatch {
			if strings.HasPrefix(path, p) {
				return true
			}
		} else {
			if path == p {
				return true
			}
		}
	}
	return false
}
