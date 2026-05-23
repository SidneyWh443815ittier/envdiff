package differ_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/differ"
	"github.com/user/envdiff/internal/filter"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestRun_NoDiff(t *testing.T) {
	base := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")
	comp := writeTempEnv(t, "FOO=bar\nBAZ=qux\n")

	results, err := differ.Run(base, []string{comp}, differ.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if len(r.Diff.Missing)+len(r.Diff.Extra)+len(r.Diff.Mismatched) != 0 {
		t.Errorf("expected no diff, got %+v", r.Diff)
	}
}

func TestRun_WithDiff(t *testing.T) {
	base := writeTempEnv(t, "FOO=bar\nBAZ=qux\nSECRET=x\n")
	comp := writeTempEnv(t, "FOO=changed\nEXTRA=yes\n")

	results, err := differ.Run(base, []string{comp}, differ.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := results[0]
	if len(r.Diff.Missing) != 2 {
		t.Errorf("expected 2 missing keys, got %d", len(r.Diff.Missing))
	}
	if len(r.Diff.Extra) != 1 {
		t.Errorf("expected 1 extra key, got %d", len(r.Diff.Extra))
	}
	if len(r.Diff.Mismatched) != 1 {
		t.Errorf("expected 1 mismatch, got %d", len(r.Diff.Mismatched))
	}
}

func TestRun_BaseNotFound(t *testing.T) {
	_, err := differ.Run(filepath.Join(t.TempDir(), "missing.env"), []string{}, differ.Options{})
	if err == nil {
		t.Error("expected error for missing base file")
	}
}

func TestRun_WithFilter(t *testing.T) {
	base := writeTempEnv(t, "APP_FOO=1\nDB_HOST=localhost\n")
	comp := writeTempEnv(t, "APP_FOO=2\n")

	opts := differ.Options{FilterOpts: filter.Options{Prefix: "APP_"}}
	results, err := differ.Run(base, []string{comp}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := results[0]
	if len(r.Diff.Missing) != 0 {
		t.Errorf("DB_HOST should be filtered out, got missing: %v", r.Diff.Missing)
	}
}

func TestHasIssues(t *testing.T) {
	base := writeTempEnv(t, "A=1\n")
	comp := writeTempEnv(t, "B=2\n")

	results, err := differ.Run(base, []string{comp}, differ.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if !differ.HasIssues(results) {
		t.Error("expected HasIssues to return true")
	}
}
