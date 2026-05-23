package reporter

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/envdiff/internal/comparator"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Report writes a human-readable diff report to the given writer.
func Report(w io.Writer, result comparator.Result, baseFile, targetFile string) {
	if result.IsClean() {
		fmt.Fprintf(w, "✅ No differences found between %s and %s\n", baseFile, targetFile)
		return
	}

	fmt.Fprintf(w, "🔍 Comparing %s → %s\n\n", baseFile, targetFile)

	if len(result.Missing) > 0 {
		keys := sortedKeys(result.Missing)
		fmt.Fprintf(w, "❌ Missing keys in %s (%d):\n", targetFile, len(keys))
		for _, k := range keys {
			fmt.Fprintf(w, "   - %s\n", k)
		}
		fmt.Fprintln(w)
	}

	if len(result.Extra) > 0 {
		keys := sortedKeys(result.Extra)
		fmt.Fprintf(w, "➕ Extra keys in %s (%d):\n", targetFile, len(keys))
		for _, k := range keys {
			fmt.Fprintf(w, "   + %s\n", k)
		}
		fmt.Fprintln(w)
	}

	if len(result.Mismatched) > 0 {
		keys := sortedMismatchKeys(result.Mismatched)
		fmt.Fprintf(w, "⚠️  Mismatched values (%d):\n", len(keys))
		for _, k := range keys {
			m := result.Mismatched[k]
			fmt.Fprintf(w, "   ~ %s: %q → %q\n", k, m.Base, m.Target)
		}
		fmt.Fprintln(w)
	}
}

// ExitCode returns 1 if there are any differences, 0 otherwise.
// Useful for CI integration.
func ExitCode(result comparator.Result) int {
	if result.IsClean() {
		return 0
	}
	return 1
}

// ReportToStdout writes the report to os.Stdout.
func ReportToStdout(result comparator.Result, baseFile, targetFile string) {
	Report(os.Stdout, result, baseFile, targetFile)
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedMismatchKeys(m map[string]comparator.Mismatch) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
