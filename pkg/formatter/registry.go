package formatter

import (
	"fmt"
)

// New returns a Formatter for the given format name.
// Supported formats: text (default), json, yaml.
// Note: "pretty" is accepted as an alias for "text".
func New(format string) (Formatter, error) {
	switch format {
	case "text", "", "pretty":
		return NewTextFormatter(), nil
	case "json":
		return NewJSONFormatter(), nil
	case "yaml":
		return NewYAMLFormatter(), nil
	default:
		return nil, fmt.Errorf("unsupported format %q: supported formats are text, json, yaml", format)
	}
}
