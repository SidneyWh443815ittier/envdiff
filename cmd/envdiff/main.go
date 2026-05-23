package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/envdiff/internal/comparator"
	"github.com/envdiff/internal/loader"
	"github.com/envdiff/internal/reporter"
)

func main() {
	baseFlag := flag.String("base", "", "Path to the base .env file (required)")
	targetFlag := flag.String("target", "", "Path to the target .env file to compare against (required)")
	ciFlag := flag.Bool("ci", false, "Exit with non-zero status code if differences are found")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff -base <file> -target <file> [options]\n\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *baseFlag == "" || *targetFlag == "" {
		fmt.Fprintln(os.Stderr, "error: -base and -target flags are required")
		flag.Usage()
		os.Exit(2)
	}

	baseEnv, err := loader.LoadEnvFile(*baseFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading base file: %v\n", err)
		os.Exit(2)
	}

	targetEnv, err := loader.LoadEnvFile(*targetFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading target file: %v\n", err)
		os.Exit(2)
	}

	result := comparator.Compare(baseEnv, targetEnv)
	reporter.ReportToStdout(result)

	if *ciFlag {
		os.Exit(reporter.ExitCode(result))
	}
}
