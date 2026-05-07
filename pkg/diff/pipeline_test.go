package diff

import (
	"errors"
	"testing"
)

func TestPipeline_Empty(t *testing.T) {
	changes := []Change{
		{Path: "a.b", Type: ChangeTypeAdded, NewValue: 1},
	}
	out, err := NewPipeline().Run(changes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 change, got %d", len(out))
	}
}

func TestPipeline_StepsAreChained(t *testing.T) {
	double := func(changes []Change) ([]Change, error) {
		out := make([]Change, 0, len(changes)*2)
		out = append(out, changes...)
		out = append(out, changes...)
		return out, nil
	}
	changes := []Change{{Path: "x", Type: ChangeTypeAdded}}
	out, err := NewPipeline().Use(double, double).Run(changes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 1 → double → 2 → double → 4
	if len(out) != 4 {
		t.Fatalf("expected 4 changes, got %d", len(out))
	}
}

func TestPipeline_StopsOnError(t *testing.T) {
	sentinel := errors.New("step error")
	fail := func(changes []Change) ([]Change, error) {
		return nil, sentinel
	}
	neverReached := func(changes []Change) ([]Change, error) {
		t.Error("this step should not have been called")
		return changes, nil
	}
	_, err := NewPipeline().Use(fail, neverReached).Run([]Change{})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestIgnoreStep(t *testing.T) {
	changes := []Change{
		{Path: "secret.key", Type: ChangeTypeModified},
		{Path: "public.value", Type: ChangeTypeModified},
	}
	rules := []IgnoreRule{{Path: "secret.key"}}
	out, err := NewPipeline().Use(IgnoreStep(rules)).Run(changes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 change after ignore, got %d", len(out))
	}
	if out[0].Path != "public.value" {
		t.Errorf("unexpected path: %s", out[0].Path)
	}
}

func TestRedactStep(t *testing.T) {
	changes := []Change{
		{Path: "db.password", Type: ChangeTypeModified, OldValue: "hunter2", NewValue: "s3cr3t"},
	}
	rules := []RedactRule{{Path: "db.password"}}
	out, err := NewPipeline().Use(RedactStep(rules)).Run(changes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].NewValue != "***" {
		t.Errorf("expected redacted value, got %v", out[0].NewValue)
	}
}

func TestAnnotateStep(t *testing.T) {
	changes := []Change{
		{Path: "app.version", Type: ChangeTypeModified},
	}
	rules := []AnnotateRule{
		{PathPrefix: "app", Annotations: map[string]string{"owner": "platform"}},
	}
	out, err := NewPipeline().Use(AnnotateStep(rules)).Run(changes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Annotations["owner"] != "platform" {
		t.Errorf("annotation not set: %v", out[0].Annotations)
	}
}
