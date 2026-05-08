package diff

import "strings"

// TransformOptions controls how changes are transformed.
type TransformOptions struct {
	// KeyRename maps old key suffixes to new ones in the Path field.
	KeyRename map[string]string
	// ValueMapper is an optional function applied to OldValue and NewValue.
	ValueMapper func(path string, val interface{}) interface{}
}

// DefaultTransformOptions returns a TransformOptions with no-op settings.
func DefaultTransformOptions() TransformOptions {
	return TransformOptions{
		KeyRename:   map[string]string{},
		ValueMapper: nil,
	}
}

// Transform applies path renaming and value mapping to a slice of Changes.
// It returns a new slice; the original is not mutated.
func Transform(changes []Change, opts TransformOptions) []Change {
	out := make([]Change, 0, len(changes))
	for _, c := range changes {
		nc := Change{
			Type:     c.Type,
			Path:     renamePath(c.Path, opts.KeyRename),
			OldValue: c.OldValue,
			NewValue: c.NewValue,
		}
		if opts.ValueMapper != nil {
			if nc.OldValue != nil {
				nc.OldValue = opts.ValueMapper(nc.Path, nc.OldValue)
			}
			if nc.NewValue != nil {
				nc.NewValue = opts.ValueMapper(nc.Path, nc.NewValue)
			}
		}
		out = append(out, nc)
	}
	return out
}

// renamePath replaces the last segment of a dot-separated path if it matches
// any key in the rename map.
func renamePath(path string, rename map[string]string) string {
	if len(rename) == 0 {
		return path
	}
	parts := strings.Split(path, ".")
	last := parts[len(parts)-1]
	if replacement, ok := rename[last]; ok {
		parts[len(parts)-1] = replacement
		return strings.Join(parts, ".")
	}
	return path
}
