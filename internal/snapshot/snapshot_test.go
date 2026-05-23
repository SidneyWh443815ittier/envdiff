package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/snapshot"
)

func makeResult() comparator.Result {
	return comparator.Result{
		Missing: []string{"DB_HOST"},
		Extra:   []string{"OLD_KEY"},
		Mismatched: map[string][2]string{
			"PORT": {"3000", "4000"},
		},
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	snap := snapshot.Snapshot{
		BaseFile:  ".env",
		CompFiles: []string{".env.prod"},
		Result:    makeResult(),
	}

	if err := snapshot.Save(path, snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if loaded.BaseFile != snap.BaseFile {
		t.Errorf("BaseFile: got %q, want %q", loaded.BaseFile, snap.BaseFile)
	}
	if loaded.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	if len(loaded.Result.Missing) != 1 || loaded.Result.Missing[0] != "DB_HOST" {
		t.Errorf("Missing keys not preserved: %v", loaded.Result.Missing)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestSave_InvalidPath(t *testing.T) {
	snap := snapshot.Snapshot{CreatedAt: time.Now()}
	err := snapshot.Save("/nonexistent/dir/snap.json", snap)
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json{"), 0o644)
	_, err := snapshot.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestDiff_NewIssues(t *testing.T) {
	prev := comparator.Result{
		Missing: []string{"OLD_MISSING"},
		Extra:   []string{},
	}
	curr := comparator.Result{
		Missing: []string{"OLD_MISSING", "NEW_MISSING"},
		Extra:   []string{"NEW_EXTRA"},
	}

	report := snapshot.Diff(prev, curr)

	if len(report.NewMissing) != 1 || report.NewMissing[0] != "NEW_MISSING" {
		t.Errorf("NewMissing: got %v", report.NewMissing)
	}
	if len(report.NewExtra) != 1 || report.NewExtra[0] != "NEW_EXTRA" {
		t.Errorf("NewExtra: got %v", report.NewExtra)
	}
	if !report.HasDrift() {
		t.Error("expected HasDrift to be true")
	}
}

func TestDiff_NoDrift(t *testing.T) {
	result := makeResult()
	report := snapshot.Diff(result, result)
	if report.HasDrift() {
		t.Errorf("expected no drift, got %+v", report)
	}
}
