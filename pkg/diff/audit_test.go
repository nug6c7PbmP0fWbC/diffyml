package diff

import (
	"testing"
	"time"
)

var auditTime = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func TestAudit_Empty(t *testing.T) {
	log := Audit(nil, auditTime, DefaultAuditOptions())
	if len(log.Entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(log.Entries))
	}
}

func TestAudit_AllOperations(t *testing.T) {
	changes := []Change{
		{Type: ChangeAdded, Path: "a", NewValue: 1},
		{Type: ChangeRemoved, Path: "b", OldValue: 2},
		{Type: ChangeModified, Path: "c", OldValue: 3, NewValue: 4},
	}
	log := Audit(changes, auditTime, DefaultAuditOptions())
	if len(log.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(log.Entries))
	}
	if log.Entries[0].Operation != "added" {
		t.Errorf("expected added, got %s", log.Entries[0].Operation)
	}
	if log.Entries[1].Operation != "removed" {
		t.Errorf("expected removed, got %s", log.Entries[1].Operation)
	}
	if log.Entries[2].Operation != "modified" {
		t.Errorf("expected modified, got %s", log.Entries[2].Operation)
	}
}

func TestAudit_FilterByOperation(t *testing.T) {
	changes := []Change{
		{Type: ChangeAdded, Path: "x", NewValue: 10},
		{Type: ChangeRemoved, Path: "y", OldValue: 20},
		{Type: ChangeModified, Path: "z", OldValue: 1, NewValue: 2},
	}
	opts := AuditOptions{Actor: "ci", Operations: []string{"added"}}
	log := Audit(changes, auditTime, opts)
	if len(log.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log.Entries))
	}
	if log.Entries[0].Path != "x" {
		t.Errorf("unexpected path: %s", log.Entries[0].Path)
	}
}

func TestAudit_ActorStored(t *testing.T) {
	changes := []Change{
		{Type: ChangeAdded, Path: "k", NewValue: "v"},
	}
	opts := AuditOptions{Actor: "alice"}
	log := Audit(changes, auditTime, opts)
	if log.Entries[0].Actor != "alice" {
		t.Errorf("expected actor alice, got %s", log.Entries[0].Actor)
	}
}

func TestAudit_TimestampPreserved(t *testing.T) {
	changes := []Change{
		{Type: ChangeModified, Path: "p", OldValue: "a", NewValue: "b"},
	}
	log := Audit(changes, auditTime, DefaultAuditOptions())
	if !log.Entries[0].Timestamp.Equal(auditTime) {
		t.Errorf("timestamp mismatch")
	}
}

func TestAudit_ValuesPreserved(t *testing.T) {
	changes := []Change{
		{Type: ChangeModified, Path: "val", OldValue: "old", NewValue: "new"},
	}
	log := Audit(changes, auditTime, DefaultAuditOptions())
	e := log.Entries[0]
	if e.OldValue != "old" || e.NewValue != "new" {
		t.Errorf("values not preserved: old=%v new=%v", e.OldValue, e.NewValue)
	}
}
