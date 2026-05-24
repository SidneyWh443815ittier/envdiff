// Package resolver resolves environment variable references within .env files,
// expanding values that reference other keys (e.g. BASE_URL=${SCHEME}://${HOST}).
package resolver

import (
	"fmt"
	"regexp"
	"strings"
)

// Options controls resolver behaviour.
type Options struct {
	// MaxDepth limits recursive expansion to prevent infinite loops.
	MaxDepth int
	// FailOnMissing returns an error when a referenced key is not present.
	FailOnMissing bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		MaxDepth:      10,
		FailOnMissing: false,
	}
}

var refPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// Resolve expands variable references in env values using the same map.
// Values like FOO=${BAR} are replaced with the value of BAR.
// Unresolved references are left as-is unless FailOnMissing is set.
func Resolve(env map[string]string, opts Options) (map[string]string, error) {
	out := make(map[string]string, len(env))
	for k, v := range env {
		resolved, err := expand(v, env, opts, 0)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}

func expand(value string, env map[string]string, opts Options, depth int) (string, error) {
	if depth > opts.MaxDepth {
		return value, fmt.Errorf("max expansion depth %d exceeded", opts.MaxDepth)
	}
	var expandErr error
	result := refPattern.ReplaceAllStringFunc(value, func(match string) string {
		if expandErr != nil {
			return match
		}
		key := strings.TrimSpace(match[2 : len(match)-1])
		val, ok := env[key]
		if !ok {
			if opts.FailOnMissing {
				expandErr = fmt.Errorf("referenced key %q not found", key)
			}
			return match
		}
		expanded, err := expand(val, env, opts, depth+1)
		if err != nil {
			expandErr = err
			return match
		}
		return expanded
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}
