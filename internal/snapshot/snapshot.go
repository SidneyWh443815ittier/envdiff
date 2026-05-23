// Package snapshot provides functionality to save and load .env comparison
// snapshots to disk, enabling drift detection between runs.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envdiff/internal/comparator"
)

// Snapshot captures the result of a comparison at a point in time.
type Snapshot struct {
	CreatedAt time.Time          `json:"created_at"`
	BaseFile  string             `json:"base_file"`
	CompFiles []string           `json:"comp_files"`
	Result    comparator.Result  `json:"result"`
}

// Save writes a snapshot to the given file path as JSON.
func Save(path string, snap Snapshot) error {
	snap.CreatedAt = time.Now().UTC()
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write %s: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: read %s: %w", path, err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: unmarshal: %w", err)
	}
	return snap, nil
}

// Diff compares a new result against a previously saved snapshot and returns
// keys that have newly appeared or disappeared since the snapshot was taken.
func Diff(prev, curr comparator.Result) DriftReport {
	report := DriftReport{}

	for _, k := range curr.Missing {
		if !contains(prev.Missing, k) {
			report.NewMissing = append(report.NewMissing, k)
		}
	}
	for _, k := range curr.Extra {
		if !contains(prev.Extra, k) {
			report.NewExtra = append(report.NewExtra, k)
		}
	}
	return report
}

// DriftReport holds keys that are new issues since the last snapshot.
type DriftReport struct {
	NewMissing []string `json:"new_missing"`
	NewExtra   []string `json:"new_extra"`
}

// HasDrift returns true if the report contains any new issues.
func (d DriftReport) HasDrift() bool {
	return len(d.NewMissing) > 0 || len(d.NewExtra) > 0
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
