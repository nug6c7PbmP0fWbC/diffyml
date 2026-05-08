package diff

import (
	"testing"
)

func TestTruncate_Empty(t *testing.T) {
	out := Truncate(nil, DefaultTruncateOptions())
	if len(out) != 0 {
		t.Fatalf("expected empty slice, got %d items", len(out))
	}
}

func TestTruncate_ShortValuesUnchanged(t *testing.T) {
	changes := []Change{
		{Path: "key", Type: ChangeModified, Before: "hello", After: "world"},
	}
	out := Truncate(changes, DefaultTruncateOptions())
	if out[0].After != "world" {
		t.Fatalf("expected 'world', got %q", out[0].After)
	}
}

func TestTruncate_LongValueTrimmed(t *testing.T) {
	long := "abcdefghij" // 10 chars
	opts := TruncateOptions{MaxLength: 5, Ellipsis: "..."}
	changes := []Change{
		{Path: "k", Type: ChangeModified, Before: "x", After: long},
	}
	out := Truncate(changes, opts)
	got, ok := out[0].After.(string)
	if !ok {
		t.Fatal("After should be string")
	}
	if got != "abcde..." {
		t.Fatalf("expected 'abcde...', got %q", got)
	}
}

func TestTruncate_BeforeAlsoTrimmed(t *testing.T) {
	long := "0123456789"
	opts := TruncateOptions{MaxLength: 4, Ellipsis: "~~", OnlyValues: false}
	changes := []Change{
		{Path: "k", Type: ChangeModified, Before: long, After: "short"},
	}
	out := Truncate(changes, opts)
	got, ok := out[0].Before.(string)
	if !ok {
		t.Fatal("Before should be string")
	}
	if got != "0123~~" {
		t.Fatalf("expected '0123~~', got %q", got)
	}
}

func TestTruncate_OnlyValues_SkipsBefore(t *testing.T) {
	long := "abcdefghij"
	opts := TruncateOptions{MaxLength: 3, Ellipsis: ".", OnlyValues: true}
	changes := []Change{
		{Path: "k", Type: ChangeModified, Before: long, After: long},
	}
	out := Truncate(changes, opts)
	if out[0].Before.(string) != long {
		t.Fatalf("Before should be untouched when OnlyValues=true")
	}
	if out[0].After.(string) != "abc." {
		t.Fatalf("expected 'abc.', got %q", out[0].After.(string))
	}
}

func TestTruncate_NonStringUnchanged(t *testing.T) {
	opts := TruncateOptions{MaxLength: 2, Ellipsis: "..."}
	changes := []Change{
		{Path: "num", Type: ChangeModified, Before: 42, After: 99},
	}
	out := Truncate(changes, opts)
	if out[0].After.(int) != 99 {
		t.Fatalf("non-string value should be unchanged")
	}
}

func TestTruncate_ZeroMaxLengthNoOp(t *testing.T) {
	long := "this is a very long string value"
	opts := TruncateOptions{MaxLength: 0, Ellipsis: "..."}
	changes := []Change{
		{Path: "k", Type: ChangeModified, Before: long, After: long},
	}
	out := Truncate(changes, opts)
	if out[0].After.(string) != long {
		t.Fatal("MaxLength=0 should be a no-op")
	}
}

func TestTruncate_DefaultEllipsis(t *testing.T) {
	long := "abcdefghij"
	opts := TruncateOptions{MaxLength: 5, Ellipsis: ""} // empty → default "..."
	changes := []Change{
		{Path: "k", Type: ChangeAdded, Before: nil, After: long},
	}
	out := Truncate(changes, opts)
	got := out[0].After.(string)
	if got != "abcde..." {
		t.Fatalf("expected 'abcde...', got %q", got)
	}
}
