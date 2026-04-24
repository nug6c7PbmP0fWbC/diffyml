package diff

import (
	"testing"
)

func TestAnnotate_NoRules(t *testing.T) {
	changes := []Change{
		{Path: "a.b", Type: ChangeTypeAdded},
	}
	result := Annotate(changes, nil)
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %v", result)
	}
}

func TestAnnotate_MatchAll(t *testing.T) {
	changes := []Change{
		{Path: "x.y", Type: ChangeTypeModified},
		{Path: "a.b", Type: ChangeTypeAdded},
	}
	rules := []AnnotationRule{
		{PathPrefix: "", ChangeType: "", Severity: "info", Message: "any change"},
	}
	result := Annotate(changes, rules)
	if len(result) != 2 {
		t.Fatalf("expected 2 annotated paths, got %d", len(result))
	}
}

func TestAnnotate_FilterByType(t *testing.T) {
	changes := []Change{
		{Path: "a", Type: ChangeTypeAdded},
		{Path: "b", Type: ChangeTypeRemoved},
	}
	rules := []AnnotationRule{
		{ChangeType: "removed", Severity: "warning", Message: "key removed"},
	}
	result := Annotate(changes, rules)
	if _, ok := result["b"]; !ok {
		t.Fatal("expected annotation on 'b'")
	}
	if _, ok := result["a"]; ok {
		t.Fatal("did not expect annotation on 'a'")
	}
}

func TestAnnotate_FilterByPathPrefix(t *testing.T) {
	changes := []Change{
		{Path: "service.port", Type: ChangeTypeModified},
		{Path: "database.host", Type: ChangeTypeModified},
	}
	rules := []AnnotationRule{
		{PathPrefix: "service", Severity: "critical", Message: "service config changed"},
	}
	result := Annotate(changes, rules)
	if len(result["service.port"]) != 1 {
		t.Fatal("expected annotation on 'service.port'")
	}
	if len(result["database.host"]) != 0 {
		t.Fatal("did not expect annotation on 'database.host'")
	}
}

func TestAnnotate_MultipleRulesOnSamePath(t *testing.T) {
	changes := []Change{
		{Path: "app.secret", Type: ChangeTypeModified},
	}
	rules := []AnnotationRule{
		{PathPrefix: "app", Severity: "info", Message: "app changed"},
		{PathPrefix: "app.secret", Severity: "critical", Message: "secret rotated"},
	}
	result := Annotate(changes, rules)
	if len(result["app.secret"]) != 2 {
		t.Fatalf("expected 2 annotations, got %d", len(result["app.secret"]))
	}
}

func TestMatchPrefix_Empty(t *testing.T) {
	if !matchPrefix("any.path", "") {
		t.Fatal("empty prefix should match anything")
	}
}

func TestMatchPrefix_Exact(t *testing.T) {
	if !matchPrefix("foo.bar", "foo") {
		t.Fatal("prefix 'foo' should match 'foo.bar'")
	}
	if matchPrefix("baz", "foo") {
		t.Fatal("prefix 'foo' should not match 'baz'")
	}
}
