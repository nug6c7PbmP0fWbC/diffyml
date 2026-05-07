package diff

// ClusterOptions controls how changes are clustered.
type ClusterOptions struct {
	// MaxDistance is the maximum edit-distance between two paths
	// for them to be considered part of the same cluster.
	MaxDistance int
	// MinSize is the minimum number of changes required to form a cluster.
	MinSize int
}

// DefaultClusterOptions returns sensible defaults.
func DefaultClusterOptions() ClusterOptions {
	return ClusterOptions{
		MaxDistance: 2,
		MinSize:     2,
	}
}

// Cluster groups related changes into named buckets based on shared path
// prefixes and edit-distance heuristics.
func Cluster(changes []Change, opts ClusterOptions) map[string][]Change {
	result := make(map[string][]Change)

	for _, c := range changes {
		key := clusterKey(c.Path, opts.MaxDistance, changes)
		result[key] = append(result[key], c)
	}

	// Remove clusters that are too small.
	if opts.MinSize > 1 {
		for k, v := range result {
			if len(v) < opts.MinSize {
				delete(result, k)
			}
		}
	}

	return result
}

// clusterKey returns the representative key for a change path.
func clusterKey(path string, maxDist int, all []Change) string {
	best := path
	bestDist := maxDist + 1

	for _, c := range all {
		if c.Path == path {
			continue
		}
		d := editDistance(path, c.Path)
		if d < bestDist {
			bestDist = d
			best = c.Path
		}
	}

	if bestDist <= maxDist && best < path {
		return best
	}
	return path
}

// editDistance computes the Levenshtein distance between two strings.
func editDistance(a, b string) int {
	la, lb := len(a), len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
	row := make([]int, lb+1)
	for j := 0; j <= lb; j++ {
		row[j] = j
	}
	for i := 1; i <= la; i++ {
		prev := row[0]
		row[0] = i
		for j := 1; j <= lb; j++ {
			tmp := row[j]
			if a[i-1] == b[j-1] {
				row[j] = prev
			} else {
				row[j] = 1 + min3(prev, row[j], row[j-1])
			}
			prev = tmp
		}
	}
	return row[lb]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
