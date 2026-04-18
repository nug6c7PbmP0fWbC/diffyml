package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/diffyml/diffyml/pkg/filter"
	"github.com/diffyml/diffyml/pkg/reporter"
)

func TestReporter_TextFormat(t *testing.T) {
	r, err := reporter.New(reporter.Config{Format: "text"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a := map[string]interface{}{"key": "old"}
	b := map[string]interface{}{"key": "new"}

	var buf bytes.Buffer
	if err := r.Report(&buf, a, b); err != nil {
		t.Fatalf("Report failed: %v", err)
	}

	if !strings.Contains(buf.String(), "key") {
		t.Errorf("expected output to contain 'key', got: %s", buf.String())
	}
}

func TestReporter_JSONFormat(t *testing.T) {
	r, err := reporter.New(reporter.Config{Format: "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a := map[string]interface{}{"x": 1}
	b := map[string]interface{}{"x": 2}

	var buf bytes.Buffer
	if err := r.Report(&buf, a, b); err != nil {
		t.Fatalf("Report failed: %v", err)
	}

	if !strings.Contains(buf.String(), "\"path\"") {
		t.Errorf("expected JSON output, got: %s", buf.String())
	}
}

func TestReporter_WithFilter(t *testing.T) {
	r, err := reporter.New(reporter.Config{
		Format:    "text",
		FilterCfg: filter.Config{Types: []string{"added"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a := map[string]interface{}{}
	b := map[string]interface{}{"newkey": "val", "other": "x"}

	var buf bytes.Buffer
	if err := r.Report(&buf, a, b); err != nil {
		t.Fatalf("Report failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("expected non-empty output for added changes")
	}
}

func TestReporter_InvalidFormat(t *testing.T) {
	_, err := reporter.New(reporter.Config{Format: "xml"})
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
