package formatter_test

import (
	"bytes"
	"testing"

	"github.com/diffyml/diffyml/pkg/diff"
	"github.com/diffyml/diffyml/pkg/formatter"
)

func TestTextFormatter_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.NewTextFormatter(&buf)
	if err := f.Write(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := buf.String()
	if got != "No differences found.\n" {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestTextFormatter_Added(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.NewTextFormatter(&buf)
	changes := []diff.Change{
		{Type: diff.Added, Path: []string{"server", "port"}, To: 8080},
	}
	if err := f.Write(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "+ server.port: 8080\n"
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}

func TestTextFormatter_Removed(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.NewTextFormatter(&buf)
	changes := []diff.Change{
		{Type: diff.Removed, Path: []string{"db"}, From: "postgres"},
	}
	if err := f.Write(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "- db: postgres\n"
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}

func TestTextFormatter_Modified(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.NewTextFormatter(&buf)
	changes := []diff.Change{
		{Type: diff.Modified, Path: []string{"timeout"}, From: 30, To: 60},
	}
	if err := f.Write(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "~ timeout: 30 -> 60\n"
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}

// TestTextFormatter_MultipleChanges verifies that multiple changes are all written correctly.
func TestTextFormatter_MultipleChanges(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.NewTextFormatter(&buf)
	changes := []diff.Change{
		{Type: diff.Added, Path: []string{"host"}, To: "localhost"},
		{Type: diff.Removed, Path: []string{"debug"}, From: true},
		{Type: diff.Modified, Path: []string{"workers"}, From: 2, To: 4},
	}
	if err := f.Write(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "+ host: localhost\n- debug: true\n~ workers: 2 -> 4\n"
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}
