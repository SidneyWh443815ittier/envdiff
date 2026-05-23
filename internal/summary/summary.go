package summary

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/comparator"
)

// Stats holds aggregate counts from a comparison result.
type Stats struct {
	Missing    int
	Extra      int
	Mismatched int
	Total      int
}

// Compute derives Stats from a comparator.Result.
func Compute(r comparator.Result) Stats {
	return Stats{
		Missing:    len(r.Missing),
		Extra:      len(r.Extra),
		Mismatched: len(r.Mismatched),
		Total:      len(r.Missing) + len(r.Extra) + len(r.Mismatched),
	}
}

// Clean reports whether the result has no differences.
func (s Stats) Clean() bool {
	return s.Total == 0
}

// Write prints a human-readable summary line to w.
func Write(w io.Writer, s Stats) {
	if s.Clean() {
		fmt.Fprintln(w, "✔ No differences found.")
		return
	}

	parts := []string{}
	if s.Missing > 0 {
		parts = append(parts, fmt.Sprintf("%d missing", s.Missing))
	}
	if s.Extra > 0 {
		parts = append(parts, fmt.Sprintf("%d extra", s.Extra))
	}
	if s.Mismatched > 0 {
		parts = append(parts, fmt.Sprintf("%d mismatched", s.Mismatched))
	}

	fmt.Fprintf(w, "✖ %d issue(s) found: %s\n", s.Total, strings.Join(parts, ", "))
}
