/*
Package diff provides the Mask function for hiding sensitive values in change sets.

# Overview

Mask replaces the OldValue and NewValue of matched changes with a configurable
placeholder (default "***"), while preserving the Path and Type fields so that
audit trails remain intact.

# Usage

	opts := diff.DefaultMaskOptions()
	opts.Paths = []string{"db.password", "auth.token"}
	opts.PrefixMatch = true   // match any path that starts with the listed prefixes
	opts.Placeholder = "[HIDDEN]"

	masked := diff.Mask(changes, opts)

# Notes

  - Mask does not remove changes; it only obscures their values.
  - Use alongside Redact when you need to fully suppress certain paths.
  - The MaskFormatter in pkg/formatter renders masked changes and appends
    a "[masked]" annotation so readers know a value was hidden.
*/
package diff
