// Package patch provides utilities for applying a set of diff.Change entries
// to a YAML document represented as a map[string]interface{}.
//
// It supports two directions:
//
//   - Forward: transforms a document from its old state to the new state
//     by applying added/modified values and removing deleted keys.
//
//   - Reverse: undoes a set of changes, restoring the document to its
//     previous state.
//
// # Example
//
//	changes, _ := diff.Compare(oldDoc, newDoc)
//	err := patch.Apply(oldDoc, changes, patch.Forward)
//	// oldDoc now reflects newDoc
package patch
