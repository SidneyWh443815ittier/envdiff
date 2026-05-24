package auditor_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/auditor"
	"github.com/user/envdiff/internal/comparator"
)

func makeEntry(base string, hadIssues bool) auditor.Entry {
	return auditor.Entry{
		Timestamp: time.Now().UTC().Truncate(time.Second),
		BaseFile:  base,
		CompFiles: []string{"comp.env"},
		Result: comparator.Result{
			Missing: map[string]struct{}{"KEY_A": {}},
		},
		HadIssues: hadIssues,
	}
}

func TestRecord_And_Load(t *testing.T) {
	dir := t.TempDir()
	opts := auditor.Options{LogPath: filepath.Join(dir, "audit.log")}

	e := makeEntry("base.env", true)
	if err := auditor.Record(e, opts); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entries, err := auditor.Load(opts)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].BaseFile != "base.env" {
		t.Errorf("BaseFile mismatch: %s", entries[0].BaseFile)
	}
	if !entries[0].HadIssues {
		t.Error("expected HadIssues=true")
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	dir := t.TempDir()
	opts := auditor.Options{LogPath: filepath.Join(dir, "audit.log")}

	for i := 0; i < 3; i++ {
		if err := auditor.Record(makeEntry("base.env", false), opts); err != nil {
			t.Fatalf("Record %d: %v", i, err)
		}
	}

	entries, err := auditor.Load(opts)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestLoad_NoFile(t *testing.T) {
	dir := t.TempDir()
	opts := auditor.Options{LogPath: filepath.Join(dir, "missing.log")}

	entries, err := auditor.Load(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries, got %v", entries)
	}
}

func TestRecord_InvalidPath(t *testing.T) {
	opts := auditor.Options{LogPath: "/nonexistent/dir/audit.log"}
	err := auditor.Record(makeEntry("base.env", false), opts)
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := auditor.DefaultOptions()
	if opts.LogPath == "" {
		t.Error("expected non-empty LogPath")
	}
}

func init() {
	_ = os.Getenv // suppress unused import
}
