package diff

// Intersect returns only the changes whose paths appear in ALL provided
// change slices. The values (Before/After) are taken from the first slice.
// If fewer than two slices are provided the first slice is returned as-is.
func Intersect(sets ...[]Change) []Change {
	if len(sets) == 0 {
		return nil
	}
	if len(sets) == 1 {
		return sets[0]
	}

	// Build path sets for every slice beyond the first.
	presence := make([]map[string]struct{}, len(sets)-1)
	for i, s := range sets[1:] {
		pm := make(map[string]struct{}, len(s))
		for _, c := range s {
			pm[c.Path] = struct{}{}
		}
		presence[i] = pm
	}

	var result []Change
	for _, c := range sets[0] {
		inAll := true
		for _, pm := range presence {
			if _, ok := pm[c.Path]; !ok {
				inAll = false
				break
			}
		}
		if inAll {
			result = append(result, c)
		}
	}
	return result
}
