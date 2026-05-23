// Package exporter provides functionality to export diff results
// to various file formats such as JSON and plain text.
package exporter

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/user/envdiff/internal/comparator"
)

// Format represents the output format for export.
type Format string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"
)

// Options configures the exporter.
type Options struct {
	Format   Format
	FilePath string
}

// DefaultOptions returns sensible export defaults.
func DefaultOptions() Options {
	return Options{
		Format:   FormatText,
		FilePath: "envdiff-report.txt",
	}
}

// Export writes the comparison result to a file using the specified format.
func Export(result comparator.Result, opts Options) error {
	f, err := os.Create(opts.FilePath)
	if err != nil {
		return fmt.Errorf("exporter: create file %q: %w", opts.FilePath, err)
	}
	defer f.Close()

	switch opts.Format {
	case FormatJSON:
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		if err := enc.Encode(result); err != nil {
			return fmt.Errorf("exporter: encode json: %w", err)
		}
	case FormatText:
		if err := writeText(f, result); err != nil {
			return fmt.Errorf("exporter: write text: %w", err)
		}
	default:
		return fmt.Errorf("exporter: unknown format %q", opts.Format)
	}

	return nil
}

func writeText(f *os.File, result comparator.Result) error {
	if len(result.Missing) == 0 && len(result.Extra) == 0 && len(result.Mismatched) == 0 {
		_, err := fmt.Fprintln(f, "No differences found.")
		return err
	}

	missing := sortedKeys(result.Missing)
	for _, k := range missing {
		if _, err := fmt.Fprintf(f, "MISSING  %s\n", k); err != nil {
			return err
		}
	}

	extra := sortedKeys(result.Extra)
	for _, k := range extra {
		if _, err := fmt.Fprintf(f, "EXTRA    %s\n", k); err != nil {
			return err
		}
	}

	for _, k := range sortedMismatchKeys(result.Mismatched) {
		m := result.Mismatched[k]
		if _, err := fmt.Fprintf(f, "MISMATCH %s (base=%q, comp=%q)\n", k, m.Base, m.Comp); err != nil {
			return err
		}
	}

	return nil
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
