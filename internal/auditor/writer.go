package auditor

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// WriteReport writes a human-readable audit report to w.
func WriteReport(w io.Writer, entries []Entry) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "No audit entries found.")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tBASE\tCOMP FILES\tISSUES")
	fmt.Fprintln(tw, "---------\t----\t----------\t------")

	// Sort by timestamp ascending.
	sorted := make([]Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.Before(sorted[j].Timestamp)
	})

	for _, e := range sorted {
		issues := "none"
		if e.HadIssues {
			issues = "yes"
		}
		compStr := joinStrings(e.CompFiles)
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.BaseFile,
			compStr,
			issues,
		)
	}
	return tw.Flush()
}

// IssueCount returns the number of entries that had issues.
func IssueCount(entries []Entry) int {
	count := 0
	for _, e := range entries {
		if e.HadIssues {
			count++
		}
	}
	return count
}

func joinStrings(ss []string) string {
	if len(ss) == 0 {
		return ""
	}
	out := ss[0]
	for _, s := range ss[1:] {
		out += "," + s
	}
	return out
}
