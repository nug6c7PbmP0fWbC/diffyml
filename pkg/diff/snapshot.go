package diff

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a saved state of a YAML document at a point in time.
type Snapshot struct {
	Timestamp time.Time              `json:"timestamp"`
	Label     string                 `json:"label"`
	Data      map[string]interface{} `json:"data"`
}

// SaveSnapshot writes a snapshot of the given data to a JSON file.
func SaveSnapshot(path, label string, data map[string]interface{}) error {
	s := Snapshot{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Data:      data,
	}
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal error: %w", err)
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return fmt.Errorf("snapshot: write error: %w", err)
	}
	return nil
}

// LoadSnapshot reads a snapshot from a JSON file.
func LoadSnapshot(path string) (*Snapshot, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read error: %w", err)
	}
	var s Snapshot
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal error: %w", err)
	}
	return &s, nil
}

// DiffSnapshot compares a snapshot's data against a current data map and
// returns the list of changes.
func DiffSnapshot(snap *Snapshot, current map[string]interface{}) []Change {
	return Compare(snap.Data, current)
}
