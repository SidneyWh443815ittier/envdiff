package comparator

// Mismatch holds the base and target values for a key that exists in both
// environments but with differing values.
type Mismatch struct {
	Base   string
	Target string
}

// Result holds the outcome of comparing two environment maps.
type Result struct {
	// Missing contains keys present in base but absent in target.
	Missing map[string]string
	// Extra contains keys present in target but absent in base.
	Extra map[string]string
	// Mismatched contains keys present in both but with different values.
	Mismatched map[string]Mismatch
}

// IsClean returns true when there are no differences between the two env maps.
func (r Result) IsClean() bool {
	return len(r.Missing) == 0 && len(r.Extra) == 0 && len(r.Mismatched) == 0
}

// Compare compares base and target environment maps and returns a Result
// describing any differences found.
func Compare(base, target map[string]string) Result {
	result := Result{
		Missing:    make(map[string]string),
		Extra:      make(map[string]string),
		Mismatched: make(map[string]Mismatch),
	}

	for k, baseVal := range base {
		targetVal, ok := target[k]
		if !ok {
			result.Missing[k] = baseVal
			continue
		}
		if baseVal != targetVal {
			result.Mismatched[k] = Mismatch{Base: baseVal, Target: targetVal}
		}
	}

	for k, targetVal := range target {
		if _, ok := base[k]; !ok {
			result.Extra[k] = targetVal
		}
	}

	return result
}
