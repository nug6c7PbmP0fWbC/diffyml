// Package diff provides utilities for comparing, transforming, and annotating
// YAML structures.
//
// # Pipeline
//
// Pipeline allows you to compose multiple diff transformation steps into a
// single reusable chain.  Each step is a plain function with the signature:
//
//	func(changes []Change) ([]Change, error)
//
// Steps are executed in the order they are registered.  If any step returns a
// non-nil error the pipeline stops immediately and propagates that error to the
// caller.
//
// # Built-in steps
//
// Several helper constructors wrap existing diff operations as pipeline steps:
//
//   - IgnoreStep  – drops changes that match IgnoreRule patterns.
//   - RedactStep  – replaces sensitive values with a placeholder.
//   - AnnotateStep – attaches metadata to matching changes.
//
// # Example
//
//	p := diff.NewPipeline().
//	    Use(diff.IgnoreStep(ignoreRules)).
//	    Use(diff.RedactStep(redactRules)).
//	    Use(diff.AnnotateStep(annotateRules))
//
//	result, err := p.Run(changes)
package diff
