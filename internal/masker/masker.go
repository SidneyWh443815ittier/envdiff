// Package masker provides utilities for masking sensitive env values
// in output, replacing them with a configurable mask string.
package masker

import "strings"

// DefaultMask is the string used to replace sensitive values.
const DefaultMask = "****"

// Options controls masking behaviour.
type Options struct {
	// Enabled toggles masking on or off.
	Enabled bool
	// Mask is the replacement string for sensitive values.
	Mask string
	// SensitiveSubstrings are substrings that, when found in a key
	// (case-insensitive), cause the value to be masked.
	SensitiveSubstrings []string
}

// DefaultOptions returns a sensible default Options.
func DefaultOptions() Options {
	return Options{
		Enabled: true,
		Mask:    DefaultMask,
		SensitiveSubstrings: []string{
			"password", "passwd", "secret", "token",
			"api_key", "apikey", "private", "credential",
		},
	}
}

// Mask returns a new map where sensitive values are replaced with the
// configured mask string. The original map is never modified.
func Mask(env map[string]string, opts Options) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if opts.Enabled && isSensitive(k, opts.SensitiveSubstrings) {
			out[k] = opts.Mask
		} else {
			out[k] = v
		}
	}
	return out
}

// IsSensitiveKey reports whether key is considered sensitive given the
// provided list of substrings.
func IsSensitiveKey(key string, substrings []string) bool {
	return isSensitive(key, substrings)
}

func isSensitive(key string, substrings []string) bool {
	lower := strings.ToLower(key)
	for _, sub := range substrings {
		if strings.Contains(lower, strings.ToLower(sub)) {
			return true
		}
	}
	return false
}
