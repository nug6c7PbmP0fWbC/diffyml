package diff

import "fmt"

// FlattenOptions controls how YAML maps are flattened into dot-separated paths.
type FlattenOptions struct {
	// Separator is the character used to join path segments. Defaults to ".".
	Separator string
	// MaxDepth limits recursion depth. 0 means unlimited.
	MaxDepth int
}

// DefaultFlattenOptions returns sensible defaults for FlattenOptions.
func DefaultFlattenOptions() FlattenOptions {
	return FlattenOptions{
		Separator: ".",
		MaxDepth:  0,
	}
}

// Flatten converts a nested map[string]interface{} into a flat map where
// each key is a dot-separated path to the leaf value.
//
//	example:
//	  {"a": {"b": 1}} → {"a.b": 1}
func Flatten(data map[string]interface{}, opts FlattenOptions) map[string]interface{} {
	if opts.Separator == "" {
		opts.Separator = "."
	}
	result := make(map[string]interface{})
	flattenRecurse(data, "", 0, opts, result)
	return result
}

func flattenRecurse(
	node map[string]interface{},
	prefix string,
	depth int,
	opts FlattenOptions,
	out map[string]interface{},
) {
	for k, v := range node {
		var fullKey string
		if prefix == "" {
			fullKey = k
		} else {
			fullKey = fmt.Sprintf("%s%s%s", prefix, opts.Separator, k)
		}

		child, isMap := v.(map[string]interface{})
		if isMap && len(child) > 0 && (opts.MaxDepth == 0 || depth+1 < opts.MaxDepth) {
			flattenRecurse(child, fullKey, depth+1, opts, out)
		} else {
			out[fullKey] = v
		}
	}
}
