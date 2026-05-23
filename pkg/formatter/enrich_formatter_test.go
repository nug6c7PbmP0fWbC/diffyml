package formatter

import (
	"bytes"
	"strings"
	"testing"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

func TestEnrichFormatter_NoChanges(t *testing.T) {
	f := NewEnrichFormatter()
	var buf bytes.Buffer
	if err := f.Format(&buf, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "no changes") {
		t.Errorf("expected 'no changes', got %q", buf.String())
	}
}

func TestEnrichFormatter_NoMetadata(t *testing.T) {
	f := NewEnrichFormatter()
	var buf bytes.Buffer
	changes := []diff.Change{
		{Path: "app.name", Type: diff.ChangeModified},
	}
	if err := f.Format(&buf, changes); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "app.name") {
		t.Errorf("expected path in output, got %q", out)
	}
	if strings.Contains(out, "[") {
		t.Errorf("expected no metadata brackets when metadata is empty")
	}
}

func TestEnrichFormatter_WithMetadata(t *testing.T) {
	f := NewEnrichFormatter()
	var buf bytes.Buffer
	changes := []diff.Change{
		{
			Path:     "db.host",
			Type:     diff.ChangeAdded,
			Metadata: map[string]string{"owner": "dba", "env": "prod"},
		},
	}
	if err := f.Format(&buf, changes); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "env=prod") {
		t.Errorf("expected env=prod in output, got %q", out)
	}
	if !strings.Contains(out, "owner=dba") {
		t.Errorf("expected owner=dba in output, got %q", out)
	}
}

func TestEnrichFormatter_MetadataSorted(t *testing.T) {
	f := NewEnrichFormatter()
	var buf bytes.Buffer
	changes := []diff.Change{
		{
			Path:     "x",
			Type:     diff.ChangeRemoved,
			Metadata: map[string]string{"z": "last", "a": "first"},
		},
	}
	if err := f.Format(&buf, changes); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	aIdx := strings.Index(out, "a=first")
	zIdx := strings.Index(out, "z=last")
	if aIdx < 0 || zIdx < 0 || aIdx > zIdx {
		t.Errorf("expected sorted metadata: a before z, got %q", out)
	}
}

func TestEnrichFormatter_SymbolsCorrect(t *testing.T) {
	f := NewEnrichFormatter()
	var buf bytes.Buffer
	changes := []diff.Change{
		{Path: "a", Type: diff.ChangeAdded},
		{Path: "b", Type: diff.ChangeRemoved},
		{Path: "c", Type: diff.ChangeModified},
	}
	if err := f.Format(&buf, changes); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	expected := []string{"+", "-", "~"}
	for i, line := range lines {
		if !strings.HasPrefix(line, expected[i]) {
			t.Errorf("line %d: expected prefix %q, got %q", i, expected[i], line)
		}
	}
}
