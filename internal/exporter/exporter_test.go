package exporter_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/comparator"
	"github.com/user/envdiff/internal/exporter"
)

func cleanResult() comparator.Result {
	return comparator.Result{}
}

func diffResult() comparator.Result {
	return comparator.Result{
		Missing: map[string]string{"DB_HOST": "localhost"},
		Extra:   map[string]string{"UNUSED_KEY": "value"},
		Mismatched: map[string]comparator.Mismatch{
			"APP_ENV": {Base: "production", Comp: "staging"},
		},
	}
}

func TestExport_TextClean(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.txt")
	opts := exporter.Options{Format: exporter.FormatText, FilePath: tmp}

	if err := exporter.Export(cleanResult(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	if !strings.Contains(string(data), "No differences") {
		t.Errorf("expected clean message, got: %s", data)
	}
}

func TestExport_TextDiff(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.txt")
	opts := exporter.Options{Format: exporter.FormatText, FilePath: tmp}

	if err := exporter.Export(diffResult(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	content := string(data)

	if !strings.Contains(content, "MISSING  DB_HOST") {
		t.Errorf("expected MISSING line, got: %s", content)
	}
	if !strings.Contains(content, "EXTRA    UNUSED_KEY") {
		t.Errorf("expected EXTRA line, got: %s", content)
	}
	if !strings.Contains(content, "MISMATCH APP_ENV") {
		t.Errorf("expected MISMATCH line, got: %s", content)
	}
}

func TestExport_JSONDiff(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.json")
	opts := exporter.Options{Format: exporter.FormatJSON, FilePath: tmp}

	if err := exporter.Export(diffResult(), opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(tmp)
	var result comparator.Result
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if _, ok := result.Missing["DB_HOST"]; !ok {
		t.Errorf("expected DB_HOST in missing")
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "out.xyz")
	opts := exporter.Options{Format: "xml", FilePath: tmp}

	err := exporter.Export(cleanResult(), opts)
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestExport_InvalidPath(t *testing.T) {
	opts := exporter.Options{Format: exporter.FormatText, FilePath: "/nonexistent/dir/out.txt"}
	err := exporter.Export(cleanResult(), opts)
	if err == nil {
		t.Fatal("expected error for invalid path")
	}
}
