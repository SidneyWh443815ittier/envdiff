// Package sorter provides utilities for sorting and ranking .env keys
// by various criteria such as alphabetical order, change frequency, or severity.
package sorter

import (
	"sort"

	"github.com/user/envdiff/internal/comparator"
)

// SortBy defines the sort strategy.
type SortBy int

const (
	SortByKey SortBy = iota
	SortBySeverity
	SortByCategory
)

// Options controls how results are sorted.
type Options struct {
	SortBy    SortBy
	Descending bool
}

// DefaultOptions returns the default sort options.
func DefaultOptions() Options {
	return Options{
		SortBy:    SortByKey,
		Descending: false,
	}
}

// severityRank assigns a numeric rank to each category for severity sorting.
func severityRank(category string) int {
	switch category {
	case "missing":
		return 3
	case "extra":
		return 2
	case "mismatch":
		return 1
	default:
		return 0
	}
}

// Entry represents a flat, sortable diff entry.
type Entry struct {
	Key      string
	Category string // "missing", "extra", "mismatch"
	BaseVal  string
	CompVal  string
}

// Flatten converts a comparator.Result into a slice of sortable Entry values.
func Flatten(result comparator.Result) []Entry {
	var entries []Entry
	for _, k := range result.Missing {
		entries = append(entries, Entry{Key: k, Category: "missing"})
	}
	for _, k := range result.Extra {
		entries = append(entries, Entry{Key: k, Category: "extra"})
	}
	for k, v := range result.Mismatched {
		entries = append(entries, Entry{Key: k, Category: "mismatch", BaseVal: v[0], CompVal: v[1]})
	}
	return entries
}

// Sort sorts a slice of Entry values according to the given Options.
func Sort(entries []Entry, opts Options) []Entry {
	sorted := make([]Entry, len(entries))
	copy(sorted, entries)

	sort.SliceStable(sorted, func(i, j int) bool {
		var less bool
		switch opts.SortBy {
		case SortBySeverity:
			ri, rj := severityRank(sorted[i].Category), severityRank(sorted[j].Category)
			if ri != rj {
				less = ri > rj
			} else {
				less = sorted[i].Key < sorted[j].Key
			}
		case SortByCategory:
			if sorted[i].Category != sorted[j].Category {
				less = sorted[i].Category < sorted[j].Category
			} else {
				less = sorted[i].Key < sorted[j].Key
			}
		default: // SortByKey
			less = sorted[i].Key < sorted[j].Key
		}
		if opts.Descending {
			return !less
		}
		return less
	})
	return sorted
}
