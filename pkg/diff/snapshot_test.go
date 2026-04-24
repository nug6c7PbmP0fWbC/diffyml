package diff

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveAndLoadSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	data := map[string]interface{}{
		"version": "1.0",
		"debug":   true,
	}

	if err := SaveSnapshot(path, "test-label", data); err != nil {
		t.Fatalf("SaveSnapshot error: %v", err)
	}

	snap, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot error: %v", err)
	}

	if snap.Label != "test-label" {
		t.Errorf("expected label 'test-label', got %q", snap.Label)
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	if snap.Data["version"] != "1.0" {
		t.Errorf("expected version '1.0', got %v", snap.Data["version"])
	}
}

func TestLoadSnapshot_NotFound(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestDiffSnapshot_DetectsChanges(t *testing.T) {
	snap := &Snapshot{
		Timestamp: time.Now(),
		Label:     "baseline",
		Data: map[string]interface{}{
			"key": "old",
		},
	}

	current := map[string]interface{}{
		"key": "new",
	}

	changes := DiffSnapshot(snap, current)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != Modified {
		t.Errorf("expected Modified, got %v", changes[0].Type)
	}
}

func TestDiffSnapshot_NoChanges(t *testing.T) {
	snap := &Snapshot{
		Timestamp: time.Now(),
		Label:     "stable",
		Data: map[string]interface{}{
			"key": "value",
		},
	}

	current := map[string]interface{}{
		"key": "value",
	}

	changes := DiffSnapshot(snap, current)
	if len(changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(changes))
	}
}

func TestSaveSnapshot_InvalidPath(t *testing.T) {
	err := SaveSnapshot("/nonexistent/dir/snap.json", "x", map[string]interface{}{})
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestLoadSnapshot_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o644)

	_, err := LoadSnapshot(path)
	if err == nil {
		t.Fatal("expected unmarshal error")
	}
}
