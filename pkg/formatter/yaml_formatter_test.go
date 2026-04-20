package formatter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/diffyml/diffyml/pkg/diff"
	"github.com/diffyml/diffyml/pkg/formatter"
)

func TestYAMLFormatter_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.NewYAMLFormatter(&buf)
	if err := f.Format(nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "changes: []") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestYAMLFormatter_Added(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.NewYAMLFormatter(&buf)
	changes := []diff.Change{{Path: "app.name", Type: diff.Added, NewValue: "myapp"}}
	if err := f.Format(changes); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "type: added") {
		t.Errorf("missing type: %s", out)
	}
	if !strings.Contains(out, "new_value: myapp") {
		t.Errorf("missing new_value: %s", out)
	}
}

func TestYAMLFormatter_Removed(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.NewYAMLFormatter(&buf)
	changes := []diff.Change{{Path: "app.debug", Type: diff.Removed, OldValue: "true"}}
	if err := f.Format(changes); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "old_value: true") {
		t.Errorf("missing old_value: %s", out)
	}
	// Also verify that removed changes don't include a new_value field
	if strings.Contains(out, "new_value:") {
		t.Errorf("removed change should not have new_value: %s", out)
	}
}

func TestYAMLFormatter_Modified(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.NewYAMLFormatter(&buf)
	changes := []diff.Change{{Path: "server.port", Type: diff.Modified, OldValue: "80", NewValue: "443"}}
	if err := f.Format(changes); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "old_value: 80") || !strings.Contains(out, "new_value: 443") {
		t.Errorf("unexpected output: %s", out)
	}
}
