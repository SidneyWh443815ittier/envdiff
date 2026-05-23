package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/envdiff/internal/parser"
)

// EnvMap is a map of environment variable key-value pairs.
type EnvMap = map[string]string

// LoadEnvFiles loads one or more .env files and returns their parsed contents.
// If multiple files are provided, keys from later files override earlier ones.
func LoadEnvFiles(paths []string) (EnvMap, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("no env files provided")
	}

	merged := make(EnvMap)

	for _, path := range paths {
		clean := filepath.Clean(path)
		if _, err := os.Stat(clean); os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", clean)
		}

		envMap, err := parser.ParseFile(clean)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", clean, err)
		}

		for k, v := range envMap {
			merged[k] = v
		}
	}

	return merged, nil
}

// LoadEnvFile loads a single .env file and returns its parsed contents.
func LoadEnvFile(path string) (EnvMap, error) {
	return LoadEnvFiles([]string{path})
}
