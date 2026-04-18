package loader

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadFile reads a YAML file from disk and returns a parsed map.
func LoadFile(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("loader: failed to read file %q: %w", path, err)
	}
	return ParseBytes(data)
}

// ParseBytes parses raw YAML bytes into a map.
func ParseBytes(data []byte) (map[string]interface{}, error) {
	var out map[string]interface{}
	if err := yaml.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("loader: failed to parse YAML: %w", err)
	}
	if out == nil {
		out = make(map[string]interface{})
	}
	return out, nil
}
