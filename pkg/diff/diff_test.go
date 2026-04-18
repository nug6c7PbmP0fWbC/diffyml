package diff

import (
	"testing"
)

func TestCompare_Added(t *testing.T) {
	oldMap := map[string]interface{}{"a": "1"}
	newMap := map[string]interface{}{"a": "1", "b": "2"}
	result := Compare(oldMap, newMap)
	if !result.HasChanges() {
		t.Fatal("expected changes")
	}
	found := false
	for _, c := range result.Changes {
		if c.Path == "b" && c.Type == Added {
			found = true
		}
	}
	if !found {
		t.Error("expected 'b' to be added")
	}
}

func TestCompare_Removed(t *testing.T) {
	oldMap := map[string]interface{}{"a": "1", "b": "2"}
	newMap := map[string]interface{}{"a": "1"}
	result := Compare(oldMap, newMap)
	found := false
	for _, c := range result.Changes {
		if c.Path == "b" && c.Type == Removed {
			found = true
		}
	}
	if !found {
		t.Error("expected 'b' to be removed")
	}
}

func TestCompare_Modified(t *testing.T) {
	oldMap := map[string]interface{}{"a": "1"}
	newMap := map[string]interface{}{"a": "2"}
	result := Compare(oldMap, newMap)
	found := false
	for _, c := range result.Changes {
		if c.Path == "a" && c.Type == Modified && c.OldVal == "1" && c.NewVal == "2" {
			found = true
		}
	}
	if !found {
		t.Error("expected 'a' to be modified")
	}
}

func TestCompare_Nested(t *testing.T) {
	oldMap := map[string]interface{}{
		"service": map[string]interface{}{"port": "8080"},
	}
	newMap := map[string]interface{}{
		"service": map[string]interface{}{"port": "9090"},
	}
	result := Compare(oldMap, newMap)
	found := false
	for _, c := range result.Changes {
		if c.Path == "service.port" && c.Type == Modified {
			found = true
		}
	}
	if !found {
		t.Error("expected 'service.port' to be modified")
	}
}

func TestCompare_NoChanges(t *testing.T) {
	oldMap := map[string]interface{}{"x": "y"}
	newMap := map[string]interface{}{"x": "y"}
	result := Compare(oldMap, newMap)
	if result.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestSummary(t *testing.T) {
	result := &Result{
		Changes: []Change{
			{Path: "foo", Type: Added, NewVal: "bar"},
			{Path: "baz", Type: Removed, OldVal: "qux"},
		},
	}
	s := result.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
