// Package merger provides utilities for merging multiple .env maps
// into a single unified map, with configurable conflict resolution.
package merger

import "fmt"

// Strategy defines how key conflicts are resolved during a merge.
type Strategy int

const (
	// StrategyFirst keeps the value from the first file that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast overwrites with the value from the last file that defines the key.
	StrategyLast
	// StrategyError returns an error if the same key appears with different values.
	StrategyError
)

// Options configures the merge behaviour.
type Options struct {
	Strategy Strategy
}

// DefaultOptions returns Options with the StrategyLast conflict resolution.
func DefaultOptions() Options {
	return Options{Strategy: StrategyLast}
}

// Merge combines multiple env maps into one according to the given Options.
// The maps are processed in order; conflict resolution depends on opts.Strategy.
func Merge(maps []map[string]string, opts Options) (map[string]string, error) {
	result := make(map[string]string)

	for _, m := range maps {
		for k, v := range m {
			existing, exists := result[k]
			if !exists {
				result[k] = v
				continue
			}

			switch opts.Strategy {
			case StrategyFirst:
				// keep existing value — do nothing
			case StrategyLast:
				result[k] = v
			case StrategyError:
				if existing != v {
					return nil, fmt.Errorf("merger: conflict on key %q: %q vs %q", k, existing, v)
				}
			}
		}
	}

	return result, nil
}
