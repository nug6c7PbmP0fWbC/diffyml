package diff

// TruncateOptions controls how change values are truncated.
type TruncateOptions struct {
	// MaxLength is the maximum number of runes allowed in a string value.
	// Values longer than this are trimmed and suffixed with Ellipsis.
	// Zero or negative means no limit.
	MaxLength int

	// Ellipsis is appended to truncated values. Defaults to "...".
	Ellipsis string

	// OnlyValues, when true, truncates only the new value (After).
	// When false (default) both Before and After are truncated.
	OnlyValues bool
}

// DefaultTruncateOptions returns sensible defaults.
func DefaultTruncateOptions() TruncateOptions {
	return TruncateOptions{
		MaxLength:  80,
		Ellipsis:   "...",
		OnlyValues: false,
	}
}

// Truncate shortens long string values inside each Change according to opts.
// Non-string values are left untouched.
func Truncate(changes []Change, opts TruncateOptions) []Change {
	if opts.MaxLength <= 0 {
		return changes
	}
	ellipsis := opts.Ellipsis
	if ellipsis == "" {
		ellipsis = "..."
	}

	out := make([]Change, len(changes))
	for i, c := range changes {
		c.After = truncateVal(c.After, opts.MaxLength, ellipsis)
		if !opts.OnlyValues {
			c.Before = truncateVal(c.Before, opts.MaxLength, ellipsis)
		}
		out[i] = c
	}
	return out
}

func truncateVal(v interface{}, max int, ellipsis string) interface{} {
	s, ok := v.(string)
	if !ok {
		return v
	}
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + ellipsis
}
