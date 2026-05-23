package diff

import (
	"testing"
)

func TestHighlight_Empty(t *testing.T) {
	out := Highlight(nil, DefaultHighlightOptions())
	if len(out) != 0 {
		t.Fatalf("expected empty slice, got %d", len(out))
	}
}

func TestHighlight_Added(t *testing.T) {
	changes := []Change{
		{Type: ChangeAdded, Path: "a.b"},
	}
	opts := DefaultHighlightOptions()
	out := Highlight(changes, opts)
	want := "[+] a.b"
	if got := out[0].Metadata[opts.Tag]; got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestHighlight_Removed(t *testing.T) {
	changes := []Change{
		{Type: ChangeRemoved, Path: "x.y"},
	}
	opts := DefaultHighlightOptions()
	out := Highlight(changes, opts)
	want := "[-] x.y"
	if got := out[0].Metadata[opts.Tag]; got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestHighlight_Modified(t *testing.T) {
	changes := []Change{
		{Type: ChangeModified, Path: "config.timeout"},
	}
	opts := DefaultHighlightOptions()
	out := Highlight(changes, opts)
	want := "[~] config.timeout"
	if got := out[0].Metadata[opts.Tag]; got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestHighlight_PreservesExistingMetadata(t *testing.T) {
	changes := []Change{
		{Type: ChangeAdded, Path: "k", Metadata: map[string]string{"owner": "team-a"}},
	}
	opts := DefaultHighlightOptions()
	out := Highlight(changes, opts)
	if out[0].Metadata["owner"] != "team-a" {
		t.Error("existing metadata key 'owner' was lost")
	}
	if out[0].Metadata[opts.Tag] == "" {
		t.Error("highlight tag was not set")
	}
}

func TestHighlight_DoesNotMutateOriginal(t *testing.T) {
	orig := []Change{
		{Type: ChangeAdded, Path: "p"},
	}
	opts := DefaultHighlightOptions()
	Highlight(orig, opts)
	if orig[0].Metadata != nil {
		t.Error("original change was mutated")
	}
}

func TestHighlight_CustomOptions(t *testing.T) {
	changes := []Change{
		{Type: ChangeAdded, Path: "z"},
	}
	opts := HighlightOptions{
		PrefixAdded:    "ADD",
		PrefixRemoved:  "DEL",
		PrefixModified: "MOD",
		Tag:            "hl",
	}
	out := Highlight(changes, opts)
	want := "ADD z"
	if got := out[0].Metadata["hl"]; got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
