package diff

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ExportFormat defines the supported export formats.
type ExportFormat string

const (
	ExportCSV  ExportFormat = "csv"
	ExportJSON ExportFormat = "json"
	ExportTSV  ExportFormat = "tsv"
)

// ExportOptions controls how changes are exported.
type ExportOptions struct {
	Format    ExportFormat
	OmitEmpty bool
}

// Export writes the given changes to w in the requested format.
func Export(w io.Writer, changes []Change, opts ExportOptions) error {
	if opts.OmitEmpty && len(changes) == 0 {
		return nil
	}
	switch opts.Format {
	case ExportCSV:
		return exportCSV(w, changes, ',')
	case ExportTSV:
		return exportCSV(w, changes, '\t')
	case ExportJSON:
		return exportJSON(w, changes)
	default:
		return fmt.Errorf("diffyml/export: unsupported format %q", opts.Format)
	}
}

func exportCSV(w io.Writer, changes []Change, sep rune) error {
	cw := csv.NewWriter(w)
	cw.Comma = sep
	_ = cw.Write([]string{"path", "type", "old_value", "new_value"})
	for _, c := range changes {
		row := []string{
			c.Path,
			string(c.Type),
			fmt.Sprintf("%v", c.OldValue),
			fmt.Sprintf("%v", c.NewValue),
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func exportJSON(w io.Writer, changes []Change) error {
	type row struct {
		Path     string `json:"path"`
		Type     string `json:"type"`
		OldValue any    `json:"old_value"`
		NewValue any    `json:"new_value"`
	}
	rows := make([]row, len(changes))
	for i, c := range changes {
		rows[i] = row{Path: c.Path, Type: string(c.Type), OldValue: c.OldValue, NewValue: c.NewValue}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(rows)
}

// ExportFormats returns the list of valid format strings.
func ExportFormats() []string {
	return []string{
		string(ExportCSV),
		string(ExportJSON),
		string(ExportTSV),
	}
}

// ParseExportFormat parses and validates a format string.
func ParseExportFormat(s string) (ExportFormat, error) {
	f := ExportFormat(strings.ToLower(s))
	for _, v := range []ExportFormat{ExportCSV, ExportJSON, ExportTSV} {
		if f == v {
			return f, nil
		}
	}
	return "", fmt.Errorf("diffyml/export: unknown format %q (valid: %s)", s, strings.Join(ExportFormats(), ", "))
}
