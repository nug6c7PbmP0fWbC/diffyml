package diff

import (
	"testing"
	"time"
)

func makeLoader(data map[string]interface{}) func() (map[string]interface{}, error) {
	return func() (map[string]interface{}, error) {
		return data, nil
	}
}

func TestWatcher_DetectsChange(t *testing.T) {
	a := map[string]interface{}{"key": "old"}
	b := map[string]interface{}{"key": "new"}

	w := NewWatcher(20 * time.Millisecond)
	go w.Watch("a.yml", "b.yml", makeLoader(a), makeLoader(b))
	defer w.Stop()

	select {
	case evt := <-w.Events():
		if len(evt.Changes) == 0 {
			t.Fatal("expected changes, got none")
		}
		if evt.FileA != "a.yml" || evt.FileB != "b.yml" {
			t.Errorf("unexpected file names: %s, %s", evt.FileA, evt.FileB)
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatcher_NoEventWhenIdentical(t *testing.T) {
	data := map[string]interface{}{"key": "value"}

	w := NewWatcher(20 * time.Millisecond)
	go w.Watch("a.yml", "b.yml", makeLoader(data), makeLoader(data))
	defer w.Stop()

	select {
	case evt := <-w.Events():
		t.Fatalf("unexpected event: %+v", evt)
	case <-time.After(150 * time.Millisecond):
		// expected: no events
	}
}

func TestWatcher_NoDuplicateEvents(t *testing.T) {
	a := map[string]interface{}{"x": 1}
	b := map[string]interface{}{"x": 2}

	w := NewWatcher(20 * time.Millisecond)
	go w.Watch("a.yml", "b.yml", makeLoader(a), makeLoader(b))
	defer w.Stop()

	count := 0
	timeout := time.After(200 * time.Millisecond)
loop:
	for {
		select {
		case <-w.Events():
			count++
		case <-timeout:
			break loop
		}
	}

	if count != 1 {
		t.Errorf("expected exactly 1 event, got %d", count)
	}
}

func TestHashChanges_Deterministic(t *testing.T) {
	changes := []Change{
		{Path: "a.b", Type: ChangeModified, OldValue: "x", NewValue: "y"},
	}
	h1 := hashChanges(changes)
	h2 := hashChanges(changes)
	if h1 != h2 {
		t.Errorf("hash not deterministic: %s != %s", h1, h2)
	}
}
