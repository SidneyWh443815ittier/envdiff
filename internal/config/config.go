package config

import (
	"errors"
	"flag"
	"strings"

	"github.com/user/envdiff/internal/output"
)

// Config holds the parsed CLI configuration.
type Config struct {
	BaseFile    string
	CompFiles   []string
	IgnoreKeys  []string
	Prefix      string
	Format      output.Format
	FailOnDiff  bool
}

type multiString []string

func (m *multiString) String() string  { return strings.Join(*m, ",") }
func (m *multiString) Set(v string) error { *m = append(*m, v); return nil }

// Parse parses os.Args-style arguments and returns a Config.
func Parse(args []string) (*Config, error) {
	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)

	var compFiles multiString
	var ignoreKeys multiString
	prefix := fs.String("prefix", "", "Only compare keys with this prefix")
	format := fs.String("format", "text", "Output format: text, json, markdown")
	failOnDiff := fs.Bool("fail", false, "Exit with non-zero code if differences found")

	fs.Var(&compFiles, "comp", "Comparison .env file (repeatable)")
	fs.Var(&ignoreKeys, "ignore", "Key to ignore (repeatable)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	remaining := fs.Args()
	if len(remaining) < 1 {
		return nil, errors.New("usage: envdiff [flags] <base.env> [comp.env ...]")
	}

	base := remaining[0]
	if len(compFiles) == 0 {
		if len(remaining) < 2 {
			return nil, errors.New("at least one comparison file is required")
		}
		compFiles = append(compFiles, remaining[1:]...)
	}

	fmt := output.Format(*format)
	if fmt != output.FormatText && fmt != output.FormatJSON && fmt != output.FormatMarkdown {
		return nil, errors.New("invalid format: must be text, json, or markdown")
	}

	return &Config{
		BaseFile:   base,
		CompFiles:  []string(compFiles),
		IgnoreKeys: []string(ignoreKeys),
		Prefix:     *prefix,
		Format:     fmt,
		FailOnDiff: *failOnDiff,
	}, nil
}
