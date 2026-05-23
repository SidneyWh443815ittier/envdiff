// Package pipeline wires together the differ, validator, summary, and
// output formatter into a single callable unit used by the CLI.
package pipeline

import (
	"io"

	"github.com/user/envdiff/internal/config"
	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/filter"
	"github.com/user/envdiff/internal/output"
	"github.com/user/envdiff/internal/summary"
)

// RunResult captures the overall outcome of a pipeline execution.
type RunResult struct {
	HasIssues bool
	ExitCode  int
}

// Run executes the full envdiff pipeline using the provided config and
// writes formatted output to w. It returns a RunResult indicating whether
// any issues were found.
func Run(cfg config.Config, w io.Writer) (RunResult, error) {
	differOpts := differ.Options{
		FilterOpts: filter.Options{
			Prefix:     cfg.Prefix,
			IgnoreKeys: cfg.IgnoreKeys,
		},
	}

	results, err := differ.Run(cfg.BaseFile, cfg.CompFiles, differOpts)
	if err != nil {
		return RunResult{ExitCode: 2}, err
	}

	fmt := output.New(cfg.Format)
	for _, r := range results {
		if err := fmt.Write(w, r.BaseFile, r.CompFile, r.Diff); err != nil {
			return RunResult{ExitCode: 2}, err
		}
	}

	hasIssues := differ.HasIssues(results)

	var allDiffs []interface{}
	_ = allDiffs // summary operates per-result
	for _, r := range results {
		summary.Write(w, summary.Compute(r.Diff))
	}

	exitCode := 0
	if hasIssues && cfg.FailOnDiff {
		exitCode = 1
	}

	return RunResult{HasIssues: hasIssues, ExitCode: exitCode}, nil
}
