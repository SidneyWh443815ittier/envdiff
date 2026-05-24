// Package templater generates a .env.template file from an existing env map,
// replacing all values with empty strings or placeholder comments.
package templater

import (
	"fmt"
	"os"\n	"sort"
	"strings"
)

// Options controls template generation behaviour.
type Options struct {
	// Placeholder is written as the value for every key.
	// Defaults to an empty string.
	Placeholder string
	// AddComments prepends a comment above each key when true.
	AddComments bool
}

// DefaultOptions returns sensible defaults for template generation.
func DefaultOptions() Options {
	return Options{
		Placeholder: "",
		AddComments: false,
	}
}

// Generate builds a template string from the supplied env map.
// Keys are emitted in sorted order so output is deterministic.
func Generate(env map[string]string, opts Options) string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		if opts.AddComments {
			fmt.Fprintf(&sb, "# %s\n", k)
		}
		fmt.Fprintf(&sb, "%s=%s\n", k, opts.Placeholder)
	}
	return sb.String()
}

// WriteFile writes the generated template to the given file path.
// The file is created or truncated if it already exists.
func WriteFile(path string, env map[string]string, opts Options) error {
	content := Generate(env, opts)
	return os.WriteFile(path, []byte(content), 0o644)
}
