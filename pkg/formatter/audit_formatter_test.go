package formatter

import (
	"strings"
	"testing"
	"time"

	"github.com/szhekpisov/diffyml/pkg/diff"
)

var testAuditTime = time.Date(2024, 6, 1, 9, 0, 0, 0, time.UTC)

func TestAuditFormatter_NoEntries(t *testing.T) {
	var sb strings.Builder
	f := NewAuditFormatter(&sb)
	if err := f.Format(diff.AuditLog{}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sb.String(), "no entries") {
		t.Errorf("expected 'no entries', got: %s", sb.String())
	}
}

func TestAuditFormatter_HeaderPresent(t *testing.T) {
	log := diff.AuditLog{
		Entries: []diff.AuditEntry{
			{Timestamp: testAuditTime, Operation: "added", Path: "x", NewValue: 1, Actor: "bot"},
		},
	}
	var sb strings.Builder
	f := NewAuditFormatter(&sb)
	if err := f.Format(log); err != nil {
		t.Fatal(err)
	}
	out := sb.String()
	if !strings.Contains(out, "audit log") {
		t.Errorf("expected header, got: %s", out)
	}
}

func TestAuditFormatter_AddedLine(t *testing.T) {
	log := diff.AuditLog{
		Entries: []diff.AuditEntry{
			{Timestamp: testAuditTime, Operation: "added", Path: "foo", NewValue: "bar", Actor: "ci"},
		},
	}
	var sb strings.Builder
	f := NewAuditFormatter(&sb)
	_ = f.Format(log)
	out := sb.String()
	if !strings.Contains(out, "ADDED") || !strings.Contains(out, "foo") || !strings.Contains(out, "bar") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestAuditFormatter_RemovedLine(t *testing.T) {
	log := diff.AuditLog{
		Entries: []diff.AuditEntry{
			{Timestamp: testAuditTime, Operation: "removed", Path: "gone", OldValue: 42, Actor: "user"},
		},
	}
	var sb strings.Builder
	f := NewAuditFormatter(&sb)
	_ = f.Format(log)
	out := sb.String()
	if !strings.Contains(out, "REMOVED") || !strings.Contains(out, "gone") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestAuditFormatter_ModifiedLine(t *testing.T) {
	log := diff.AuditLog{
		Entries: []diff.AuditEntry{
			{Timestamp: testAuditTime, Operation: "modified", Path: "cfg.timeout", OldValue: 30, NewValue: 60, Actor: "admin"},
		},
	}
	var sb strings.Builder
	f := NewAuditFormatter(&sb)
	_ = f.Format(log)
	out := sb.String()
	if !strings.Contains(out, "MODIFIED") || !strings.Contains(out, "cfg.timeout") {
		t.Errorf("unexpected output: %s", out)
	}
	if !strings.Contains(out, "30") || !strings.Contains(out, "60") {
		t.Errorf("expected old/new values in output: %s", out)
	}
}
