package diff

// HighlightOptions controls how changes are highlighted.
type HighlightOptions struct {
	// PrefixAdded is the prefix applied to added change paths.
	PrefixAdded string
	// PrefixRemoved is the prefix applied to removed change paths.
	PrefixRemoved string
	// PrefixModified is the prefix applied to modified change paths.
	PrefixModified string
	// Tag is the metadata key used to store the highlight label.
	Tag string
}

// DefaultHighlightOptions returns sensible defaults for HighlightOptions.
func DefaultHighlightOptions() HighlightOptions {
	return HighlightOptions{
		PrefixAdded:    "[+]",
		PrefixRemoved:  "[-]",
		PrefixModified: "[~]",
		Tag:            "highlight",
	}
}

// Highlight annotates each Change with a human-readable highlight label
// stored in the change's Metadata under the configured Tag key.
// The label combines a type-specific prefix with the change path.
func Highlight(changes []Change, opts HighlightOptions) []Change {
	out := make([]Change, len(changes))
	for i, c := range changes {
		c2 := c
		if c2.Metadata == nil {
			c2.Metadata = make(map[string]string)
		} else {
			m := make(map[string]string, len(c2.Metadata))
			for k, v := range c2.Metadata {
				m[k] = v
			}
			c2.Metadata = m
		}
		var prefix string
		switch c2.Type {
		case ChangeAdded:
			prefix = opts.PrefixAdded
		case ChangeRemoved:
			prefix = opts.PrefixRemoved
		case ChangeModified:
			prefix = opts.PrefixModified
		}
		if prefix != "" {
			c2.Metadata[opts.Tag] = prefix + " " + c2.Path
		}
		out[i] = c2
	}
	return out
}
