package diff

import (
	"testing"
)

func TestValidateSchema_NoViolations(t *testing.T) {
	changes := []Change{
		{Path: "app.name", Type: ChangeAdded, NewValue: "myapp"},
	}
	rule := SchemaRule{}
	violations := ValidateSchema(changes, rule)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestValidateSchema_ProtectedPathModified(t *testing.T) {
	changes := []Change{
		{Path: "app.version", Type: ChangeModified, OldValue: "1.0", NewValue: "2.0"},
	}
	rule := SchemaRule{ProtectedPaths: []string{"app.version"}}
	violations := ValidateSchema(changes, rule)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Path != "app.version" {
		t.Errorf("expected path app.version, got %s", violations[0].Path)
	}
}

func TestValidateSchema_ProtectedPathRemoved(t *testing.T) {
	changes := []Change{
		{Path: "db.host", Type: ChangeRemoved, OldValue: "localhost"},
	}
	rule := SchemaRule{ProtectedPaths: []string{"db.host"}}
	violations := ValidateSchema(changes, rule)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidateSchema_ForbidType(t *testing.T) {
	changes := []Change{
		{Path: "feature.flag", Type: ChangeRemoved, OldValue: true},
	}
	rule := SchemaRule{ForbidType: ChangeRemoved, ForbidTypeSet: true}
	violations := ValidateSchema(changes, rule)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidateSchema_RequiredOnAdd_Missing(t *testing.T) {
	changes := []Change{
		{Path: "service.name", Type: ChangeAdded, NewValue: "api"},
	}
	rule := SchemaRule{RequiredOnAdd: []string{"service.port"}}
	violations := ValidateSchema(changes, rule)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Path != "service.port" {
		t.Errorf("unexpected path: %s", violations[0].Path)
	}
}

func TestValidateSchema_RequiredOnAdd_Present(t *testing.T) {
	changes := []Change{
		{Path: "service.name", Type: ChangeAdded, NewValue: "api"},
		{Path: "service.port", Type: ChangeAdded, NewValue: 8080},
	}
	rule := SchemaRule{RequiredOnAdd: []string{"service.port"}}
	violations := ValidateSchema(changes, rule)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestSchemaViolation_Error(t *testing.T) {
	v := SchemaViolation{Path: "a.b", Message: "not allowed"}
	expected := "a.b: not allowed"
	if v.Error() != expected {
		t.Errorf("expected %q, got %q", expected, v.Error())
	}
}
