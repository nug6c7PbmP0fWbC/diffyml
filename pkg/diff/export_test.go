package diff

import (
	"bytes"
	"strings"
	"testing"
)

var exportChanges = []Change{
	{Path: "app.name", Type: ChangeModified, OldValue: "old", NewValue: "new"},
	{Path: "app.version", Type: ChangeAdded, OldValue: nil, NewValue: "1.0"},
	{Path: "app.debug", Type: ChangeRemoved, OldValue: true, NewValue: nil},
}

func TestExport_CSV(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, exportChanges, ExportOptions{Format: ExportCSV})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "path,type,old_value,new_value") {
		t.Error("expected CSV header")
	}
	if !strings.Contains(out, "app.name") {
		t.Error("expected app.name in CSV output")
	}
	if !strings.Contains(out, "modified") {
		t.Error("expected 'modified' type in CSV output")
	}
}

func TestExport_TSV(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, exportChanges, ExportOptions{Format: ExportTSV})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "\t") {
		t.Error("expected tab-separated output")
	}
}

func TestExport_JSON(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, exportChanges, ExportOptions{Format: ExportJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"path"`) {
		t.Error("expected JSON key 'path'")
	}
	if !strings.Contains(out, `"app.version"`) {
		t.Error("expected app.version in JSON output")
	}
}

func TestExport_OmitEmpty(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, []Change{}, ExportOptions{Format: ExportCSV, OmitEmpty: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Error("expected empty output when OmitEmpty is true and no changes")
	}
}

func TestExport_InvalidFormat(t *testing.T) {
	var buf bytes.Buffer
	err := Export(&buf, exportChanges, ExportOptions{Format: "xml"})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestParseExportFormat_Valid(t *testing.T) {
	for _, s := range []string{"csv", "json", "tsv", "CSV", "JSON"} {
		_, err := ParseExportFormat(s)
		if err != nil {
			t.Errorf("expected %q to be valid, got error: %v", s, err)
		}
	}
}

func TestParseExportFormat_Invalid(t *testing.T) {
	_, err := ParseExportFormat("toml")
	if err == nil {
		t.Error("expected error for invalid format 'toml'")
	}
}
