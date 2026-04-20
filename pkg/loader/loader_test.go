package loader

import (
	"os"
	"testing"
)

func TestParseBytes_Valid(t *testing.T) {
	input := []byte("key: value\nnested:\n  inner: 42\n")
	m, err := ParseBytes(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["key"] != "value" {
		t.Errorf("expected key=value, got %v", m["key"])
	}
	nested, ok := m["nested"].(map[string]interface{})
	if !ok {
		t.Fatal("expected nested to be a map")
	}
	// Note: gopkg.in/yaml.v2 unmarshals integers as int, not int64
	if nested["inner"] != 42 {
		t.Errorf("expected inner=42, got %v", nested["inner"])
	}
}

func TestParseBytes_Empty(t *testing.T) {
	m, err := ParseBytes([]byte{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 0 {
		t.Error("expected empty map for empty input")
	}
}

func TestParseBytes_Invalid(t *testing.T) {
	_, err := ParseBytes([]byte("key: [unclosed"))
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := LoadFile("/nonexistent/path/file.yml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadFile_Valid(t *testing.T) {
	f, err := os.CreateTemp("", "diffyml-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString("hello: world\n")
	f.Close()

	m, err := LoadFile(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["hello"] != "world" {
		t.Errorf("expected hello=world, got %v", m["hello"])
	}
}

// TestLoadFile_EmptyFile verifies that loading an empty YAML file returns an empty map without error.
func TestLoadFile_EmptyFile(t *testing.T) {
	f, err := os.CreateTemp("", "diffyml-empty-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Close()

	m, err := LoadFile(f.Name())
	if err != nil {
		t.Fatalf("unexpected error for empty file: %v", err)
	}
	if len(m) != 0 {
		t.Errorf("expected empty map for empty file, got %v", m)
	}
}
