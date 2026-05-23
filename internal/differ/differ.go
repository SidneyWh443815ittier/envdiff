// Package differ provides environment-aware diffing between a base
// env file and one or more comparison env files.
package differ

import (
	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/filter"
	"github.com/user/envdiff/internal/loader"
)

// Result holds the diff outcome for a single comparison pair.
type Result struct {
	BaseFile string
	CompFile string
	Diff     comparator.Result
}

// Options controls how diffs are performed.
type Options struct {
	FilterOpts filter.Options
}

// Run loads the base env file and each comp file, compares them,
// applies any filters, and returns one Result per comp file.
func Run(baseFile string, compFiles []string, opts Options) ([]Result, error) {
	base, err := loader.LoadEnvFile(baseFile)
	if err != nil {
		return nil, err
	}

	results := make([]Result, 0, len(compFiles))
	for _, cf := range compFiles {
		comp, err := loader.LoadEnvFile(cf)
		if err != nil {
			return nil, err
		}

		diff := comparator.Compare(base, comp)
		diff = filter.Apply(diff, opts.FilterOpts)

		results = append(results, Result{
			BaseFile: baseFile,
			CompFile: cf,
			Diff:     diff,
		})
	}

	return results, nil
}

// HasIssues returns true if any result contains at least one diff entry.
func HasIssues(results []Result) bool {
	for _, r := range results {
		if len(r.Diff.Missing) > 0 || len(r.Diff.Extra) > 0 || len(r.Diff.Mismatched) > 0 {
			return true
		}
	}
	return false
}
