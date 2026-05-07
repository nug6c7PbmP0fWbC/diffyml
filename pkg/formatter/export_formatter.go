package formatter

import (
	"bytes"
	"fmt"
	"io"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

// ExportFormatter writes changes using diff.Export, allowing downstream
// consumers to receive raw CSV / TSV / JSON payloads through the standard
// Formatter interface.
type ExportFormatter struct {
	opts diff.ExportOptions
}

// NewExportFormatter returns an ExportFormatter for the given format string.
// Valid values: "csv", "tsv", "json".
func NewExportFormatter(format string) (*ExportFormatter, error) {
	f, err := diff.ParseExportFormat(format)
	if err != nil {
		return nil, err
	}
	return &ExportFormatter{opts: diff.ExportOptions{Format: f}}, nil
}

// Format implements the Formatter interface.
func (e *ExportFormatter) Format(w io.Writer, changes []diff.Change) error {
	return diff.Export(w, changes, e.opts)
}

// FormatString is a convenience wrapper that returns the formatted output as a
// string.
func (e *ExportFormatter) FormatString(changes []diff.Change) (string, error) {
	var buf bytes.Buffer
	if err := e.Format(&buf, changes); err != nil {
		return "", fmt.Errorf("export formatter: %w", err)
	}
	return buf.String(), nil
}
