package snapshot

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// Manager provides higher-level operations on top of Save/Load.
type Manager struct {
	Dir string
}

// NewManager creates a Manager that stores snapshots under dir.
func NewManager(dir string) *Manager {
	return &Manager{Dir: dir}
}

// SaveNamed stores a snapshot under a logical name (e.g. "production").
func (m *Manager) SaveNamed(name string, snap Snapshot) error {
	path := m.pathFor(name)
	return Save(path, snap)
}

// LoadNamed retrieves a snapshot by logical name.
func (m *Manager) LoadNamed(name string) (Snapshot, error) {
	return Load(m.pathFor(name))
}

// Compare loads two named snapshots and returns the drift report between them.
// prev is the baseline snapshot name and curr is the name of the snapshot to
// compare against it.
func (m *Manager) Compare(prev, curr string) (DriftReport, error) {
	prevSnap, err := m.LoadNamed(prev)
	if err != nil {
		return DriftReport{}, fmt.Errorf("loading previous snapshot %q: %w", prev, err)
	}
	currSnap, err := m.LoadNamed(curr)
	if err != nil {
		return DriftReport{}, fmt.Errorf("loading current snapshot %q: %w", curr, err)
	}
	return Diff(prevSnap, currSnap), nil
}

// WriteReport writes a human-readable drift report to w.
func WriteReport(w io.Writer, prev, curr Snapshot, report DriftReport) {
	fmt.Fprintf(w, "Snapshot drift report\n")
	fmt.Fprintf(w, "  Previous : %s\n", formatTime(prev.CreatedAt))
	fmt.Fprintf(w, "  Current  : %s\n", formatTime(curr.CreatedAt))

	if !report.HasDrift() {
		fmt.Fprintln(w, "  Status   : no drift detected")
		return
	}

	fmt.Fprintln(w, "  Status   : drift detected")

	if len(report.NewMissing) > 0 {
		sort.Strings(report.NewMissing)
		fmt.Fprintln(w, "  New missing keys:")
		for _, k := range report.NewMissing {
			fmt.Fprintf(w, "    - %s\n", k)
		}
	}

	if len(report.NewExtra) > 0 {
		sort.Strings(report.NewExtra)
		fmt.Fprintln(w, "  New extra keys:")
		for _, k := range report.NewExtra {
			fmt.Fprintf(w, "    + %s\n", k)
		}
	}
}

func (m *Manager) pathFor(name string) string {
	return fmt.Sprintf("%s/%s.snapshot.json", m.Dir, name)
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "unknown"
	}
	return t.UTC().Format(time.RFC3339)
}
