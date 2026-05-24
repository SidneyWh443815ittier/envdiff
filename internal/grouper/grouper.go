// Package grouper organises sorted diff entries into named groups
// based on key prefix, category, or a custom label function.
package grouper

import (
	"strings"

	"github.com/user/envdiff/internal/sorter"
)

// Strategy defines how entries are grouped.
type Strategy int

const (
	GroupByPrefix   Strategy = iota // group by KEY prefix before "_"
	GroupByCategory                 // group by diff category
)

// Group holds a named collection of entries.
type Group struct {
	Name    string
	Entries []sorter.Entry
}

// Options controls grouping behaviour.
type Options struct {
	Strategy Strategy
}

// DefaultOptions returns sensible grouper defaults.
func DefaultOptions() Options {
	return Options{Strategy: GroupByCategory}
}

// Group partitions entries into named groups according to opts.
func Group(entries []sorter.Entry, opts Options) []Group {
	index := map[string]*Group{}
	var order []string

	for _, e := range entries {
		var key string
		switch opts.Strategy {
		case GroupByPrefix:
			key = prefixOf(e.Key)
		default: // GroupByCategory
			key = e.Category
		}
		if _, ok := index[key]; !ok {
			index[key] = &Group{Name: key}
			order = append(order, key)
		}
		index[key].Entries = append(index[key].Entries, e)
	}

	result := make([]Group, 0, len(order))
	for _, name := range order {
		result = append(result, *index[name])
	}
	return result
}

// prefixOf returns the segment before the first underscore, or the full key.
func prefixOf(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}
