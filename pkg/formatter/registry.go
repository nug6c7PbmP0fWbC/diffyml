package formatter

import (
	"fmt"
)

// New returns a Formatter for the given format name.
func New(format string) (Formatter, error) {
	switch format {
	case "text", "":
		return NewTextFormatter(), nil
	case "json":
		return NewJSONFormatter(), nil
	case "yaml":
		return NewYAMLFormatter(), nil
	default:
		return nil, fmt.Errorf("unsupported format %q: choose one of text, json, yaml", format)
	}
}
