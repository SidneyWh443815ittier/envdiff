// Package interpolator expands shell-style variable references within env values.
// For example, a value like "${BASE_URL}/api" will be expanded using the
// provided env map, leaving unresolved references intact or failing fast
// depending on the options supplied.
package interpolator

import (
	"fmt"
	"regexp"
	"strings"
)

// refPattern matches ${VAR} and $VAR style references.
var refPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Z_][A-Z0-9_]*)`)

// Options controls interpolator behaviour.
type Options struct {
	// FailOnMissing causes Interpolate to return an error when a referenced
	// variable is not present in the env map.
	FailOnMissing bool
	// Placeholder is substituted for missing references when FailOnMissing is
	// false. Defaults to empty string.
	Placeholder string
}

// DefaultOptions returns a safe default configuration.
func DefaultOptions() Options {
	return Options{
		FailOnMissing: false,
		Placeholder:   "",
	}
}

// Interpolate expands variable references in every value of env using the
// values already present in env. It returns a new map and does not mutate the
// original.
func Interpolate(env map[string]string, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		expanded, err := expandValue(v, env, opts)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		out[k] = expanded
	}
	return out, nil
}

func expandValue(value string, env map[string]string, opts Options) (string, error) {
	var expandErr error
	result := refPattern.ReplaceAllStringFunc(value, func(match string) string {
		if expandErr != nil {
			return match
		}
		name := extractName(match)
		if resolved, ok := env[name]; ok {
			return resolved
		}
		if opts.FailOnMissing {
			expandErr = fmt.Errorf("undefined variable %q", name)
			return match
		}
		return opts.Placeholder
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}

func extractName(match string) string {
	match = strings.TrimPrefix(match, "$")
	match = strings.TrimPrefix(match, "{")
	match = strings.TrimSuffix(match, "}")
	return match
}
