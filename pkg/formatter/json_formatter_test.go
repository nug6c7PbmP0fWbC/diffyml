package formatter

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

func TestJSONFormatter_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	f := NewJSONFormatter(&buf)
	if err := f.Format([]diff.Change{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["total"].(float64) != 0 {
		t.Errorf("expected total 0, got %v", out["total"])
	}
}

func TestJSONFormatter_Added(t *testing.T) {
	var buf bytes.Buffer
	f := NewJSONFormatter(&buf)
	changes := []diff.Change{
		{Type: diff.Added, Path: "server.port", NewValue: 8080},
	}
	if err := f.Format(changes); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	json.Unmarshal(buf.Bytes(), &out)
	if out["total"].(float64) != 1 {
		t.Errorf("expected total 1")
	}
	changesList := out["changes"].([]interface{})
	c := changesList[0].(map[string]interface{})
	if c["type"] != "added" {
		t.Errorf("expected type added, got %v", c["type"])
	}
	if c["path"] != "server.port" {
		t.Errorf("expected path server.port, got %v", c["path"])
	}
}

func TestJSONFormatter_Modified(t *testing.T) {
	var buf bytes.Buffer
	f := NewJSONFormatter(&buf)
	changes := []diff.Change{
		{Type: diff.Modified, Path: "app.name", OldValue: "foo", NewValue: "bar"},
	}
	f.Format(changes)
	var out map[string]interface{}
	json.Unmarshal(buf.Bytes(), &out)
	changesList := out["changes"].([]interface{})
	c := changesList[0].(map[string]interface{})
	if c["old_value"] != "foo" || c["new_value"] != "bar" {
		t.Errorf("unexpected values: %v", c)
	}
}

func TestJSONFormatter_Removed(t *testing.T) {
	var buf bytes.Buffer
	f := NewJSONFormatter(&buf)
	changes := []diff.Change{
		{Type: diff.Removed, Path: "db.host", OldValue: "localhost"},
	}
	f.Format(changes)
	var out map[string]interface{}
	json.Unmarshal(buf.Bytes(), &out)
	changesList := out["changes"].([]interface{})
	c := changesList[0].(map[string]interface{})
	if c["type"] != "removed" {
		t.Errorf("expected removed, got %v", c["type"])
	}
	if _, ok := c["new_value"]; ok {
		t.Errorf("new_value should be omitted for removed changes")
	}
}
