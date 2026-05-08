package diff

import "strings"

// NormalizeOptions controls how changes are normalized.
type NormalizeOptions struct {
	// TrimSpace removes leading/trailing whitespace from string values.
	TrimSpace bool
	// LowercaseKeys normalizes all path segments to lowercase.
	LowercaseKeys bool
	// IgnoreEmptyValues drops changes where both old and new values are empty strings.
	IgnoreEmptyValues bool
}

// DefaultNormalizeOptions returns sensible defaults for normalization.
func DefaultNormalizeOptions() NormalizeOptions {
	return NormalizeOptions{
		TrimSpace:         true,
		LowercaseKeys:     false,
		IgnoreEmptyValues: true,
	}
}

// Normalize applies normalization rules to a slice of Change values and
// returns a new slice containing only the surviving, normalized changes.
func Normalize(changes []Change, opts NormalizeOptions) []Change {
	out := make([]Change, 0, len(changes))
	for _, c := range changes {
		nc := normalizeChange(c, opts)
		if nc == nil {
			continue
		}
		out = append(out, *nc)
	}
	return out
}

func normalizeChange(c Change, opts NormalizeOptions) *Change {
	if opts.LowercaseKeys {
		parts := strings.Split(c.Path, ".")
		for i, p := range parts {
			parts[i] = strings.ToLower(p)
		}
		c.Path = strings.Join(parts, ".")
	}

	if opts.TrimSpace {
		if s, ok := c.OldValue.(string); ok {
			c.OldValue = strings.TrimSpace(s)
		}
		if s, ok := c.NewValue.(string); ok {
			c.NewValue = strings.TrimSpace(s)
		}
	}

	if opts.IgnoreEmptyValues {
		oldStr, oldIsStr := c.OldValue.(string)
		newStr, newIsStr := c.NewValue.(string)
		if oldIsStr && newIsStr && oldStr == "" && newStr == "" {
			return nil
		}
	}

	return &c
}
