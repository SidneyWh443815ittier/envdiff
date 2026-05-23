package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envdiff/internal/loader"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

func TestLoadEnvFile_Basic(t *testing.T) {
	path := writeTempEnv(t, "KEY=value\nFOO=bar\n")
	env, err := loader.LoadEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY"] != "value" || env["FOO"] != "bar" {
		t.Errorf("unexpected env map: %v", env)
	}
}

func TestLoadEnvFile_NotFound(t *testing.T) {
	_, err := loader.LoadEnvFile("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadEnvFiles_NoPaths(t *testing.T) {
	_, err := loader.LoadEnvFiles([]string{})
	if err == nil {
		t.Fatal("expected error for empty paths, got nil")
	}
}

func TestLoadEnvFiles_MergesFiles(t *testing.T) {
	path1 := writeTempEnv(t, "KEY=original\nONLY_IN_FIRST=yes\n")
	path2 := writeTempEnv(t, "KEY=overridden\nONLY_IN_SECOND=yes\n")

	env, err := loader.LoadEnvFiles([]string{path1, path2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["KEY"] != "overridden" {
		t.Errorf("expected KEY=overridden, got %s", env["KEY"])
	}
	if env["ONLY_IN_FIRST"] != "yes" {
		t.Errorf("expected ONLY_IN_FIRST=yes, got %s", env["ONLY_IN_FIRST"])
	}
	if env["ONLY_IN_SECOND"] != "yes" {
		t.Errorf("expected ONLY_IN_SECOND=yes, got %s", env["ONLY_IN_SECOND"])
	}
}

func TestLoadEnvFiles_OneFileMissing(t *testing.T) {
	path1 := writeTempEnv(t, "KEY=value\n")
	_, err := loader.LoadEnvFiles([]string{path1, "/does/not/exist/.env"})
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
