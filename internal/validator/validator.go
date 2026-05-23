// Package validator provides validation logic for .env file keys and values.
package validator

import (
	"fmt"
	"strings"
)

// Issue represents a single validation problem found in an env map.
type Issue struct {
	Key     string
	Message string
}

// String returns a human-readable representation of the issue.
func (i Issue) String() string {
	return fmt.Sprintf("[%s] %s", i.Key, i.Message)
}

// Options controls which validations are performed.
type Options struct {
	// WarnEmptyValues flags keys whose value is an empty string.
	WarnEmptyValues bool
	// WarnKeyFormat flags keys that contain lowercase letters (convention: ALL_CAPS).
	WarnKeyFormat bool
}

// DefaultOptions returns Options with all validations enabled.
func DefaultOptions() Options {
	return Options{
		WarnEmptyValues: true,
		WarnKeyFormat:   true,
	}
}

// Validate inspects the provided env map and returns a slice of Issues.
// An empty slice means no problems were found.
func Validate(env map[string]string, opts Options) []Issue {
	var issues []Issue

	for k, v := range env {
		if opts.WarnKeyFormat && !isUpperSnakeCase(k) {
			issues = append(issues, Issue{
				Key:     k,
				Message: "key should be UPPER_SNAKE_CASE",
			})
		}
		if opts.WarnEmptyValues && strings.TrimSpace(v) == "" {
			issues = append(issues, Issue{
				Key:     k,
				Message: "value is empty",
			})
		}
	}

	return issues
}

// isUpperSnakeCase returns true when every character in s is an uppercase
// letter, digit, or underscore, and the key is non-empty.
func isUpperSnakeCase(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !((r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}
	return true
}
