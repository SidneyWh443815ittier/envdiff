// Package patcher applies a set of key-value changes to an existing .env file,
// preserving comments, blank lines, and ordering of unchanged entries.
package patcher

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Options controls patcher behaviour.
type Options struct {
	// CreateMissing adds keys that do not exist in the target file.
	CreateMissing bool
	// QuoteValues wraps new/updated values in double-quotes.
	QuoteValues bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		CreateMissing: true,
		QuoteValues:   false,
	}
}

// Patch reads path, applies changes, and writes the result back to path.
// Keys present in changes but absent from the file are appended when
// CreateMissing is true. Existing comments and blank lines are preserved.
func Patch(path string, changes map[string]string, opts Options) error {
	lines, err := readLines(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("patcher: read %s: %w", path, err)
	}

	applied := make(map[string]bool)
	var out []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			out = append(out, line)
			continue
		}
		idx := strings.IndexByte(trimmed, '=')
		if idx < 0 {
			out = append(out, line)
			continue
		}
		key := strings.TrimSpace(trimmed[:idx])
		if val, ok := changes[key]; ok {
			out = append(out, formatLine(key, val, opts.QuoteValues))
			applied[key] = true
		} else {
			out = append(out, line)
		}
	}

	if opts.CreateMissing {
		for k, v := range changes {
			if !applied[k] {
				out = append(out, formatLine(k, v, opts.QuoteValues))
			}
		}
	}

	return writeLines(path, out)
}

func formatLine(key, value string, quote bool) string {
	if quote {
		return fmt.Sprintf(`%s="%s"`, key, value)
	}
	return fmt.Sprintf("%s=%s", key, value)
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines, sc.Err()
}

func writeLines(path string, lines []string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("patcher: create %s: %w", path, err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, l := range lines {
		if _, err := fmt.Fprintln(w, l); err != nil {
			return err
		}
	}
	return w.Flush()
}
