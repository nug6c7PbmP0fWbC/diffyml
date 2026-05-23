package diff

import (
	"testing"
)

func TestEnrich_NoRules(t *testing.T) {
	changes := []Change{{Path: "a", Type: ChangeAdded}}
	out := Enrich(changes, nil)
	if len(out) != 1 || out[0].Metadata != nil {
		t.Fatal("expected unchanged change with no rules")
	}
}

func TestEnrich_MatchAll(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeRemoved},
	}
	rules := []EnrichRule{
		{Meta: map[string]string{"env": "prod"}},
	}
	out := Enrich(changes, rules)
	for _, c := range out {
		if c.Metadata["env"] != "prod" {
			t.Errorf("path %s: expected env=prod, got %v", c.Path, c.Metadata)
		}
	}
}

func TestEnrich_FilterByType(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded},
		{Path: "b", Type: ChangeRemoved},
	}
	rules := []EnrichRule{
		{Type: ChangeAdded, Meta: map[string]string{"tag": "new"}},
	}
	out := Enrich(changes, rules)
	if out[0].Metadata["tag"] != "new" {
		t.Errorf("expected tag=new on added change")
	}
	if out[1].Metadata != nil && out[1].Metadata["tag"] != "" {
		t.Errorf("expected no tag on removed change")
	}
}

func TestEnrich_FilterByPathPrefix(t *testing.T) {
	changes := []Change{
		{Path: "db.host", Type: ChangeModified},
		{Path: "app.name", Type: ChangeModified},
	}
	rules := []EnrichRule{
		{PathPrefix: "db", Meta: map[string]string{"owner": "dba"}},
	}
	out := Enrich(changes, rules)
	if out[0].Metadata["owner"] != "dba" {
		t.Errorf("expected owner=dba on db.host")
	}
	if out[1].Metadata != nil && out[1].Metadata["owner"] != "" {
		t.Errorf("expected no owner on app.name")
	}
}

func TestEnrich_MultipleRules_LastWins(t *testing.T) {
	changes := []Change{{Path: "x", Type: ChangeAdded}}
	rules := []EnrichRule{
		{Meta: map[string]string{"priority": "low"}},
		{Meta: map[string]string{"priority": "high"}},
	}
	out := Enrich(changes, rules)
	if out[0].Metadata["priority"] != "high" {
		t.Errorf("expected last rule to win, got %v", out[0].Metadata["priority"])
	}
}

func TestEnrich_PreservesExistingMetadata(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeAdded, Metadata: map[string]string{"existing": "yes"}},
	}
	rules := []EnrichRule{
		{Meta: map[string]string{"new": "val"}},
	}
	out := Enrich(changes, rules)
	if out[0].Metadata["existing"] != "yes" {
		t.Errorf("expected existing metadata to be preserved")
	}
	if out[0].Metadata["new"] != "val" {
		t.Errorf("expected new metadata to be added")
	}
}
