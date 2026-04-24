package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge_AddKey(t *testing.T) {
	base := map[string]interface{}{
		"name": "alice",
	}
	override := map[string]interface{}{
		"age": 30,
	}
	result, err := Merge(base, override)
	assert.NoError(t, err)
	assert.Equal(t, "alice", result["name"])
	assert.Equal(t, 30, result["age"])
}

func TestMerge_OverrideKey(t *testing.T) {
	base := map[string]interface{}{
		"name": "alice",
		"role": "admin",
	}
	override := map[string]interface{}{
		"role": "viewer",
	}
	result, err := Merge(base, override)
	assert.NoError(t, err)
	assert.Equal(t, "viewer", result["role"])
	assert.Equal(t, "alice", result["name"])
}

func TestMerge_NestedKeys(t *testing.T) {
	base := map[string]interface{}{
		"database": map[string]interface{}{
			"host": "localhost",
			"port": 5432,
		},
	}
	override := map[string]interface{}{
		"database": map[string]interface{}{
			"port": 5433,
		},
	}
	result, err := Merge(base, override)
	assert.NoError(t, err)
	db, ok := result["database"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "localhost", db["host"])
	assert.Equal(t, 5433, db["port"])
}

func TestMerge_NilBase(t *testing.T) {
	override := map[string]interface{}{"key": "val"}
	result, err := Merge(nil, override)
	assert.NoError(t, err)
	assert.Equal(t, "val", result["key"])
}

func TestMerge_NilOverride(t *testing.T) {
	base := map[string]interface{}{"key": "val"}
	result, err := Merge(base, nil)
	assert.NoError(t, err)
	assert.Equal(t, "val", result["key"])
}

func TestMerge_DoesNotMutateBase(t *testing.T) {
	base := map[string]interface{}{"key": "original"}
	override := map[string]interface{}{"key": "changed"}
	_, err := Merge(base, override)
	assert.NoError(t, err)
	assert.Equal(t, "original", base["key"])
}
