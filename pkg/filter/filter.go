package filter

import (
	"strings"

	"github.com/diffyml/diffyml/pkg/diff"
)

// Options holds filtering configuration.
type Options struct {
	// Paths restricts output to changes whose path has one of these prefixes.
	Paths []string
	// Types restricts output to specific change types (added, removed, modified).
	Types []diff.ChangeType
}

// Apply returns only the changes that match the given options.
// If an option slice is empty it is treated as "match all".
func Apply(changes []diff.Change, opts Options) []diff.Change {
	var result []diff.Change
	for _, c := range changes {
		if !matchPath(c.Path, opts.Paths) {
			continue
		}
		if !matchType(c.Type, opts.Types) {
			continue
		}
		result = append(result, c)
	}
	return result
}

func matchPath(path string, prefixes []string) bool {
	if len(prefixes) == 0 {
		return true
	}
	for _, p := range prefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func matchType(ct diff.ChangeType, types []diff.ChangeType) bool {
	if len(types) == 0 {
		return true
	}
	for _, t := range types {
		if ct == t {
			return true
		}
	}
	return false
}
