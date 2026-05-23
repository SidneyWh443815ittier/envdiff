package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Config holds all CLI configuration parsed from flags and arguments.
type Config struct {
	BaseFile   string
	CompFiles  []string
	Prefix     string
	IgnoreKeys []string
	Quiet      bool
	FailOnDiff bool
	JSONOutput bool
}

// Parse parses command-line arguments and returns a Config.
// It returns an error if required arguments are missing.
func Parse(args []string) (*Config, error) {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var prefix string
	var ignoreRaw string
	var quiet bool
	var failOnDiff bool
	var jsonOutput bool

	fs.StringVar(&prefix, "prefix", "", "Only compare keys with this prefix")
	fs.StringVar(&ignoreRaw, "ignore", "", "Comma-separated list of keys to ignore")
	fs.BoolVar(&quiet, "quiet", false, "Suppress output, only use exit code")
	fs.BoolVar(&failOnDiff, "fail", false, "Exit with non-zero code if any diff is found")
	fs.BoolVar(&jsonOutput, "json", false, "Output results as JSON")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	remaining := fs.Args()
	if len(remaining) < 2 {
		return nil, fmt.Errorf("usage: envdiff [flags] <base.env> <compare.env> [compare2.env ...]")
	}

	var ignoreKeys []string
	if ignoreRaw != "" {
		for _, k := range strings.Split(ignoreRaw, ",") {
			trimmed := strings.TrimSpace(k)
			if trimmed != "" {
				ignoreKeys = append(ignoreKeys, trimmed)
			}
		}
	}

	return &Config{
		BaseFile:   remaining[0],
		CompFiles:  remaining[1:],
		Prefix:     prefix,
		IgnoreKeys: ignoreKeys,
		Quiet:      quiet,
		FailOnDiff: failOnDiff,
		JSONOutput: jsonOutput,
	}, nil
}
