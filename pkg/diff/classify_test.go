package diff

import (
	"testing"
)

func TestClassify_NoRules(t *testing.T) {
	changes := []Change{
		{Path: "a.b", Type: ChangeAdded},
	}
	out := Classify(changes, nil)
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if out[0].Class != ClassUnclassified {
		t.Errorf("expected unclassified, got %s", out[0].Class)
	}
}

func TestClassify_NoChanges(t *testing.T) {
	rules := []ClassifyRule{{Class: ClassBreaking}}
	out := Classify(nil, rules)
	if len(out) != 0 {
		t.Fatalf("expected 0 results, got %d", len(out))
	}
}

func TestClassify_MatchAll(t *testing.T) {
	changes := []Change{
		{Path: "x", Type: ChangeAdded},
		{Path: "y", Type: ChangeRemoved},
	}
	rules := []ClassifyRule{{Class: ClassAdditive}}
	out := Classify(changes, rules)
	for _, c := range out {
		if c.Class != ClassAdditive {
			t.Errorf("expected additive, got %s", c.Class)
		}
	}
}

func TestClassify_FilterByType(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeRemoved},
		{Path: "b", Type: ChangeAdded},
	}
	rules := []ClassifyRule{
		{Type: ChangeRemoved, Class: ClassBreaking},
	}
	out := Classify(changes, rules)
	if out[0].Class != ClassBreaking {
		t.Errorf("expected breaking for removed, got %s", out[0].Class)
	}
	if out[1].Class != ClassUnclassified {
		t.Errorf("expected unclassified for added, got %s", out[1].Class)
	}
}

func TestClassify_FilterByPathPrefix(t *testing.T) {
	changes := []Change{
		{Path: "api.endpoint", Type: ChangeModified},
		{Path: "docs.readme", Type: ChangeModified},
	}
	rules := []ClassifyRule{
		{PathPrefix: "api", Class: ClassBreaking},
		{PathPrefix: "docs", Class: ClassCosmetic},
	}
	out := Classify(changes, rules)
	if out[0].Class != ClassBreaking {
		t.Errorf("expected breaking, got %s", out[0].Class)
	}
	if out[1].Class != ClassCosmetic {
		t.Errorf("expected cosmetic, got %s", out[1].Class)
	}
}

func TestClassify_FirstRuleWins(t *testing.T) {
	changes := []Change{
		{Path: "a.b", Type: ChangeAdded},
	}
	rules := []ClassifyRule{
		{Class: ClassDeprecated},
		{Class: ClassBreaking},
	}
	out := Classify(changes, rules)
	if out[0].Class != ClassDeprecated {
		t.Errorf("expected deprecated (first rule), got %s", out[0].Class)
	}
}

func TestClassify_ExactPathMatch(t *testing.T) {
	changes := []Change{
		{Path: "foo", Type: ChangeAdded},
		{Path: "foobar", Type: ChangeAdded},
	}
	rules := []ClassifyRule{
		{PathPrefix: "foo", Class: ClassBreaking},
	}
	out := Classify(changes, rules)
	if out[0].Class != ClassBreaking {
		t.Errorf("expected breaking for exact match, got %s", out[0].Class)
	}
	if out[1].Class != ClassUnclassified {
		t.Errorf("expected unclassified for foobar, got %s", out[1].Class)
	}
}
