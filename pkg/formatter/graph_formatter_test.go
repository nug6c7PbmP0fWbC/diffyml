package formatter

import (
	"strings"
	"testing"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

func TestGraphFormatter_NoChanges(t *testing.T) {
	f := NewGraphFormatter()
	out := f.Format(nil)
	if out != "(no changes)\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestGraphFormatter_FlatChanges(t *testing.T) {
	f := NewGraphFormatter()
	changes := []diff.Change{
		{Path: "alpha", Type: diff.ChangeAdded},
		{Path: "beta", Type: diff.ChangeRemoved},
	}
	out := f.Format(changes)
	if !strings.Contains(out, "- alpha") {
		t.Errorf("expected alpha in output, got: %s", out)
	}
	if !strings.Contains(out, "- beta") {
		t.Errorf("expected beta in output, got: %s", out)
	}
}

func TestGraphFormatter_NestedChanges(t *testing.T) {
	f := NewGraphFormatter()
	changes := []diff.Change{
		{Path: "server", Type: diff.ChangeModified},
		{Path: "server.port", Type: diff.ChangeModified},
		{Path: "server.host", Type: diff.ChangeAdded},
	}
	out := f.Format(changes)
	if !strings.Contains(out, "Change Graph:") {
		t.Error("expected header in output")
	}
	serverIdx := strings.Index(out, "- server\n")
	portIdx := strings.Index(out, "- server.port")
	if serverIdx == -1 {
		t.Error("expected server node in output")
	}
	if portIdx == -1 {
		t.Error("expected server.port node in output")
	}
	if serverIdx > portIdx {
		t.Error("expected server to appear before server.port")
	}
}

func TestGraphFormatter_HeaderPresent(t *testing.T) {
	f := NewGraphFormatter()
	changes := []diff.Change{
		{Path: "x", Type: diff.ChangeAdded},
	}
	out := f.Format(changes)
	if !strings.HasPrefix(out, "Change Graph:") {
		t.Errorf("expected output to start with 'Change Graph:', got: %s", out)
	}
}
