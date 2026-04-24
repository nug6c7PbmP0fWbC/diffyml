package diff

import (
	"testing"
)

func TestApplyIgnore_NoRules(t *testing.T) {
	changes := []Change{
		{Path: "a.b", Type: ChangeAdded},
		{Path: "x.y", Type: ChangeRemoved},
	}
	cfg := IgnoreConfig{}
	result := ApplyIgnore(changes, cfg)
	if len(result) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(result))
	}
}

func TestApplyIgnore_ExactMatch(t *testing.T) {
	changes := []Change{
		{Path: "metadata.name", Type: ChangeModified},
		{Path: "spec.replicas", Type: ChangeModified},
	}
	cfg := IgnoreConfig{Rules: []IgnoreRule{{Path: "metadata.name"}}}
	result := ApplyIgnore(changes, cfg)
	if len(result) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result))
	}
	if result[0].Path != "spec.replicas" {
		t.Errorf("unexpected path: %s", result[0].Path)
	}
}

func TestApplyIgnore_PrefixMatch(t *testing.T) {
	changes := []Change{
		{Path: "metadata.annotations.kubectl", Type: ChangeAdded},
		{Path: "metadata.annotations.other", Type: ChangeRemoved},
		{Path: "metadata.labels", Type: ChangeModified},
		{Path: "spec.image", Type: ChangeModified},
	}
	cfg := IgnoreConfig{Rules: []IgnoreRule{{Path: "metadata.annotations"}}}
	result := ApplyIgnore(changes, cfg)
	if len(result) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(result))
	}
}

func TestApplyIgnore_NoPrefixFalsePositive(t *testing.T) {
	// "metadata.ann" should NOT suppress "metadata.annotations"
	changes := []Change{
		{Path: "metadata.annotations", Type: ChangeAdded},
	}
	cfg := IgnoreConfig{Rules: []IgnoreRule{{Path: "metadata.ann"}}}
	result := ApplyIgnore(changes, cfg)
	if len(result) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result))
	}
}

func TestApplyIgnore_MultipleRules(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b.c", Type: ChangeRemoved},
		{Path: "d.e.f", Type: ChangeModified},
	}
	cfg := IgnoreConfig{Rules: []IgnoreRule{{Path: "a"}, {Path: "d.e"}}}
	result := ApplyIgnore(changes, cfg)
	if len(result) != 1 {
		t.Fatalf("expected 1 change, got %d", len(result))
	}
	if result[0].Path != "b.c" {
		t.Errorf("unexpected path: %s", result[0].Path)
	}
}

func TestApplyIgnore_EmptyRulePath(t *testing.T) {
	changes := []Change{
		{Path: "a.b", Type: ChangeAdded},
	}
	cfg := IgnoreConfig{Rules: []IgnoreRule{{Path: ""}}}
	result := ApplyIgnore(changes, cfg)
	if len(result) != 1 {
		t.Fatalf("empty rule should not suppress anything, got %d changes", len(result))
	}
}
