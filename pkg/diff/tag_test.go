package diff

import (
	"testing"
)

func TestTagChanges_NoRules(t *testing.T) {
	changes := []Change{
		{Path: "a.b", Type: Added, NewValue: 1},
	}
	result := TagChanges(changes, nil)
	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}
	if len(result[0].Tags) != 0 {
		t.Errorf("expected no tags, got %v", result[0].Tags)
	}
}

func TestTagChanges_MatchAll(t *testing.T) {
	changes := []Change{
		{Path: "x", Type: Added},
		{Path: "y", Type: Removed},
	}
	rules := []TagRule{
		{Tag: Tag{Key: "env", Value: "prod"}},
	}
	result := TagChanges(changes, rules)
	for _, tc := range result {
		if len(tc.Tags) != 1 || tc.Tags[0].Key != "env" {
			t.Errorf("expected env tag on %s, got %v", tc.Path, tc.Tags)
		}
	}
}

func TestTagChanges_FilterByType(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: Added},
		{Path: "b", Type: Removed},
		{Path: "c", Type: Modified},
	}
	rules := []TagRule{
		{Types: []ChangeType{Added}, Tag: Tag{Key: "new", Value: "true"}},
	}
	result := TagChanges(changes, rules)
	tagged := 0
	for _, tc := range result {
		if len(tc.Tags) > 0 {
			tagged++
			if tc.Type != Added {
				t.Errorf("unexpected tag on type %v", tc.Type)
			}
		}
	}
	if tagged != 1 {
		t.Errorf("expected 1 tagged change, got %d", tagged)
	}
}

func TestTagChanges_FilterByPathPrefix(t *testing.T) {
	changes := []Change{
		{Path: "db.host", Type: Modified},
		{Path: "db.port", Type: Modified},
		{Path: "app.name", Type: Modified},
	}
	rules := []TagRule{
		{PathPrefix: "db", Tag: Tag{Key: "component", Value: "database"}},
	}
	result := TagChanges(changes, rules)
	for _, tc := range result {
		if tc.Path == "db.host" || tc.Path == "db.port" {
			if len(tc.Tags) != 1 || tc.Tags[0].Value != "database" {
				t.Errorf("expected database tag on %s", tc.Path)
			}
		} else {
			if len(tc.Tags) != 0 {
				t.Errorf("unexpected tag on %s", tc.Path)
			}
		}
	}
}

func TestTagChanges_MultipleRules(t *testing.T) {
	changes := []Change{
		{Path: "db.password", Type: Modified},
	}
	rules := []TagRule{
		{PathPrefix: "db", Tag: Tag{Key: "component", Value: "database"}},
		{PathPrefix: "db.password", Tag: Tag{Key: "sensitive", Value: "true"}},
	}
	result := TagChanges(changes, rules)
	if len(result[0].Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(result[0].Tags))
	}
}
