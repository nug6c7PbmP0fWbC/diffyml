package diff

// PartitionOptions controls how changes are split into partitions.
type PartitionOptions struct {
	// MaxPerPartition is the maximum number of changes per partition.
	// Zero means no limit.
	MaxPerPartition int
	// GroupByType splits added/removed/modified into separate partitions.
	GroupByType bool
}

// DefaultPartitionOptions returns sensible defaults.
func DefaultPartitionOptions() PartitionOptions {
	return PartitionOptions{
		MaxPerPartition: 50,
		GroupByType:     false,
	}
}

// Partition splits a slice of Changes into buckets according to opts.
// Each bucket is a []Change. If both MaxPerPartition and GroupByType are
// active, changes are first grouped by type, then chunked by size.
func Partition(changes []Change, opts PartitionOptions) [][]Change {
	if len(changes) == 0 {
		return nil
	}

	var groups [][]Change

	if opts.GroupByType {
		order := []ChangeType{Added, Removed, Modified}
		for _, ct := range order {
			var bucket []Change
			for _, c := range changes {
				if c.Type == ct {
					bucket = append(bucket, c)
				}
			}
			if len(bucket) > 0 {
				groups = append(groups, bucket)
			}
		}
	} else {
		groups = [][]Change{changes}
	}

	if opts.MaxPerPartition <= 0 {
		return groups
	}

	var result [][]Change
	for _, g := range groups {
		for i := 0; i < len(g); i += opts.MaxPerPartition {
			end := i + opts.MaxPerPartition
			if end > len(g) {
				end = len(g)
			}
			result = append(result, g[i:end])
		}
	}
	return result
}
