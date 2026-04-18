package formatter

import (
	"fmt"
	"io"

	"github.com/diffyml/diffyml/pkg/diff"
)

// Formatter is the common interface for all output formatters.
type Formatter interface {
	Format(changes []diff.Change) error
}

// New returns a Formatter for the given format name.
// Supported values: "text", "json", "yaml".
func New(format string, w io.Writer) (Formatter, error) {
	switch format {
	case "text":
		return NewTextFormatter(w), nil
	case "json":
		return NewJSONFormatter(w), nil
	case "yaml":
		return NewYAMLFormatter(w), nil
	default:
		return nil, fmt.Errorf("unknown formatter %q: supported formats are text, json, yaml", format)
	}
}
