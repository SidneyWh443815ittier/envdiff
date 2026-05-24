// Package auditor records and reports a history of diff runs for audit purposes.
package auditor

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envdiff/internal/comparator"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp time.Time          `json:"timestamp"`
	BaseFile  string             `json:"base_file"`
	CompFiles []string           `json:"comp_files"`
	Result    comparator.Result  `json:"result"`
	HadIssues bool               `json:"had_issues"`
}

// Options configures auditor behaviour.
type Options struct {
	LogPath string
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		LogPath: "envdiff-audit.log",
	}
}

// Record appends an audit entry to the log file.
func Record(entry Entry, opts Options) error {
	f, err := os.OpenFile(opts.LogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("auditor: open log: %w", err)
	}
	defer f.Close()

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("auditor: marshal entry: %w", err)
	}
	_, err = fmt.Fprintln(f, string(data))
	return err
}

// Load reads all audit entries from the log file.
func Load(opts Options) ([]Entry, error) {
	data, err := os.ReadFile(opts.LogPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("auditor: read log: %w", err)
	}

	var entries []Entry
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("auditor: parse entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
