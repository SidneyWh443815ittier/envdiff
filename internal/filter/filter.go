package filter

import (
	"strings"

	"github.com/user/envdiff/internal/comparator"
)

// Options holds filtering configuration.
type Options struct {
	// Prefix filters keys to only those starting with the given prefix.
	Prefix string
	// IgnoreKeys is a set of key names to exclude from results.
	IgnoreKeys map[string]struct{}
}

// NewOptions creates an Options with the given prefix and ignore list.
func NewOptions(prefix string, ignoreKeys []string) Options {
	ignore := make(map[string]struct{}, len(ignoreKeys))
	for _, k := range ignoreKeys {
		ignore[k] = struct{}{}
	}
	return Options{
		Prefix:     prefix,
		IgnoreKeys: ignore,
	}
}

// Apply filters a comparator.Result according to the given Options.
// Keys that do not match the prefix or are in the ignore list are removed.
func Apply(result comparator.Result, opts Options) comparator.Result {
	filtered := comparator.Result{
		Missing:   filterKeys(result.Missing, opts),
		Extra:     filterKeys(result.Extra, opts),
		Mismatched: filterMismatched(result.Mismatched, opts),
	}
	return filtered
}

func filterKeys(keys []string, opts Options) []string {
	var out []string
	for _, k := range keys {
		if shouldInclude(k, opts) {
			out = append(out, k)
		}
	}
	return out
}

func filterMismatched(mm []comparator.Mismatch, opts Options) []comparator.Mismatch {
	var out []comparator.Mismatch
	for _, m := range mm {
		if shouldInclude(m.Key, opts) {
			out = append(out, m)
		}
	}
	return out
}

func shouldInclude(key string, opts Options) bool {
	if _, ignored := opts.IgnoreKeys[key]; ignored {
		return false
	}
	if opts.Prefix != "" && !strings.HasPrefix(key, opts.Prefix) {
		return false
	}
	return true
}
