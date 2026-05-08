package diff

import (
	"testing"
)

func makePartitionChanges() []Change {
	return []Change{
		{Path: "a", Type: Added},
		{Path: "b", Type: Added},
		{Path: "c", Type: Removed},
		{Path: "d", Type: Modified},
		{Path: "e", Type: Modified},
		{Path: "f", Type: Added},
	}
}

func TestPartition_Empty(t *testing.T) {
	result := Partition(nil, DefaultPartitionOptions())
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestPartition_NoLimit(t *testing.T) {
	opts := PartitionOptions{MaxPerPartition: 0, GroupByType: false}
	changes := makePartitionChanges()
	result := Partition(changes, opts)
	if len(result) != 1 {
		t.Fatalf("expected 1 partition, got %d", len(result))
	}
	if len(result[0]) != len(changes) {
		t.Errorf("expected %d changes, got %d", len(changes), len(result[0]))
	}
}

func TestPartition_MaxCaps(t *testing.T) {
	opts := PartitionOptions{MaxPerPartition: 2, GroupByType: false}
	changes := makePartitionChanges() // 6 changes
	result := Partition(changes, opts)
	if len(result) != 3 {
		t.Fatalf("expected 3 partitions, got %d", len(result))
	}
	for i, p := range result {
		if len(p) != 2 {
			t.Errorf("partition %d: expected 2 changes, got %d", i, len(p))
		}
	}
}

func TestPartition_GroupByType(t *testing.T) {
	opts := PartitionOptions{MaxPerPartition: 0, GroupByType: true}
	changes := makePartitionChanges()
	result := Partition(changes, opts)
	// Added: a,b,f  Removed: c  Modified: d,e
	if len(result) != 3 {
		t.Fatalf("expected 3 partitions (by type), got %d", len(result))
	}
	for _, p := range result {
		ct := p[0].Type
		for _, c := range p {
			if c.Type != ct {
				t.Errorf("mixed types in partition: expected %v, got %v", ct, c.Type)
			}
		}
	}
}

func TestPartition_GroupByTypeAndMax(t *testing.T) {
	opts := PartitionOptions{MaxPerPartition: 2, GroupByType: true}
	changes := makePartitionChanges()
	// Added(3)->2 partitions, Removed(1)->1, Modified(2)->1 => total 4
	result := Partition(changes, opts)
	if len(result) != 4 {
		t.Fatalf("expected 4 partitions, got %d", len(result))
	}
}

func TestPartition_SingleChange(t *testing.T) {
	opts := DefaultPartitionOptions()
	changes := []Change{{Path: "x", Type: Added}}
	result := Partition(changes, opts)
	if len(result) != 1 || len(result[0]) != 1 {
		t.Errorf("expected 1 partition with 1 change")
	}
}
