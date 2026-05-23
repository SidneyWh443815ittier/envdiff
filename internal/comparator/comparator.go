package comparator

// Result holds the outcome of comparing two env maps.
type Result struct {
	Missing  []string          // keys present in base but absent in target
	Extra    []string          // keys present in target but absent in base
	Mismatch map[string]Diff   // keys present in both but with different values
}

// Diff captures the differing values for a single key.
type Diff struct {
	Base   string
	Target string
}

// Compare compares a base env map against a target env map and returns a
// Result describing missing keys, extra keys, and value mismatches.
func Compare(base, target map[string]string) Result {
	result := Result{
		Mismatch: make(map[string]Diff),
	}

	for key, baseVal := range base {
		targetVal, ok := target[key]
		if !ok {
			result.Missing = append(result.Missing, key)
			continue
		}
		if baseVal != targetVal {
			result.Mismatch[key] = Diff{Base: baseVal, Target: targetVal}
		}
	}

	for key := range target {
		if _, ok := base[key]; !ok {
			result.Extra = append(result.Extra, key)
		}
	}

	return result
}

// HasDiff returns true when any differences were found.
func (r Result) HasDiff() bool {
	return len(r.Missing) > 0 || len(r.Extra) > 0 || len(r.Mismatch) > 0
}
