package cmd

// CompareResult holds the outcome of comparing two sets of env entries.
type CompareResult struct {
	OnlyInA   []string
	OnlyInB   []string
	Different []string
	Same      []string
}

// compareStores compares two maps of env entries and returns a CompareResult.
func compareStores(a, b map[string]string) CompareResult {
	result := CompareResult{}

	for k, va := range a {
		if vb, ok := b[k]; !ok {
			result.OnlyInA = append(result.OnlyInA, k)
		} else if va != vb {
			result.Different = append(result.Different, k)
		} else {
			result.Same = append(result.Same, k)
		}
	}

	for k := range b {
		if _, ok := a[k]; !ok {
			result.OnlyInB = append(result.OnlyInB, k)
		}
	}

	return result
}
