package formatter

import (
	"strings"
	"testing"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

var efChanges = []diff.Change{
	{Path: "db.host", Type: diff.ChangeModified, OldValue: "localhost", NewValue: "prod-db"},
	{Path: "db.port", Type: diff.ChangeAdded, OldValue: nil, NewValue: 5432},
}

func TestExportFormatter_CSV(t *testing.T) {
	f, err := NewExportFormatter("csv")
	if err != nil {
		t.Fatalf("NewExportFormatter: %v", err)
	}
	out, err := f.FormatString(efChanges)
	if err != nil {
		t.Fatalf("FormatString: %v", err)
	}
	if !strings.Contains(out, "db.host") {
		t.Error("expected db.host in CSV output")
	}
	if !strings.Contains(out, "path,type,old_value,new_value") {
		t.Error("expected CSV header row")
	}
}

func TestExportFormatter_TSV(t *testing.T) {
	f, err := NewExportFormatter("tsv")
	if err != nil {
		t.Fatalf("NewExportFormatter: %v", err)
	}
	out, err := f.FormatString(efChanges)
	if err != nil {
		t.Fatalf("FormatString: %v", err)
	}
	if !strings.Contains(out, "\t") {
		t.Error("expected tab separator in TSV output")
	}
}

func TestExportFormatter_JSON(t *testing.T) {
	f, err := NewExportFormatter("json")
	if err != nil {
		t.Fatalf("NewExportFormatter: %v", err)
	}
	out, err := f.FormatString(efChanges)
	if err != nil {
		t.Fatalf("FormatString: %v", err)
	}
	if !strings.HasPrefix(strings.TrimSpace(out), "[") {
		t.Error("expected JSON array output")
	}
	if !strings.Contains(out, `"db.port"`) {
		t.Error("expected db.port in JSON output")
	}
}

func TestExportFormatter_InvalidFormat(t *testing.T) {
	_, err := NewExportFormatter("xml")
	if err == nil {
		t.Error("expected error for unsupported format 'xml'")
	}
}

func TestExportFormatter_EmptyChanges(t *testing.T) {
	f, err := NewExportFormatter("csv")
	if err != nil {
		t.Fatalf("NewExportFormatter: %v", err)
	}
	out, err := f.FormatString([]diff.Change{})
	if err != nil {
		t.Fatalf("FormatString: %v", err)
	}
	// header only
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 line (header only), got %d", len(lines))
	}
}
