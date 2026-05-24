package diff

import (
	"testing"
)

func TestZoom_Empty(t *testing.T) {
	result := Zoom(nil, ZoomOptions{Prefix: "app"})
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestZoom_NoPrefix(t *testing.T) {
	changes := []Change{
		{Path: "a.b", Type: ChangeAdded},
	}
	result := Zoom(changes, ZoomOptions{Prefix: ""})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestZoom_FiltersOutNonMatching(t *testing.T) {
	changes := []Change{
		{Path: "app.db.host", Type: ChangeModified},
		{Path: "app.name", Type: ChangeModified},
		{Path: "app.db.port", Type: ChangeAdded},
	}
	result := Zoom(changes, ZoomOptions{Prefix: "app.db"})
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	if result[0].Path != "app.db.host" || result[1].Path != "app.db.port" {
		t.Fatalf("unexpected paths: %v", result)
	}
}

func TestZoom_StripPrefix(t *testing.T) {
	changes := []Change{
		{Path: "app.db.host", Type: ChangeModified},
		{Path: "app.db.port", Type: ChangeAdded},
	}
	result := Zoom(changes, ZoomOptions{Prefix: "app.db", StripPrefix: true})
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
	if result[0].Path != "host" {
		t.Errorf("expected 'host', got %q", result[0].Path)
	}
	if result[1].Path != "port" {
		t.Errorf("expected 'port', got %q", result[1].Path)
	}
}

func TestZoom_ExactPrefixMatch(t *testing.T) {
	changes := []Change{
		{Path: "app.db", Type: ChangeModified},
	}
	result := Zoom(changes, ZoomOptions{Prefix: "app.db", StripPrefix: true})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Path != "" {
		t.Errorf("expected empty path, got %q", result[0].Path)
	}
}

func TestZoom_NoPrefixFalsePositive(t *testing.T) {
	// "app.database" should NOT match prefix "app.db"
	changes := []Change{
		{Path: "app.database.host", Type: ChangeModified},
		{Path: "app.db.host", Type: ChangeAdded},
	}
	result := Zoom(changes, ZoomOptions{Prefix: "app.db"})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Path != "app.db.host" {
		t.Errorf("unexpected path: %q", result[0].Path)
	}
}

func TestZoom_PreservesChangeType(t *testing.T) {
	changes := []Change{
		{Path: "svc.port", Type: ChangeRemoved, Before: 8080},
	}
	result := Zoom(changes, ZoomOptions{Prefix: "svc", StripPrefix: true})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Type != ChangeRemoved {
		t.Errorf("expected ChangeRemoved, got %v", result[0].Type)
	}
	if result[0].Before != 8080 {
		t.Errorf("expected Before=8080, got %v", result[0].Before)
	}
}
