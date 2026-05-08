package diff

import (
	"testing"
)

func TestEnrich_NoRules(t *testing.T) {
	changes := []Change{
		{Path: "a.b", Type: ChangeTypeAdded, Value: "v"},
	}
	result := Enrich(changes, nil)
	if len(result) != 1 {
		t.Fatalf("expected 1 enriched change, got %d", len(result))
	}
	if len(result[0].Meta) != 0 {
		t.Errorf("expected empty meta, got %v", result[0].Meta)
	}
}

func TestEnrich_MatchAll(t *testing.T) {
	changes := []Change{
		{Path: "x", Type: ChangeTypeAdded, Value: "1"},
		{Path: "y", Type: ChangeTypeRemoved, Value: "2"},
	}
	rules := []EnrichRule{
		{Meta: map[string]string{"owner": "team-a"}},
	}
	result := Enrich(changes, rules)
	for _, ec := range result {
		if ec.Meta["owner"] != "team-a" {
			t.Errorf("expected owner=team-a, got %q", ec.Meta["owner"])
		}
	}
}

func TestEnrich_FilterByType(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeTypeAdded, Value: "1"},
		{Path: "b", Type: ChangeTypeRemoved, Value: "2"},
	}
	rules := []EnrichRule{
		{Type: "added", Meta: map[string]string{"flag": "new"}},
	}
	result := Enrich(changes, rules)
	if result[0].Meta["flag"] != "new" {
		t.Errorf("expected flag=new for added change")
	}
	if _, ok := result[1].Meta["flag"]; ok {
		t.Errorf("removed change should not have flag meta")
	}
}

func TestEnrich_FilterByPathPrefix(t *testing.T) {
	changes := []Change{
		{Path: "db.host", Type: ChangeTypeModified, Value: "new"},
		{Path: "app.port", Type: ChangeTypeModified, Value: "8080"},
	}
	rules := []EnrichRule{
		{PathPrefix: "db", Meta: map[string]string{"sensitive": "true"}},
	}
	result := Enrich(changes, rules)
	if result[0].Meta["sensitive"] != "true" {
		t.Errorf("expected sensitive=true for db.host")
	}
	if _, ok := result[1].Meta["sensitive"]; ok {
		t.Errorf("app.port should not be sensitive")
	}
}

func TestEnrich_MultipleRules_LastWins(t *testing.T) {
	changes := []Change{
		{Path: "cfg.key", Type: ChangeTypeAdded, Value: "v"},
	}
	rules := []EnrichRule{
		{Meta: map[string]string{"tier": "low"}},
		{PathPrefix: "cfg", Meta: map[string]string{"tier": "high"}},
	}
	result := Enrich(changes, rules)
	if result[0].Meta["tier"] != "high" {
		t.Errorf("expected tier=high (last rule wins), got %q", result[0].Meta["tier"])
	}
}

func TestEnrich_EmptyChanges(t *testing.T) {
	result := Enrich(nil, []EnrichRule{{Meta: map[string]string{"k": "v"}}})
	if len(result) != 0 {
		t.Errorf("expected empty result for nil changes")
	}
}
