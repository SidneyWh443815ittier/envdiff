package comparator

// ValuePair holds the base and comparison values for a mismatched key.
type ValuePair struct {
	Base string
	Comp string
}

// Result contains the full diff between two env maps.
type Result struct {
	// Missing keys are present in base but absent in comp.
	Missing []string
	// Extra keys are present in comp but absent in base.
	Extra []string
	// Mismatched keys exist in both but have different values.
	Mismatched map[string]ValuePair
}

// Clean returns true when there are no differences.
func (r Result) Clean() bool {
	return len(r.Missing) == 0 && len(r.Extra) == 0 && len(r.Mismatched) == 0
}

// Compare diffs base against comp and returns a Result.
func Compare(base, comp map[string]string) Result {
	r := Result{
		Mismatched: make(map[string]ValuePair),
	}

	for k, bv := range base {
		if cv, ok := comp[k]; !ok {
			r.Missing = append(r.Missing, k)
		} else if bv != cv {
			r.Mismatched[k] = ValuePair{Base: bv, Comp: cv}
		}
	}

	for k := range comp {
		if _, ok := base[k]; !ok {
			r.Extra = append(r.Extra, k)
		}
	}

	return r
}
