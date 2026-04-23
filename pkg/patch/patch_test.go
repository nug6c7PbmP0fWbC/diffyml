package patch

import (
	"testing"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

func baseData() map[string]interface{} {
	return map[string]interface{}{
		"name":    "alice",
		"version": "1.0",
		"db": map[string]interface{}{
			"host": "localhost",
			"port": 5432,
		},
	}
}

func TestApply_Forward_Modified(t *testing.T) {
	data := baseData()
	changes := []diff.Change{
		{Type: diff.Modified, Path: "version", OldValue: "1.0", NewValue: "2.0"},
	}
	if err := Apply(data, changes, Forward); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["version"] != "2.0" {
		t.Errorf("expected version=2.0, got %v", data["version"])
	}
}

func TestApply_Reverse_Modified(t *testing.T) {
	data := baseData()
	changes := []diff.Change{
		{Type: diff.Modified, Path: "version", OldValue: "1.0", NewValue: "2.0"},
	}
	if err := Apply(data, changes, Reverse); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["version"] != "1.0" {
		t.Errorf("expected version=1.0, got %v", data["version"])
	}
}

func TestApply_Forward_Added(t *testing.T) {
	data := baseData()
	changes := []diff.Change{
		{Type: diff.Added, Path: "region", NewValue: "us-east-1"},
	}
	if err := Apply(data, changes, Forward); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %v", data["region"])
	}
}

func TestApply_Reverse_Added(t *testing.T) {
	data := baseData()
	changes := []diff.Change{
		{Type: diff.Added, Path: "region", NewValue: "us-east-1"},
	}
	if err := Apply(data, changes, Reverse); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, exists := data["region"]; exists {
		t.Error("expected region to be removed")
	}
}

func TestApply_Forward_Removed(t *testing.T) {
	data := baseData()
	changes := []diff.Change{
		{Type: diff.Removed, Path: "name", OldValue: "alice"},
	}
	if err := Apply(data, changes, Forward); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, exists := data["name"]; exists {
		t.Error("expected name to be removed")
	}
}

func TestApply_Nested_Modified(t *testing.T) {
	data := baseData()
	changes := []diff.Change{
		{Type: diff.Modified, Path: "db.host", OldValue: "localhost", NewValue: "db.prod.internal"},
	}
	if err := Apply(data, changes, Forward); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	db := data["db"].(map[string]interface{})
	if db["host"] != "db.prod.internal" {
		t.Errorf("expected db.host=db.prod.internal, got %v", db["host"])
	}
}

func TestApply_EmptyChanges(t *testing.T) {
	data := baseData()
	origName := data["name"]
	if err := Apply(data, nil, Forward); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["name"] != origName {
		t.Error("data should not change with empty changes")
	}
}
