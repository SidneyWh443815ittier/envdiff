// Package linter checks .env files for common style and correctness issues.
package linter

import (
	"fmt"
	"strings"
)

// Issue represents a single linting problem found in an env file.
type Issue struct {
	Line    int
	Key     string
	Message string
}

// Options controls which lint rules are enabled.
type Options struct {
	CheckDuplicates   bool
	CheckEmptyValues  bool
	CheckQuoting      bool
	CheckKeyFormat    bool
}

// DefaultOptions returns Options with all rules enabled.
func DefaultOptions() Options {
	return Options{
		CheckDuplicates:  true,
		CheckEmptyValues: true,
		CheckQuoting:     true,
		CheckKeyFormat:   true,
	}
}

// Lint analyses the raw lines of an env file and returns any issues found.
func Lint(lines []string, opts Options) []Issue {
	var issues []Issue
	seen := make(map[string]int)

	for i, raw := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		eqIdx := strings.Index(trimmed, "=")
		if eqIdx < 0 {
			issues = append(issues, Issue{Line: lineNum, Message: "malformed line: missing '='"})
			continue
		}

		key := strings.TrimSpace(trimmed[:eqIdx])
		val := trimmed[eqIdx+1:]

		if opts.CheckDuplicates {
			if prev, ok := seen[key]; ok {
				issues = append(issues, Issue{
					Line:    lineNum,
					Key:     key,
					Message: fmt.Sprintf("duplicate key (first seen on line %d)", prev),
				})
			} else {
				seen[key] = lineNum
			}
		}

		if opts.CheckEmptyValues && strings.TrimSpace(val) == "" {
			issues = append(issues, Issue{Line: lineNum, Key: key, Message: "empty value"})
		}

		if opts.CheckQuoting {
			if strings.Contains(val, " ") && !strings.HasPrefix(val, "\"") && !strings.HasPrefix(val, "'") {
				issues = append(issues, Issue{Line: lineNum, Key: key, Message: "value with spaces should be quoted"})
			}
		}

		if opts.CheckKeyFormat {
			if key != strings.ToUpper(key) || strings.Contains(key, " ") {
				issues = append(issues, Issue{Line: lineNum, Key: key, Message: "key should be UPPER_SNAKE_CASE"})
			}
		}
	}

	return issues
}
