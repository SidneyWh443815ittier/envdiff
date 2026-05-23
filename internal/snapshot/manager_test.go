package snapshot_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/snapshot"
)

func TestManager_SaveAndLoadNamed(t *testing.T) {
	dir := t.TempDir()
	m := snapshot.NewManager(dir)

	snap := snapshot.Snapshot{
		BaseFile:  ".env",
		CompFiles: []string{".env.staging"},
		Result:    makeResult(),
	}

	if err := m.SaveNamed("staging", snap); err != nil {
		t.Fatalf("SaveNamed: %v", err)
	}

	loaded, err := m.LoadNamed("staging")
	if err != nil {
		t.Fatalf("LoadNamed: %v", err)
	}

	if loaded.BaseFile != ".env" {
		t.Errorf("BaseFile: got %q", loaded.BaseFile)
	}
}

func TestManager_LoadNamed_NotFound(t *testing.T) {
	m := snapshot.NewManager(t.TempDir())
	_, err := m.LoadNamed("nonexistent")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWriteReport_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	prev := snapshot.Snapshot{CreatedAt: time.Now().Add(-time.Hour)}
	curr := snapshot.Snapshot{CreatedAt: time.Now()}
	report := snapshot.DriftReport{}

	snapshot.WriteReport(&buf, prev, curr, report)

	if !strings.Contains(buf.String(), "no drift detected") {
		t.Errorf("expected no-drift message, got:\n%s", buf.String())
	}
}

func TestWriteReport_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	prev := snapshot.Snapshot{CreatedAt: time.Now().Add(-time.Hour)}
	curr := snapshot.Snapshot{CreatedAt: time.Now()}
	report := snapshot.DriftReport{
		NewMissing: []string{"DB_PASS", "API_KEY"},
		NewExtra:   []string{"LEGACY_TOKEN"},
	}

	snapshot.WriteReport(&buf, prev, curr, report)
	out := buf.String()

	if !strings.Contains(out, "drift detected") {
		t.Errorf("expected drift detected, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_PASS") {
		t.Errorf("expected DB_PASS in output")
	}
	if !strings.Contains(out, "LEGACY_TOKEN") {
		t.Errorf("expected LEGACY_TOKEN in output")
	}
}

func TestWriteReport_ZeroTimes(t *testing.T) {
	var buf bytes.Buffer
	snapshot.WriteReport(&buf,
		snapshot.Snapshot{},
		snapshot.Snapshot{},
		snapshot.DriftReport{},
	)
	out := buf.String()
	if !strings.Contains(out, "unknown") {
		t.Errorf("expected 'unknown' for zero times, got:\n%s", out)
	}
}

func TestDiff_IgnoresPreexistingIssues(t *testing.T) {
	prev := comparator.Result{Missing: []string{"ALREADY_MISSING"}}
	curr := comparator.Result{Missing: []string{"ALREADY_MISSING"}}

	report := snapshot.Diff(prev, curr)
	if report.HasDrift() {
		t.Errorf("should not report pre-existing issues as drift: %+v", report)
	}
}
