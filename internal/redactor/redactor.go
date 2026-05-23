// Package redactor masks sensitive values in env maps before output.
package redactor

import "strings"

// DefaultSensitivePatterns contains common substrings that indicate a key holds a secret.
var DefaultSensitivePatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"CREDENTIALS",
	"AUTH",
}

const redactedValue = "***REDACTED***"

// Options controls redaction behaviour.
type Options struct {
	// Patterns is the list of key substrings that trigger redaction.
	// Matching is case-insensitive.
	Patterns []string
	// Enabled toggles redaction on or off.
	Enabled bool
}

// DefaultOptions returns an Options with the default sensitive patterns enabled.
func DefaultOptions() Options {
	return Options{
		Patterns: DefaultSensitivePatterns,
		Enabled:  true,
	}
}

// Redact returns a copy of env with sensitive values replaced by the redacted placeholder.
// The original map is never modified.
func Redact(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if opts.Enabled && isSensitive(k, opts.Patterns) {
			out[k] = redactedValue
		} else {
			out[k] = v
		}
	}
	return out
}

// isSensitive reports whether key contains any of the given patterns (case-insensitive).
func isSensitive(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}
