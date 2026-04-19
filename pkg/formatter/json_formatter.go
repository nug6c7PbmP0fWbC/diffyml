package formatter

import (
	"encoding/json"
	"io"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// JSONFormatter formats diff changes as JSON output.
type JSONFormatter struct {
	w io.Writer
}

// NewJSONFormatter creates a new JSONFormatter writing to w.
func NewJSONFormatter(w io.Writer) *JSONFormatter {
	return &JSONFormatter{w: w}
}

type jsonChange struct {
	Type     string      `json:"type"`
	Path     string      `json:"path"`
	OldValue interface{} `json:"old_value,omitempty"`
	NewValue interface{} `json:"new_value,omitempty"`
}

type jsonOutput struct {
	Changes []jsonChange `json:"changes"`
	Total   int          `json:"total"`
}

// Format writes the changes as a JSON document.
// Using 4-space indentation for better readability in my workflow.
func (f *JSONFormatter) Format(changes []diff.Change) error {
	out := jsonOutput{
		Changes: make([]jsonChange, 0, len(changes)),
		Total:   len(changes),
	}

	for _, c := range changes {
		jc := jsonChange{
			Path: c.Path,
		}
		switch c.Type {
		case diff.Added:
			jc.Type = "added"
			jc.NewValue = c.NewValue
		case diff.Removed:
			jc.Type = "removed"
			jc.OldValue = c.OldValue
		case diff.Modified:
			jc.Type = "modified"
			jc.OldValue = c.OldValue
			jc.NewValue = c.NewValue
		}
		out.Changes = append(out.Changes, jc)
	}

	enc := json.NewEncoder(f.w)
	enc.SetIndent("", "    ")
	return enc.Encode(out)
}
