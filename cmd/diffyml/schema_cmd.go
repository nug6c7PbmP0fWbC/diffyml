package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/szhekpisov/diffyml/pkg/diff"
	"github.com/szhekpisov/diffyml/pkg/loader"
)

// schemaConfig is the structure read from a JSON schema rules file.
type schemaConfig struct {
	ProtectedPaths []string `json:"protected_paths"`
	RequiredOnAdd  []string `json:"required_on_add"`
	ForbidType     string   `json:"forbid_type"`
}

// runSchemaValidate loads two YAML files, computes their diff, reads a schema
// rules JSON file, and prints any violations to stderr, exiting with code 1
// when violations are found.
func runSchemaValidate(baseFile, headFile, rulesFile string) {
	base, err := loader.LoadFile(baseFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading base file: %v\n", err)
		os.Exit(1)
	}

	head, err := loader.LoadFile(headFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading head file: %v\n", err)
		os.Exit(1)
	}

	changes := diff.Compare(base, head)

	ruleData, err := os.ReadFile(rulesFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading rules file: %v\n", err)
		os.Exit(1)
	}

	var cfg schemaConfig
	if err := json.Unmarshal(ruleData, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing rules file: %v\n", err)
		os.Exit(1)
	}

	rule := diff.SchemaRule{
		ProtectedPaths: cfg.ProtectedPaths,
		RequiredOnAdd:  cfg.RequiredOnAdd,
	}

	switch cfg.ForbidType {
	case "added":
		rule.ForbidType = diff.ChangeAdded
		rule.ForbidTypeSet = true
	case "removed":
		rule.ForbidType = diff.ChangeRemoved
		rule.ForbidTypeSet = true
	case "modified":
		rule.ForbidType = diff.ChangeModified
		rule.ForbidTypeSet = true
	}

	violations := diff.ValidateSchema(changes, rule)
	if len(violations) == 0 {
		fmt.Println("schema validation passed: no violations found")
		return
	}

	fmt.Fprintf(os.Stderr, "schema validation failed: %d violation(s)\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(os.Stderr, "  - %s\n", v.Error())
	}
	os.Exit(1)
}
