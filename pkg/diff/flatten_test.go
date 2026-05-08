package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlatten_Empty(t *testing.T) {
	result := Flatten(map[string]interface{}{}, DefaultFlattenOptions())
	assert.Empty(t, result)
}

func TestFlatten_Flat(t *testing.T) {
	input := map[string]interface{}{
		"key": "value",
		"num": 42,
	}
	result := Flatten(input, DefaultFlattenOptions())
	assert.Equal(t, "value", result["key"])
	assert.Equal(t, 42, result["num"])
}

func TestFlatten_Nested(t *testing.T) {
	input := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "deep",
			},
		},
	}
	result := Flatten(input, DefaultFlattenOptions())
	assert.Equal(t, "deep", result["a.b.c"])
	assert.Len(t, result, 1)
}

func TestFlatten_CustomSeparator(t *testing.T) {
	input := map[string]interface{}{
		"x": map[string]interface{}{
			"y": 99,
		},
	}
	opts := FlattenOptions{Separator: "/"}
	result := Flatten(input, opts)
	assert.Equal(t, 99, result["x/y"])
}

func TestFlatten_MaxDepth(t *testing.T) {
	input := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "leaf",
			},
		},
	}
	opts := FlattenOptions{Separator: ".", MaxDepth: 2}
	result := Flatten(input, opts)
	// At depth 2 we stop recursing, so "a.b" should hold the nested map
	assert.Contains(t, result, "a.b")
	_, isMap := result["a.b"].(map[string]interface{})
	assert.True(t, isMap)
	assert.NotContains(t, result, "a.b.c")
}

func TestFlatten_NilLeaf(t *testing.T) {
	input := map[string]interface{}{
		"key": nil,
	}
	result := Flatten(input, DefaultFlattenOptions())
	val, ok := result["key"]
	assert.True(t, ok)
	assert.Nil(t, val)
}

func TestFlatten_EmptyNestedMap(t *testing.T) {
	input := map[string]interface{}{
		"empty": map[string]interface{}{},
	}
	result := Flatten(input, DefaultFlattenOptions())
	// Empty nested map should be kept as-is (leaf)
	assert.Contains(t, result, "empty")
}
