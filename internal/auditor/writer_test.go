package auditor_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/envdiff/internal/auditor"
	"github.com/user/envdiff/internal/comparator"
)

func makeEntryAt(ts time.Time, hadIssues bool) auditor.Entry {
	return auditor.Entry{
		Timestamp: ts,
		BaseFile:  "base.env",
		CompFiles: []string{"prod.env"},
		Result:    comparator.Result{},
		HadIssues: hadIssues,
	}
}

func TestWriteReport_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	if err := auditor.WriteReport(&buf, nil); err != nil {
		t.Fatalf("WriteReport: %v", err)
	}
	if !strings.Contains(buf.String(), "No audit entries") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWriteReport_WithEntries(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	entries := []auditor.Entry{
		makeEntryAt(now, false),
		makeEntryAt(now.Add(time.Minute), true),
	}

	var buf bytes.Buffer
	if err := auditor.WriteReport(&buf, entries); err != nil {
		t.Fatalf("WriteReport: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "TIMESTAMP") {
		t.Error("expected header row")
	}
	if !strings.Contains(out, "yes") {
		t.Error("expected 'yes' for issues")
	}
	if !strings.Contains(out, "none") {
		t.Error("expected 'none' for no issues")
	}
}

func TestIssueCount(t *testing.T) {
	now := time.Now()
	entries := []auditor.Entry{
		makeEntryAt(now, true),
		makeEntryAt(now, false),
		makeEntryAt(now, true),
	}
	if got := auditor.IssueCount(entries); got != 2 {
		t.Errorf("expected 2, got %d", got)
	}
}

func TestIssueCount_Empty(t *testing.T) {
	if got := auditor.IssueCount(nil); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}
