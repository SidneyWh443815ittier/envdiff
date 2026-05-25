package patcher_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/patcher"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("readFile: %v", err)
	}
	return string(b)
}

func TestPatch_UpdatesExistingKey(t *testing.T) {
	p := writeTempEnv(t, "APP_ENV=development\nDEBUG=false\n")
	err := patcher.Patch(p, map[string]string{"APP_ENV": "production"}, patcher.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := readFile(t, p)
	if !strings.Contains(got, "APP_ENV=production") {
		t.Errorf("expected updated value, got:\n%s", got)
	}
	if !strings.Contains(got, "DEBUG=false") {
		t.Errorf("expected unchanged key to be preserved, got:\n%s", got)
	}
}

func TestPatch_PreservesCommentsAndBlanks(t *testing.T) {
	src := "# app settings\nAPP_ENV=dev\n\nDEBUG=true\n"
	p := writeTempEnv(t, src)
	err := patcher.Patch(p, map[string]string{"DEBUG": "false"}, patcher.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := readFile(t, p)
	if !strings.Contains(got, "# app settings") {
		t.Errorf("comment not preserved:\n%s", got)
	}
}

func TestPatch_AppendsMissingKey(t *testing.T) {
	p := writeTempEnv(t, "APP_ENV=dev\n")
	opts := patcher.DefaultOptions()
	opts.CreateMissing = true
	err := patcher.Patch(p, map[string]string{"NEW_KEY": "hello"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := readFile(t, p)
	if !strings.Contains(got, "NEW_KEY=hello") {
		t.Errorf("expected new key appended, got:\n%s", got)
	}
}

func TestPatch_SkipsMissingKeyWhenDisabled(t *testing.T) {
	p := writeTempEnv(t, "APP_ENV=dev\n")
	opts := patcher.DefaultOptions()
	opts.CreateMissing = false
	err := patcher.Patch(p, map[string]string{"GHOST": "value"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := readFile(t, p)
	if strings.Contains(got, "GHOST") {
		t.Errorf("expected key to be skipped, got:\n%s", got)
	}
}

func TestPatch_QuotesValues(t *testing.T) {
	p := writeTempEnv(t, "SECRET=old\n")
	opts := patcher.DefaultOptions()
	opts.QuoteValues = true
	err := patcher.Patch(p, map[string]string{"SECRET": "new value"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := readFile(t, p)
	if !strings.Contains(got, `SECRET="new value"`) {
		t.Errorf("expected quoted value, got:\n%s", got)
	}
}

func TestPatch_CreatesFileWhenMissing(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "new.env")
	opts := patcher.DefaultOptions()
	err := patcher.Patch(p, map[string]string{"BRAND_NEW": "yes"}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := readFile(t, p)
	if !strings.Contains(got, "BRAND_NEW=yes") {
		t.Errorf("expected key written to new file, got:\n%s", got)
	}
}
